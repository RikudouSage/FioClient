package model

import (
	"github.com/shopspring/decimal"
	"go.chrastecky.dev/fio-api/fio/dto"
	. "go.chrastecky.dev/fio-api/fio/types"
)

type Transaction struct {
	ID                   int64
	AccountNumber        string
	Date                 TimezonedDate
	Amount               decimal.Decimal
	Currency             string
	CounterpartyAccount  string
	CounterpartyName     string
	CounterpartyBankCode string
	CounterpartyBankName string
	ConstantSymbol       *string
	VariableSymbol       *string
	SpecificSymbol       *string
	UserIdentity         *string
	TransactionType      TransactionType
	PerformedBy          *string
	AdditionalInfo       *string
	Comment              *string
	BIC                  *string
	InstructionID        *int64
	PayerReference       *string
}

func TransactionFromApiModel(transaction dto.Transaction) Transaction {
	return Transaction{
		ID:                   transaction.ID.Value,
		Date:                 transaction.Date.Value,
		Amount:               transaction.Amount.Value,
		Currency:             transaction.Currency.Value,
		CounterpartyAccount:  transaction.CounterpartyAccount.Value,
		CounterpartyName:     transaction.CounterpartyName.Value,
		CounterpartyBankCode: transaction.CounterpartyBankCode.Value,
		CounterpartyBankName: transaction.CounterpartyBankName.Value,
		ConstantSymbol:       transaction.ConstantSymbol.Value,
		VariableSymbol:       transaction.VariableSymbol.Value,
		SpecificSymbol:       transaction.SpecificSymbol.Value,
		UserIdentity:         transaction.UserIdentity.Value,
		TransactionType:      transaction.TransactionType.Value,
		PerformedBy:          transaction.PerformedBy.Value,
		AdditionalInfo:       transaction.AdditionalInfo.Value,
		Comment:              transaction.Comment.Value,
		BIC:                  transaction.BIC.Value,
		InstructionID:        transaction.InstructionID.Value,
		PayerReference:       transaction.PayerReference.Value,
	}
}
