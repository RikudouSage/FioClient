package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/mattn/go-sqlite3"
	"github.com/samber/lo"
	"go.chrastecky.dev/fio-client/fioclient/model"
	. "go.chrastecky.dev/fio-client/fioclient/types"
)

var ErrNonUnique = errors.New("the affected row is not unique")
var ErrNoRows = errors.New("no rows found")

type Manager struct {
	db                 *sql.DB
	encryptionProvider EncryptedAPIKeyProvider
}

func NewManager(db *sql.DB, encryptionProvider EncryptedAPIKeyProvider) *Manager {
	return &Manager{
		db:                 db,
		encryptionProvider: encryptionProvider,
	}
}

func (receiver *Manager) redactAPIKey(apiKey string) string {
	if receiver.encryptionProvider != nil {
		return ""
	}

	return apiKey
}

func (receiver *Manager) wrapError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("%w: %w", ErrNoRows, err)
	}

	if sqliteErr, ok := errors.AsType[sqlite3.Error](err); ok {
		if sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique || sqliteErr.ExtendedCode == sqlite3.ErrConstraintPrimaryKey {
			return fmt.Errorf("%w: %w", ErrNonUnique, err)
		}
	}

	return err
}

func (receiver *Manager) StoreAccount(ctx context.Context, account model.Account) error {
	_, err := receiver.db.ExecContext(
		ctx,
		"insert into accounts (account_number, api_key, bank_code, currency, iban, bic) values (?, ?, ?, ?, ?, ?)",
		account.AccountNumber,
		receiver.redactAPIKey(account.APIKey),
		account.BankCode,
		account.Currency,
		account.IBAN,
		account.BIC,
	)
	if err != nil {
		return fmt.Errorf("failed creating account: %w", receiver.wrapError(err))
	}

	if receiver.encryptionProvider != nil {
		err = receiver.encryptionProvider.StoreAPIKey(account.AccountNumber, account.APIKey)
		if err != nil {
			defer func() {
				// best effort
				_ = receiver.RemoveAccountByNumber(ctx, account.AccountNumber)
			}()
			return fmt.Errorf("failed storing API key for account %s: %w", account.AccountNumber, err)
		}
	}

	return nil
}

func (receiver *Manager) FindAccountByNumber(ctx context.Context, accountNumber string) (model.Account, error) {
	var account model.Account
	if err := receiver.
		db.
		QueryRowContext(ctx, "select account_number, api_key, bank_code, currency, iban, bic from accounts where account_number = ?", accountNumber).
		Scan(&account.AccountNumber, &account.APIKey, &account.BankCode, &account.Currency, &account.IBAN, &account.BIC); err != nil {
		return model.Account{}, fmt.Errorf("failed querying account: %w", receiver.wrapError(err))
	}

	if receiver.encryptionProvider != nil {
		var err error
		account.APIKey, err = receiver.encryptionProvider.GetAPIKey(accountNumber)
		if err != nil {
			return model.Account{}, fmt.Errorf("failed getting API key for account %s: %w", account.AccountNumber, err)
		}
	}

	return account, nil
}

func (receiver *Manager) RemoveAccountByNumber(ctx context.Context, accountNumber string) error {
	if receiver.encryptionProvider != nil {
		err := receiver.encryptionProvider.RemoveAPIKey(accountNumber)
		if err != nil {
			return fmt.Errorf("failed removing API key for account %s: %w", accountNumber, err)
		}
	}

	_, err := receiver.db.ExecContext(ctx, "delete from accounts where account_number = ?", accountNumber)
	if err != nil {
		return fmt.Errorf("failed deleting account: %w", receiver.wrapError(err))
	}

	return nil
}

func (receiver *Manager) StoreTransactions(ctx context.Context, transactions []model.Transaction) error {
	values := make([]any, 0)
	query := "insert into transactions (id, account_number, date, amount, currency, counterparty_account, counterparty_name, counterparty_bank_code, counterparty_bank_name, constant_symbol, variable_symbol, specific_symbol, user_identity, transaction_type, performed_by, additional_info, comment, bic, instruction_id, payer_reference, local_only) values "
	subqueries := make([]string, 0, len(transactions))
	for _, transaction := range transactions {
		tmpVals := []any{transaction.ID, transaction.AccountNumber, transaction.Date, transaction.Amount, transaction.Currency, transaction.CounterpartyAccount, transaction.CounterpartyName, transaction.CounterpartyBankCode, transaction.CounterpartyBankName, transaction.ConstantSymbol, transaction.VariableSymbol, transaction.SpecificSymbol, transaction.UserIdentity, transaction.TransactionType, transaction.PerformedBy, transaction.AdditionalInfo, transaction.Comment, transaction.BIC, transaction.InstructionID, transaction.PayerReference, transaction.LocalOnly}
		values = append(values, tmpVals...)
		subqueries = append(subqueries, "("+strings.Join(lo.RepeatBy(len(tmpVals), func(_ int) string {
			return "?"
		}), ", ")+")")
	}

	if len(subqueries) == 0 {
		return nil
	}

	query += strings.Join(subqueries, ", \n")
	query += " on conflict (account_number, id) do update set \"date\" = excluded.date, amount = excluded.amount, currency = excluded.currency, counterparty_account = excluded.counterparty_account, counterparty_name = excluded.counterparty_name, counterparty_bank_code = excluded.counterparty_bank_code, counterparty_bank_name = excluded.counterparty_bank_name, constant_symbol = excluded.constant_symbol, variable_symbol = excluded.variable_symbol, specific_symbol = excluded.specific_symbol, user_identity = excluded.user_identity, transaction_type = excluded.transaction_type, performed_by = excluded.performed_by, additional_info = excluded.additional_info, comment = excluded.comment, bic = excluded.bic, instruction_id = excluded.instruction_id, payer_reference = excluded.payer_reference, local_only = excluded.local_only"

	if _, err := receiver.db.ExecContext(ctx, query, values...); err != nil {
		return fmt.Errorf("failed inserting transactions: %w", err)
	}

	return nil
}

func (receiver *Manager) LoadTransactions(ctx context.Context, options ...LoadOption) ([]model.Transaction, error) {
	opts := loadOptions{}
	for _, option := range options {
		option(&opts)
	}

	var queryBuilder strings.Builder
	queryBuilder.WriteString("select * from transactions")
	if len(opts.where) != 0 {
		queryBuilder.WriteString(" where ")
		subQueries := make([]string, 0, len(opts.where))
		for _, where := range opts.where {
			subQueries = append(subQueries, "("+where+")")
		}
		queryBuilder.WriteString(strings.Join(subQueries, " and "))
	}
	if opts.orderByDirection != "" && opts.orderByField != "" {
		queryBuilder.WriteString(" order by " + opts.orderByField + " " + opts.orderByDirection)
	}
	if opts.limit != 0 {
		queryBuilder.WriteString(" limit " + strconv.FormatUint(uint64(opts.limit), 10))
	}

	rows, err := receiver.db.QueryContext(ctx, queryBuilder.String(), opts.bindValues...)
	if err != nil {
		return nil, fmt.Errorf("failed loading transactions: %w", err)
	}
	defer rows.Close()

	result := make([]model.Transaction, 0)
	for rows.Next() {
		var transaction model.Transaction
		if err := rows.Scan(
			&transaction.ID,
			&transaction.AccountNumber,
			&transaction.Date,
			&transaction.Amount,
			&transaction.Currency,
			&transaction.CounterpartyAccount,
			&transaction.CounterpartyName,
			&transaction.CounterpartyBankCode,
			&transaction.CounterpartyBankName,
			&transaction.ConstantSymbol,
			&transaction.VariableSymbol,
			&transaction.SpecificSymbol,
			&transaction.UserIdentity,
			&transaction.TransactionType,
			&transaction.PerformedBy,
			&transaction.AdditionalInfo,
			&transaction.Comment,
			&transaction.BIC,
			&transaction.InstructionID,
			&transaction.PayerReference,
			&transaction.LocalOnly,
		); err != nil {
			return nil, fmt.Errorf("failed loading transactions: %w", err)
		}

		result = append(result, transaction)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed iterating transactions: %w", err)
	}

	return result, nil
}

func (receiver *Manager) GetAccounts(ctx context.Context) ([]model.Account, error) {
	result := make([]model.Account, 0)

	rows, err := receiver.db.QueryContext(ctx, "select * from accounts")
	if err != nil {
		return nil, fmt.Errorf("failed loading accounts: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var account model.Account
		if err := rows.Scan(&account.AccountNumber, &account.APIKey, &account.BankCode, &account.Currency, &account.IBAN, &account.BIC); err != nil {
			return nil, fmt.Errorf("failed loading accounts: %w", err)
		}

		if receiver.encryptionProvider != nil {
			account.APIKey, err = receiver.encryptionProvider.GetAPIKey(account.AccountNumber)
			if err != nil {
				return nil, fmt.Errorf("failed loading api key for account %s: %w", account.AccountNumber, err)
			}
		}

		result = append(result, account)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed iterating accounts: %w", err)
	}

	return result, nil
}
