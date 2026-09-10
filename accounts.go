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

var ErrAccountAlreadyExists = errors.New("the account already exists")

func (receiver *client) RegisterAccount(ctx context.Context, apiKey string) (model.Account, error) {
	var zero model.Account

	apiClient, err := fio.NewClient(apiKey)
	if err != nil {
		return zero, fmt.Errorf("failed creating a temporary api client: %w", err)
	}

	apiAccount, apiTransactions, err := apiClient.AccountInfoAndTransactions(ctx)
	if err != nil {
		return zero, fmt.Errorf("failed getting account info: %w", err)
	}

	accountModel := model.AccountFromApiModel(apiAccount)
	accountModel.ApiKey = apiKey

	err = receiver.manager.StoreAccount(accountModel)
	if err != nil {
		if errors.Is(err, db.ErrNonUnique) {
			return zero, fmt.Errorf("%w: %w", ErrAccountAlreadyExists, err)
		}
	}

	err = receiver.manager.StoreTransactions(ctx, lo.Map(apiTransactions, func(item dto.Transaction, _ int) model.Transaction {
		result := model.TransactionFromApiModel(item)
		result.AccountNumber = accountModel.AccountNumber

		return result
	}))
	if err != nil {
		defer func() {
			// best effort
			_ = receiver.manager.RemoveAccountByNumber(ctx, accountModel.AccountNumber)
			_ = apiClient.SetLastFailedTransactionDate(ctx, time.Now().Add(-90*24*time.Hour))
		}()
		return accountModel, fmt.Errorf("failed inserting transactions: %w", err)
	}

	return accountModel, nil
}

func (receiver *client) Account(ctx context.Context, accountNumber string) (account.Account, error) {
	var zero account.Account

	if accountModel, err := receiver.manager.FindAccountByNumber(ctx, accountNumber); err != nil {
		return zero, fmt.Errorf("failed getting account: %w", err)
	} else {
		return account.New(accountModel, receiver.manager), nil
	}
}

func (receiver *client) RemoveAccount(ctx context.Context, accountNumber string) error {
	if err := receiver.manager.RemoveAccountByNumber(ctx, accountNumber); err != nil {
		return fmt.Errorf("failed removing account: %w", err)
	}

	return nil
}
