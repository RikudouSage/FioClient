# Using Fio Client with SQLCipher from Go

SQLCipher support is optional. It is needed when an application uses
`fioclient.WithDatabase` to keep API tokens inside the encrypted database. An
application that supplies an `EncryptedAPIKeyProvider` may instead use
`WithSQLDatabase` with regular SQLite and does not need the build process below.

`WithDatabase` passes the database key and cipher settings through SQLite
connection parameters. Applications using it must apply the repository's small
patch for `github.com/mattn/go-sqlite3` and link against SQLCipher instead of the
system SQLite library.

## Prepare SQLCipher and the driver

Clone Fio Client, enter its Nix development shell, and build SQLCipher plus the
patched driver copy:

```shell
cd /path/to/FioClient
nix develop
make build-sqlcipher-current generate-go-sqlite3
```

This creates:

- `build/sqlcipher/current/install/lib/libsqlite3.a`
- `build/sqlcipher/current/install/include/`
- `build/go-modules/github.com/mattn/go-sqlite3/`

## Configure the consuming module

From the Go application's directory, point its module graph at the patched
driver copy:

```shell
go mod edit \
  -replace github.com/mattn/go-sqlite3=/path/to/FioClient/build/go-modules/github.com/mattn/go-sqlite3
go mod tidy
```

The replacement path is generated build output and is intentionally local. Do
not commit that absolute `replace` directive when other developers or CI use a
different checkout path; configure an equivalent path in each build
environment instead.

## Build the application

While still inside the Nix development shell, build with CGO enabled and link
the generated SQLCipher archive statically:

```shell
CGO_ENABLED=1 \
CGO_CFLAGS="-I/path/to/FioClient/build/sqlcipher/current/install/include" \
CGO_LDFLAGS="-L/path/to/FioClient/build/sqlcipher/current/install/lib -Wl,-Bstatic -lsqlite3 -Wl,-Bdynamic -L$OPENSSL_CURRENT_LIB -lcrypto -lm" \
go build -tags "libsqlite3 no_postgres no_mysql no_ydb no_clickhouse no_libsql no_mssql no_vertica" ./...
```

The order of the linker flags matters: `-Bstatic` applies only to SQLCipher,
then `-Bdynamic` restores dynamic linking for OpenSSL and `libm`.

The application can now create an encrypted database normally:

```go
client, err := fioclient.New(
	fioclient.WithDatabase(
		"fio.sqlite3",
		types.NewStringSecretKey("secret"),
	),
)
```

For cross-compilation, use the matching SQLCipher install directory, compiler,
and OpenSSL variables from `flake.nix`; the project's `build-lib-386`,
`build-lib-arm7`, and `build-lib-arm64` Make targets are working examples.
