package main

/*
#include <stdint.h>
#include <stdlib.h>
#include "headers/common.h"
#include "headers/models.h"
*/
import "C"

import (
	"fmt"
	"time"
	"unsafe"

	"go.chrastecky.dev/fio-client/fioclient/account"
)

func accountHandles(accountID C.FioClientAccountHandle, contextID C.FioClientContextHandle) (account.Account, *contextHandle, error) {
	accountValue, err := getHandle[account.Account](handle(accountID))
	if err != nil {
		return nil, nil, err
	}

	ctx, err := getHandle[*contextHandle](handle(contextID))
	if err != nil {
		return nil, nil, err
	}

	return accountValue, ctx, nil
}

//export FioClientGetTransactions
func FioClientGetTransactions(accountID C.FioClientAccountHandle, contextID C.FioClientContextHandle, out *C.FioClientTransactions) C.FioClientResult {
	if out == nil {
		setLastError(nullPointerError("out"))
		return C.FioClientFailure
	}

	*out = C.FioClientTransactions{}

	accountValue, ctx, err := accountHandles(accountID, contextID)
	if err != nil {
		setLastError(err)
		return C.FioClientFailure
	}

	transactions, err := accountValue.Transactions(ctx)
	if err != nil {
		setLastError(err)
		return C.FioClientFailure
	}

	return returnTransactions(out, transactions)
}

//export FioClientLoadNewTransactions
func FioClientLoadNewTransactions(accountID C.FioClientAccountHandle, contextID C.FioClientContextHandle, out *C.FioClientTransactions) C.FioClientResult {
	if out == nil {
		setLastError(nullPointerError("out"))
		return C.FioClientFailure
	}

	*out = C.FioClientTransactions{}

	accountValue, ctx, err := accountHandles(accountID, contextID)
	if err != nil {
		setLastError(err)
		return C.FioClientFailure
	}

	transactions, err := accountValue.LoadNewTransactions(ctx)
	if err != nil {
		setLastError(err)
		return C.FioClientFailure
	}

	return returnTransactions(out, transactions)
}

// FioClientLoadTransactionsByDate downloads and stores transactions in the
// inclusive date range. Dates use YYYY-MM-DD format. The caller owns the
// returned allocation and must release it with FioClientFreeTransactions.
//
//export FioClientLoadTransactionsByDate
func FioClientLoadTransactionsByDate(accountID C.FioClientAccountHandle, contextID C.FioClientContextHandle, startDate *C.char, endDate *C.char, out *C.FioClientTransactions) C.FioClientResult {
	if out == nil {
		setLastError(nullPointerError("out"))
		return C.FioClientFailure
	}

	*out = C.FioClientTransactions{}

	accountValue, ctx, err := accountHandles(accountID, contextID)
	if err != nil {
		setLastError(err)
		return C.FioClientFailure
	}

	start, err := dateFromC("start_date", startDate)
	if err != nil {
		setLastError(err)
		return C.FioClientFailure
	}

	end, err := dateFromC("end_date", endDate)
	if err != nil {
		setLastError(err)
		return C.FioClientFailure
	}

	transactions, err := accountValue.LoadTransactionsByDate(ctx, start, end)
	if err != nil {
		setLastError(err)
		return C.FioClientFailure
	}

	return returnTransactions(out, transactions)
}

//export FioClientGetTransaction
func FioClientGetTransaction(accountID C.FioClientAccountHandle, contextID C.FioClientContextHandle, id C.int64_t, out *C.FioClientTransaction) C.FioClientResult {
	if out == nil {
		setLastError(nullPointerError("out"))
		return C.FioClientFailure
	}

	*out = C.FioClientTransaction{}

	accountValue, ctx, err := accountHandles(accountID, contextID)
	if err != nil {
		setLastError(err)
		return C.FioClientFailure
	}

	transaction, err := accountValue.Transaction(ctx, int64(id))
	if err != nil {
		setLastError(err)
		return C.FioClientFailure
	}

	*out = transactionToC(transaction)

	clearLastError()

	return C.FioClientSuccess
}

// FioClientResetTransactionPointer resets the remote marker. An empty value
// uses the package default; otherwise the timestamp must use RFC3339.
//
//export FioClientResetTransactionPointer
func FioClientResetTransactionPointer(accountID C.FioClientAccountHandle, contextID C.FioClientContextHandle, timestamp *C.char) C.FioClientResult {
	accountValue, ctx, err := accountHandles(accountID, contextID)
	if err != nil {
		setLastError(err)
		return C.FioClientFailure
	}

	var parsed time.Time

	if timestamp != nil && C.GoString(timestamp) != "" {
		parsed, err = time.Parse(time.RFC3339, C.GoString(timestamp))
		if err != nil {
			setLastError(fmt.Errorf("invalid timestamp: %w", err))
			return C.FioClientFailure
		}
	}

	if err = accountValue.ResetTransactionPointer(ctx, parsed); err != nil {
		setLastError(err)
		return C.FioClientFailure
	}

	clearLastError()

	return C.FioClientSuccess
}

//export FioClientFreeTransaction
func FioClientFreeTransaction(transaction *C.FioClientTransaction) {
	freeTransaction(transaction)
}

//export FioClientFreeTransactions
func FioClientFreeTransactions(transactions *C.FioClientTransactions) {
	if transactions == nil {
		return
	}

	if transactions.items != nil {
		items := unsafe.Slice(transactions.items, int(transactions.length))
		for index := range items {
			freeTransaction(&items[index])
		}
		C.free(unsafe.Pointer(transactions.items))
	}

	*transactions = C.FioClientTransactions{}
}
