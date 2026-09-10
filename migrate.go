package fioclient

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/pressly/goose/v3"
	"go.chrastecky.dev/fio-client/fioclient/migrations"
)

// ErrNotSQLCipher indicates that the configured database driver does not
// support SQLCipher.
var ErrNotSQLCipher = errors.New("database driver does not support SQLCipher")

// ErrDatabaseNotEncrypted indicates that the database connection has no active
// SQLCipher key.
var ErrDatabaseNotEncrypted = errors.New("database connection has no SQLCipher key")

// ErrInvalidDatabaseKey indicates that the database cannot be read using the
// configured SQLCipher key.
var ErrInvalidDatabaseKey = errors.New("database cannot be read using the configured SQLCipher key")

func (receiver *client) migrate() error {
	goose.SetBaseFS(migrations.Assets)
	if err := goose.SetDialect("sqlite"); err != nil {
		return fmt.Errorf("failed setting dialect: %w", err)
	}

	if err := goose.Up(receiver.database, "."); err != nil {
		return fmt.Errorf("failed running migrations: %w", err)
	}

	return nil
}

func (receiver *client) verifyEncryptedSQLCipher() error {
	ctx := context.Background()

	conn, err := receiver.database.Conn(ctx)
	if err != nil {
		return fmt.Errorf("getting database connection: %w", err)
	}
	defer conn.Close()

	var version string
	err = conn.QueryRowContext(ctx, "pragma cipher_version").Scan(&version)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrNotSQLCipher
	}
	if err != nil {
		return fmt.Errorf("checking SQLCipher version: %w", err)
	}
	if version == "" {
		return ErrNotSQLCipher
	}

	var provider string
	err = conn.QueryRowContext(ctx, "pragma cipher_provider").Scan(&provider)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrDatabaseNotEncrypted
	}
	if err != nil {
		return fmt.Errorf("checking SQLCipher codec: %w", err)
	}
	if provider == "" {
		return ErrDatabaseNotEncrypted
	}

	var schemaObjects int
	err = conn.QueryRowContext(
		ctx,
		"select count(*) from sqlite_master",
	).Scan(&schemaObjects)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrInvalidDatabaseKey, err)
	}

	return nil
}
