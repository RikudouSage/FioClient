package types

type EncryptedAPIKeyProvider interface {
	GetAPIKey(accountNumber string) (string, error)
	StoreAPIKey(accountNumber string, apiKey string) error
	RemoveAPIKey(accountNumber string) error
}
