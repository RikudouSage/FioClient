package fioclient

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"

	"github.com/pressly/goose/v3"
	"go.chrastecky.dev/fio-client/fioclient/account"
	"go.chrastecky.dev/fio-client/fioclient/internal/db"
	"go.chrastecky.dev/fio-client/fioclient/model"
	. "go.chrastecky.dev/fio-client/fioclient/types"
)

// Client manages registered Fio accounts and their locally persisted
// transaction histories.
//
// A Client owns its database connection. Call Close when the client is no
// longer needed.
type Client interface {
	io.Closer

	// Accounts returns all registered accounts.
	Accounts(ctx context.Context) ([]model.Account, error)
	// Account returns the account-scoped client for accountNumber.
	Account(ctx context.Context, accountNumber string) (account.Account, error)
	// RegisterAccount validates apiKey, stores the corresponding account and
	// its initial transaction history, and returns the stored account metadata.
	// When longAccessToken is true, the initial import requests up to ten years
	// of history; otherwise it requests the preceding 90 days.
	RegisterAccount(ctx context.Context, apiKey string, longAccessToken bool) (model.Account, error)
	// RemoveAccount removes accountNumber and its associated local data.
	RemoveAccount(ctx context.Context, accountNumber string) error
}

type client struct {
	database                *sql.DB
	manager                 *db.Manager
	encryptedAPIKeyProvider EncryptedAPIKeyProvider
	gooseLogger             goose.Logger
}

// New creates a Client using options.
//
// The client must be configured with WithDatabase or WithSQLDatabase. A plain
// SQLite database additionally requires WithEncryptedAPIKeyProvider so API
// tokens are not persisted unencrypted. New applies pending database
// migrations before returning.
func New(options ...Option) (Client, error) {
	instance := &client{}

	options = append([]Option{
		WithGooseLogger(NewNoopGooseLogger()),
	}, options...)
	for _, option := range options {
		if err := option(instance); err != nil {
			return nil, fmt.Errorf("failed applying option: %w", err)
		}
	}

	if instance.database == nil {
		return nil, errors.New("the database path and encryption key must be configured, please use WithDatabase()")
	}

	isSQLCipherErr := instance.verifyEncryptedSQLCipher()
	if isSQLCipherErr != nil {
		switch {
		case errors.Is(isSQLCipherErr, ErrNotSQLCipher), errors.Is(isSQLCipherErr, ErrDatabaseNotEncrypted):
			if instance.encryptedAPIKeyProvider == nil {
				return nil, fmt.Errorf("the sql database is not a properly configured SQLCipher one and no encrypted api key provider is configured - either use SQLCipher or use WithEncryptedAPIKeyProvider(). Reported error: %w", isSQLCipherErr)
			}
		default:
			return nil, fmt.Errorf("failed verifying database encryption: %w", isSQLCipherErr)
		}
	}

	var encryptedProvider EncryptedAPIKeyProvider
	if isSQLCipherErr != nil {
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
