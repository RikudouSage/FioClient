package fioclient

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/samber/lo"
	"go.chrastecky.dev/fio-api/fio"
	"go.chrastecky.dev/fio-api/fio/dto"
	"go.chrastecky.dev/fio-client/fioclient/account"
	"go.chrastecky.dev/fio-client/fioclient/internal/db"
	"go.chrastecky.dev/fio-client/fioclient/model"
)

// ErrAccountAlreadyExists indicates that an account with the same account
// number is already registered.
var ErrAccountAlreadyExists = errors.New("the account already exists")

func (receiver *client) RegisterAccount(ctx context.Context, apiKey string, longAccessToken bool) (model.Account, error) {
	var zero model.Account

	apiClient, err := fio.NewClient(apiKey)
	if err != nil {
		return zero, fmt.Errorf("failed creating a temporary api client: %w", err)
	}

	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -90)
	if longAccessToken {
		startDate = endDate.AddDate(-10, 0, 0)
	}

	apiAccount, apiTransactions, err := apiClient.AccountInfoAndTransactionsByDate(
		ctx,
		startDate,
		endDate,
	)
	if err != nil {
		return zero, fmt.Errorf("failed getting account info: %w", err)
	}

	accountModel := model.AccountFromAPIModel(apiAccount)
	accountModel.APIKey = apiKey

	err = receiver.manager.StoreAccount(ctx, accountModel)
	if err != nil {
		if errors.Is(err, db.ErrNonUnique) {
			return zero, fmt.Errorf("%w: %w", ErrAccountAlreadyExists, err)
		}

		return zero, fmt.Errorf("failed storing account %s: %w", accountModel.AccountNumber, err)
	}

	err = receiver.manager.StoreTransactions(ctx, lo.Map(apiTransactions, func(item dto.Transaction, _ int) model.Transaction {
		result := model.TransactionFromAPIModel(item)
		result.AccountNumber = accountModel.AccountNumber

		return result
	}))
	if err != nil {
		defer func() {
			// best effort
			_ = receiver.manager.RemoveAccountByNumber(ctx, accountModel.AccountNumber)
		}()
		return accountModel, fmt.Errorf("failed inserting transactions: %w", err)
	}

	if len(apiTransactions) > 0 {
		lastID := lo.MaxBy(apiTransactions, func(a, b dto.Transaction) bool {
			return a.ID.Value > b.ID.Value
		}).ID.Value

		if err := apiClient.SetLastTransactionID(ctx, lastID); err != nil {
			return accountModel, fmt.Errorf(
				"failed setting transaction marker for account %s: %w",
				accountModel.AccountNumber,
				err,
			)
		}
	} else {
		if err := apiClient.SetLastFailedTransactionDate(ctx, time.Now()); err != nil {
			return accountModel, fmt.Errorf(
				"failed setting transaction marker for account %s: %w",
				accountModel.AccountNumber,
				err,
			)
		}
	}

	return accountModel, nil
}

func (receiver *client) Account(ctx context.Context, accountNumber string) (account.Account, error) {
	var zero account.Account

	if accountModel, err := receiver.manager.FindAccountByNumber(ctx, accountNumber); err != nil {
		return zero, fmt.Errorf("failed getting account: %w", err)
	} else {
		return account.New(accountModel, receiver.manager, receiver.partialTransactionsEnabled), nil
	}
}

func (receiver *client) RemoveAccount(ctx context.Context, accountNumber string) error {
	if err := receiver.manager.RemoveAccountByNumber(ctx, accountNumber); err != nil {
		return fmt.Errorf("failed removing account: %w", err)
	}

	return nil
}

func (receiver *client) Accounts(ctx context.Context) ([]model.Account, error) {
	return receiver.manager.GetAccounts(ctx)
}
