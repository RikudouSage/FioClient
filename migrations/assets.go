package migrations

import "embed"

// Assets contains the SQL migration files embedded into the client binary.
//
//go:embed *.sql
var Assets embed.FS
