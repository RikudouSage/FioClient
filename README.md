# Fio Client

Fio Client stores Fio bank accounts and their transactions in a local SQLite
database. It can synchronize new transactions through the
[Fio API client](https://github.com/RikudouSage/FioAPI) and exposes both a Go API
and C bindings.

## Usage

```go
package main

import (
	"context"
	"log"

	"go.chrastecky.dev/fio-client/fioclient"
	"go.chrastecky.dev/fio-client/fioclient/types"
)

func main() {
	client, err := fioclient.New(
		fioclient.WithDatabase("fio.sqlite3", types.NewStringSecretKey("secret")),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	ctx := context.Background()
	account, err := client.RegisterAccount(ctx, "fio-api-token", false)
	if err != nil {
		log.Fatal(err)
	}

	accountClient, err := client.Account(ctx, account.AccountNumber)
	if err != nil {
		log.Fatal(err)
	}

	transactions, err := accountClient.LoadNewTransactions(ctx)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("downloaded %d transactions", len(transactions))
}
```

## Protecting API tokens

Fio API tokens provide access to bank-account data and must be encrypted at
rest. The client supports two storage modes:

- `WithDatabase` opens an encrypted SQLCipher database and stores API tokens
  inside it. The application must be built with SQLCipher support.
- `WithSQLDatabase` accepts an existing database connection. If that database
  is not protected by SQLCipher, `WithEncryptedAPIKeyProvider` is required. The
  provider should store tokens in a platform credential store; only account and
  transaction data remain in SQLite.

The client rejects an unencrypted database unless an external encrypted API-key
provider is configured.

## Building with SQLCipher support

Enter the Nix development shell and build the shared library for the current
system:

```shell
nix develop
make
```

Cross-build aliases are available for its supported
architectures:

```shell
make build-lib-armv7hl
make build-lib-aarch64
make build-lib-i486
```

`libsqlite3.a` from SQLCipher is linked statically into `libfioclient.so`;
OpenSSL and the operating-system libraries remain dynamic dependencies. Use
`make release-all` to place all release artifacts under `out/`.

Go applications only need the patched SQLite driver and SQLCipher linker
configuration when they use `WithDatabase` to store API tokens in the encrypted
database. Applications that provide an `EncryptedAPIKeyProvider` can use a
regular SQLite build instead. See
[Using Fio Client with SQLCipher from Go](docs/go-sqlcipher.md).
