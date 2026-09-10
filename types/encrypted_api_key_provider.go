package types

// EncryptedAPIKeyProvider stores API tokens outside an unencrypted SQL
// database. Implementations should encrypt tokens at rest using a platform
// credential or secret-storage facility.
type EncryptedAPIKeyProvider interface {
	// GetAPIKey returns the API token stored for accountNumber.
	GetAPIKey(accountNumber string) (string, error)
	// StoreAPIKey stores apiKey for accountNumber.
	StoreAPIKey(accountNumber string, apiKey string) error
	// RemoveAPIKey removes the API token stored for accountNumber.
	RemoveAPIKey(accountNumber string) error
}
