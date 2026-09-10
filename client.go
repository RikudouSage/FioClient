package fioclient

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"

	"go.chrastecky.dev/fio-client/fioclient/account"
	"go.chrastecky.dev/fio-client/fioclient/internal/db"
	"go.chrastecky.dev/fio-client/fioclient/model"
)

type Client interface {
	io.Closer

	Account(ctx context.Context, accountNumber string) (account.Account, error)
	RegisterAccount(ctx context.Context, apiKey string) (model.Account, error)
	RemoveAccount(ctx context.Context, accountNumber string) error
}

type client struct {
	database *sql.DB
	manager  *db.Manager
}

func New(options ...Option) (Client, error) {
	instance := &client{}

	for _, option := range options {
		if err := option(instance); err != nil {
			return nil, fmt.Errorf("failed applying option: %w", err)
		}
	}

	if instance.database == nil {
		return nil, errors.New("the database path and encryption key must be configured, please use WithDatabase()")
	}

	if err := instance.migrate(); err != nil {
		return nil, fmt.Errorf("failed migrating database: %w", err)
	}

	instance.manager = db.NewManager(instance.database)

	return instance, nil
}

func (receiver *client) Close() error {
	if receiver.database != nil {
		return receiver.database.Close()
	}

	return nil
}
