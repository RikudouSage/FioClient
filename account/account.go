package account

import (
	"context"
	"fmt"
	"time"

	"github.com/samber/lo"
	"go.chrastecky.dev/fio-api/fio"
	"go.chrastecky.dev/fio-api/fio/dto"
	"go.chrastecky.dev/fio-client/fioclient/internal/db"
	"go.chrastecky.dev/fio-client/fioclient/model"
)

type Account interface {
	ResetTransactionPointer(ctx context.Context, to time.Time) error
	LoadNewTransactions(ctx context.Context) ([]model.Transaction, error)

	Transactions(ctx context.Context) ([]model.Transaction, error)
}

type account struct {
	api          fio.Client
	accountModel model.Account
	manager      *db.Manager
}

func New(model model.Account, dbManager *db.Manager) Account {
	return &account{
		api:          lo.Must(fio.NewClient(model.ApiKey)),
		accountModel: model,
		manager:      dbManager,
	}
}

func (receiver *account) ResetTransactionPointer(ctx context.Context, to time.Time) error {
	if to.IsZero() {
		to = time.Now().Add(-90 * 24 * time.Hour)
	}

	return receiver.api.SetLastFailedTransactionDate(ctx, to)
}

func (receiver *account) LoadNewTransactions(ctx context.Context) ([]model.Transaction, error) {
	newTransactionsApi, err := receiver.api.TransactionsSinceLastPull(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed loading new transactions: %w", err)
	}

	newTransactions := lo.Map(newTransactionsApi, func(item dto.Transaction, index int) model.Transaction {
		result := model.TransactionFromApiModel(item)
		result.AccountNumber = receiver.accountModel.AccountNumber

		return result
	})

	if len(newTransactions) == 0 {
		return make([]model.Transaction, 0), nil
	}

	if err := receiver.manager.StoreTransactions(ctx, newTransactions); err != nil {
		if currentTransactions, err := receiver.manager.LoadTransactions(
			ctx,
			db.WithAccountNumber(receiver.accountModel.AccountNumber),
			db.WithLimit(1),
			db.WithOrderBy("date", "desc"),
		); err == nil && len(currentTransactions) > 0 {
			first := currentTransactions[0]
			_ = receiver.api.SetLastFailedTransactionDate(ctx, first.Date.AsTime())
		}

		return nil, fmt.Errorf("failed storing new transactions: %w", err)
	}

	return newTransactions, nil
}

func (receiver *account) Transactions(ctx context.Context) ([]model.Transaction, error) {
	return receiver.manager.LoadTransactions(
		ctx,
		db.WithAccountNumber(receiver.accountModel.AccountNumber),
		db.WithOrderBy("date", "desc"),
	)
}
