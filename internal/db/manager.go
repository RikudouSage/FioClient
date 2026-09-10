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
)

var ErrNonUnique = errors.New("the affected row is not unique")
var ErrNoRows = errors.New("no rows found")

type Manager struct {
	db *sql.DB
}

func NewManager(db *sql.DB) *Manager {
	return &Manager{
		db: db,
	}
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

func (receiver *Manager) StoreAccount(account model.Account) error {
	_, err := receiver.db.Exec(
		"insert into accounts (account_number, api_key, bank_code, currency, iban, bic) values (?, ?, ?, ?, ?, ?)",
		account.AccountNumber,
		account.ApiKey,
		account.BankCode,
		account.Currency,
		account.IBAN,
		account.BIC,
	)
	if err != nil {
		return fmt.Errorf("failed creating account: %w", receiver.wrapError(err))
	}

	return nil
}

func (receiver *Manager) FindAccountByNumber(ctx context.Context, accountNumber string) (model.Account, error) {
	var account model.Account
	if err := receiver.
		db.
		QueryRowContext(ctx, "select account_number, api_key, bank_code, currency, iban, bic from accounts where account_number = ?", accountNumber).
		Scan(&account.AccountNumber, &account.ApiKey, &account.BankCode, &account.Currency, &account.IBAN, &account.BIC); err != nil {
		return model.Account{}, fmt.Errorf("failed querying account: %w", receiver.wrapError(err))
	}

	return account, nil
}

func (receiver *Manager) RemoveAccountByNumber(ctx context.Context, accountNumber string) error {
	_, err := receiver.db.ExecContext(ctx, "delete from accounts where account_number = ?", accountNumber)
	if err != nil {
		return fmt.Errorf("failed deleting account: %w", receiver.wrapError(err))
	}

	return nil
}

func (receiver *Manager) StoreTransactions(ctx context.Context, transactions []model.Transaction) error {
	values := make([]any, 0)
	query := "insert into transactions (id, account_number, date, amount, currency, counterparty_account, counterparty_name, counterparty_bank_code, counterparty_bank_name, constant_symbol, variable_symbol, specific_symbol, user_identity, transaction_type, performed_by, additional_info, comment, bic, instruction_id, payer_reference) values "
	subqueries := make([]string, 0, len(transactions))
	for _, transaction := range transactions {
		tmpVals := []any{transaction.ID, transaction.AccountNumber, transaction.Date, transaction.Amount, transaction.Currency, transaction.CounterpartyAccount, transaction.CounterpartyName, transaction.CounterpartyBankCode, transaction.CounterpartyBankName, transaction.ConstantSymbol, transaction.VariableSymbol, transaction.SpecificSymbol, transaction.UserIdentity, transaction.TransactionType, transaction.PerformedBy, transaction.AdditionalInfo, transaction.Comment, transaction.BIC, transaction.InstructionID, transaction.PayerReference}
		values = append(values, tmpVals...)
		subqueries = append(subqueries, "("+strings.Join(lo.RepeatBy(len(tmpVals), func(index int) string {
			return "?"
		}), ", ")+")")
	}

	query += strings.Join(subqueries, ", \n")

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
		); err != nil {
			return nil, fmt.Errorf("failed loading transactions: %w", err)
		}

		result = append(result, transaction)
	}

	return result, nil
}
