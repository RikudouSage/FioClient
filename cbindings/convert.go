package main

/*
#include <stdlib.h>
#include "headers/models.h"
*/
import "C"

import (
	"errors"
	"fmt"
	"time"
	"unsafe"

	"go.chrastecky.dev/fio-client/fioclient/model"
)

func dateFromC(name string, value *C.char) (time.Time, error) {
	if value == nil {
		return time.Time{}, nullPointerError(name)
	}

	date, err := time.Parse(time.DateOnly, C.GoString(value))
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid %s: %w", name, err)
	}

	return date, nil
}

func cString(value string) *C.char {
	return C.CString(value)
}

func optionalCString(value *string) *C.char {
	if value == nil {
		return nil
	}

	return C.CString(*value)
}

func accountToC(value model.Account) C.FioClientAccount {
	return C.FioClientAccount{
		account_number: cString(value.AccountNumber),
		bank_code:      cString(value.BankCode),
		currency:       cString(value.Currency),
		iban:           cString(value.IBAN),
		bic:            cString(value.BIC),
	}
}

func freeAccount(value *C.FioClientAccount) {
	if value == nil {
		return
	}
	for _, pointer := range []*C.char{value.account_number, value.bank_code, value.currency, value.iban, value.bic} {
		C.free(unsafe.Pointer(pointer))
	}

	*value = C.FioClientAccount{}
}

func transactionToC(value model.Transaction) C.FioClientTransaction {
	result := C.FioClientTransaction{
		id:                     C.int64_t(value.ID),
		account_number:         cString(value.AccountNumber),
		date:                   cString(value.Date.String()),
		amount:                 cString(value.Amount.String()),
		currency:               cString(value.Currency),
		counterparty_account:   cString(value.CounterpartyAccount),
		counterparty_name:      cString(value.CounterpartyName),
		counterparty_bank_code: cString(value.CounterpartyBankCode),
		counterparty_bank_name: cString(value.CounterpartyBankName),
		constant_symbol:        optionalCString(value.ConstantSymbol),
		variable_symbol:        optionalCString(value.VariableSymbol),
		specific_symbol:        optionalCString(value.SpecificSymbol),
		user_identity:          optionalCString(value.UserIdentity),
		transaction_type:       cString(value.TransactionType.String()),
		performed_by:           optionalCString(value.PerformedBy),
		additional_info:        optionalCString(value.AdditionalInfo),
		comment:                optionalCString(value.Comment),
		bic:                    optionalCString(value.BIC),
		payer_reference:        optionalCString(value.PayerReference),
		local_only:             C.bool(value.LocalOnly),
	}

	if value.InstructionID != nil {
		result.instruction_id = (*C.int64_t)(C.malloc(C.size_t(unsafe.Sizeof(C.int64_t(0)))))
		if result.instruction_id != nil {
			*result.instruction_id = C.int64_t(*value.InstructionID)
		}
	}

	return result
}

func freeTransaction(value *C.FioClientTransaction) {
	if value == nil {
		return
	}
	for _, pointer := range []*C.char{value.account_number, value.date, value.amount, value.currency,
		value.counterparty_account, value.counterparty_name, value.counterparty_bank_code, value.counterparty_bank_name,
		value.constant_symbol, value.variable_symbol, value.specific_symbol, value.user_identity, value.transaction_type,
		value.performed_by, value.additional_info, value.comment, value.bic, value.payer_reference} {
		C.free(unsafe.Pointer(pointer))
	}

	C.free(unsafe.Pointer(value.instruction_id))
	*value = C.FioClientTransaction{}
}

func returnTransactions(out *C.FioClientTransactions, values []model.Transaction) C.FioClientResult {
	if out == nil {
		setLastError(nullPointerError("out"))
		return C.FioClientFailure
	}

	*out = C.FioClientTransactions{}

	if len(values) == 0 {
		clearLastError()
		return C.FioClientSuccess
	}

	items := C.calloc(C.size_t(len(values)), C.size_t(unsafe.Sizeof(C.FioClientTransaction{})))
	if items == nil {
		setLastErrorMessage("failed allocating transactions")
		return C.FioClientFailure
	}

	slice := unsafe.Slice((*C.FioClientTransaction)(items), len(values))
	for index := range values {
		slice[index] = transactionToC(values[index])
	}

	out.items = (*C.FioClientTransaction)(items)
	out.length = C.size_t(len(values))

	clearLastError()

	return C.FioClientSuccess
}

func setLastErrorMessage(message string) {
	if message == "" {
		clearLastError()
		return
	}

	setLastError(errors.New(message))
}
