package model

import "go.chrastecky.dev/fio-api/fio/dto"

// Account describes a Fio account registered with the client.
type Account struct {
	// AccountNumber is the domestic account number.
	AccountNumber string
	// ApiKey is the Fio API token used to access the account.
	ApiKey string
	// BankCode is the domestic bank identifier.
	BankCode string
	// Currency is the account's ISO 4217 currency code.
	Currency string
	// IBAN is the account's international bank account number.
	IBAN string
	// BIC is the account bank's business identifier code.
	BIC string
}

// AccountFromApiModel converts Fio API account information into the locally
// persisted account model. The returned model does not contain an API key.
func AccountFromApiModel(account dto.AccountInfo) Account {
	return Account{
		AccountNumber: account.AccountID,
		BankCode:      account.BankID,
		Currency:      account.Currency,
		IBAN:          account.IBAN,
		BIC:           account.BIC,
	}
}
