package fioclient

import (
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	. "go.chrastecky.dev/fio-client/fioclient/types"
)

// Option configures a Client created by New.
type Option func(instance *client) error

// WithDatabase configures a SQLCipher-backed SQLite database at path using
// key. The option enables foreign-key enforcement and SQLCipher compatibility
// version 4 with a 4096-byte cipher page size.
func WithDatabase(path string, key SecretKey) Option {
	sqliteDsn := func(path string, secretKey SecretKey) string {
		separator := "?"
		if strings.Contains(path, "?") {
			separator = "&"
		}

		values := url.Values{}
		values.Set("_pragma_key", secretKey.String())
		values.Set("_pragma_cipher_compatibility", "4")
		values.Set("_pragma_cipher_page_size", "4096")
		values.Set("_pragma_foreign_keys", "on")

		return path + separator + values.Encode()
	}

	return func(instance *client) error {
		if key.IsZero() {
			return errors.New("secret key is empty")
		}
		db, err := sql.Open("sqlite3", sqliteDsn(path, key))
		if err != nil {
			return fmt.Errorf("failed opening database: %w", err)
		}

		instance.database = db
		return nil
	}
}

// WithSQLDatabase configures an existing SQL database connection.
//
// When db is not backed by an encrypted SQLCipher database, callers must also
// provide WithEncryptedAPIKeyProvider to keep API tokens out of the database.
// The Client takes ownership of db and closes it when Client.Close is called.
func WithSQLDatabase(db *sql.DB) Option {
	return func(instance *client) error {
		instance.database = db
		return nil
	}
}

// WithEncryptedAPIKeyProvider stores API tokens outside the SQL database by
// using provider. It is required when the configured database is not protected
// by SQLCipher.
func WithEncryptedAPIKeyProvider(provider EncryptedAPIKeyProvider) Option {
	return func(instance *client) error {
		instance.encryptedAPIKeyProvider = provider
		return nil
	}
}
