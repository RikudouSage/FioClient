package fioclient

import (
	"fmt"

	"github.com/pressly/goose/v3"
	"go.chrastecky.dev/fio-client/fioclient/migrations"
)

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
