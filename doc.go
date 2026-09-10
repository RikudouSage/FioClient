// Package fioclient manages Fio bank accounts and their transaction history in
// a local SQLite database.
//
// A Client is configured with an encrypted SQLCipher database or with a custom
// SQL database and an EncryptedAPIKeyProvider. RegisterAccount validates an API
// token, stores the account metadata and initial transaction history, and
// returns an account model. Account returns an account-scoped client used to
// synchronize and query transactions.
package fioclient
