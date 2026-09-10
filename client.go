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
	. "go.chrastecky.dev/fio-client/fioclient/types"
)

type Client interface {
	io.Closer

	Accounts(ctx context.Context) ([]model.Account, error)
	Account(ctx context.Context, accountNumber string) (account.Account, error)
	RegisterAccount(ctx context.Context, apiKey string, longAccessToken bool) (model.Account, error)
	RemoveAccount(ctx context.Context, accountNumber string) error
}

type client struct {
	database                *sql.DB
	manager                 *db.Manager
	encryptedAPIKeyProvider EncryptedAPIKeyProvider
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

	isSqlCipherErr := instance.verifyEncryptedSQLCipher()
	if isSqlCipherErr != nil {
		switch {
		case errors.Is(isSqlCipherErr, ErrNotSQLCipher), errors.Is(isSqlCipherErr, ErrDatabaseNotEncrypted):
			if instance.encryptedAPIKeyProvider == nil {
				return nil, fmt.Errorf("the sql database is not a properly configured SQLCipher one and no encrypted api key provider is configured - either use SQLCipher or use WithEncryptedAPIKeyProvider(). Reported error: %w", isSqlCipherErr)
			}
		default:
			return nil, fmt.Errorf("failed verifying database encryption: %w", isSqlCipherErr)
		}
	}

	var encryptedProvider EncryptedAPIKeyProvider
	if isSqlCipherErr != nil {
		encryptedProvider = instance.encryptedAPIKeyProvider
	}

	if err := instance.migrate(); err != nil {
		return nil, fmt.Errorf("failed migrating database: %w", err)
	}

	instance.manager = db.NewManager(instance.database, encryptedProvider)

	return instance, nil
}

func (receiver *client) Close() error {
	if receiver.database != nil {
		return receiver.database.Close()
	}

	return nil
}
