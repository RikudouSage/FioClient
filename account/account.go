package account

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/samber/lo"
	"go.chrastecky.dev/fio-api/fio"
	"go.chrastecky.dev/fio-api/fio/dto"
	"go.chrastecky.dev/fio-client/fioclient/internal/db"
	"go.chrastecky.dev/fio-client/fioclient/model"
)

// ErrTransactionNotFound indicates that an account has no locally stored
// transaction with the requested ID.
var ErrTransactionNotFound = errors.New("transaction not found")

// Account provides synchronization and transaction queries for one registered
// Fio account.
type Account interface {
	// ResetTransactionPointer moves the remote transaction-download marker to
	// the supplied time. A zero time resets it to 90 days before the call.
	ResetTransactionPointer(ctx context.Context, to time.Time) error
	// LoadNewTransactions downloads transactions since the remote marker, stores
	// them locally, and returns the downloaded transactions.
	LoadNewTransactions(ctx context.Context) ([]model.Transaction, error)

	// Transactions returns all locally stored transactions for the account in
	// descending date order.
	Transactions(ctx context.Context) ([]model.Transaction, error)
	// Transaction returns the locally stored transaction identified by
	// transactionID. It returns ErrTransactionNotFound when no match exists.
	Transaction(ctx context.Context, transactionID int64) (model.Transaction, error)
}

type account struct {
	api          fio.Client
	accountModel model.Account
	manager      *db.Manager
}

// New creates an account-scoped client for model using dbManager as its local
// transaction store.
//
// Most callers should use fioclient.Client.Account instead. New is primarily
// intended for integration with the parent client package.
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

func (receiver *account) Transaction(ctx context.Context, transactionID int64) (model.Transaction, error) {
	transactions, err := receiver.manager.LoadTransactions(
		ctx,
		db.WithAccountNumber(receiver.accountModel.AccountNumber),
		db.WithOrderBy("date", "desc"),
		db.WithWhere("id = ?", transactionID),
		db.WithLimit(1),
	)
	if err != nil {
		return model.Transaction{}, fmt.Errorf("failed loading transaction: %w", err)
	}

	if len(transactions) == 0 {
		return model.Transaction{}, ErrTransactionNotFound
	}

	return transactions[0], nil
}
