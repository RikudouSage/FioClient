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

type Option func(instance *client) error

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

func WithSQLDatabase(db *sql.DB) Option {
	return func(instance *client) error {
		instance.database = db
		return nil
	}
}

func WithEncryptedAPIKeyProvider(provider EncryptedAPIKeyProvider) Option {
	return func(instance *client) error {
		instance.encryptedAPIKeyProvider = provider
		return nil
	}
}
