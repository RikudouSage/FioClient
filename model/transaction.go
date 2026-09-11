package model

import (
	"github.com/shopspring/decimal"
	"go.chrastecky.dev/fio-api/fio/dto"
	. "go.chrastecky.dev/fio-api/fio/types"
)

// Transaction describes a booked transaction persisted by the client.
// Pointer fields are nil when the corresponding value was absent from the Fio
// API response.
type Transaction struct {
	// ID is the transaction identifier assigned by Fio.
	ID int64 `json:"id"`
	// AccountNumber identifies the registered account that owns the transaction.
	AccountNumber string `json:"account_number"`
	// Date is the booking date together with the timezone supplied by Fio.
	Date TimezonedDate `json:"date"`
	// Amount is the signed transaction amount.
	Amount decimal.Decimal `json:"amount"`
	// Currency is the transaction's ISO 4217 currency code.
	Currency string `json:"currency"`
	// CounterpartyAccount is the counterparty's account number.
	CounterpartyAccount string `json:"counterparty_account"`
	// CounterpartyName is the counterparty's reported name.
	CounterpartyName string `json:"counterparty_name"`
	// CounterpartyBankCode is the counterparty bank's domestic identifier.
	CounterpartyBankCode string `json:"counterparty_bank_code"`
	// CounterpartyBankName is the counterparty bank's reported name.
	CounterpartyBankName string `json:"counterparty_bank_name"`
	// ConstantSymbol is the optional Czech constant payment symbol.
	ConstantSymbol *string `json:"constant_symbol"`
	// VariableSymbol is the optional Czech variable payment symbol.
	VariableSymbol *string `json:"variable_symbol"`
	// SpecificSymbol is the optional Czech specific payment symbol.
	SpecificSymbol *string `json:"specific_symbol"`
	// UserIdentity is the optional user-provided transaction identity.
	UserIdentity *string `json:"user_identity"`
	// TransactionType identifies the kind of transaction reported by Fio.
	TransactionType TransactionType `json:"transaction_type"`
	// PerformedBy optionally identifies who performed the transaction.
	PerformedBy *string `json:"performed_by"`
	// AdditionalInfo contains optional additional information from Fio.
	AdditionalInfo *string `json:"additional_info"`
	// Comment contains the optional transaction comment.
	Comment *string `json:"comment"`
	// BIC is the optional counterparty bank business identifier code.
	BIC *string `json:"bic"`
	// InstructionID is the optional Fio payment instruction identifier.
	InstructionID *int64 `json:"instruction_id"`
	// PayerReference is the optional payer-supplied reference.
	PayerReference *string `json:"payer_reference"`
}

// TransactionFromAPIModel converts a Fio API transaction into the locally
// persisted transaction model. AccountNumber must be assigned by the caller.
func TransactionFromAPIModel(transaction dto.Transaction) Transaction {
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
