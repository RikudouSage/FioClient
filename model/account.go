package model

import "go.chrastecky.dev/fio-api/fio/dto"

type Account struct {
	AccountNumber string
	ApiKey        string
	BankCode      string
	Currency      string
	IBAN          string
	BIC           string
}

func AccountFromApiModel(account dto.AccountInfo) Account {
	return Account{
		AccountNumber: account.AccountID,
		BankCode:      account.BankID,
		Currency:      account.Currency,
		IBAN:          account.IBAN,
		BIC:           account.BIC,
	}
}
