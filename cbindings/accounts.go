package main

/*
#include <stdbool.h>
#include <stdlib.h>
#include "headers/common.h"
#include "headers/models.h"
*/
import "C"

import "unsafe"

// FioClientGetAccounts returns all registered accounts.
//
//export FioClientGetAccounts
func FioClientGetAccounts(clientID C.FioClientHandle, contextID C.FioClientContextHandle, out *C.FioClientAccounts) C.FioClientResult {
	if out == nil {
		setLastError(nullPointerError("out"))
		return C.FioClientFailure
	}

	*out = C.FioClientAccounts{}

	client, ctx, err := commonHandles(clientID, contextID)
	if err != nil {
		setLastError(err)
		return C.FioClientFailure
	}

	accounts, err := client.Accounts(ctx)
	if err != nil {
		setLastError(err)
		return C.FioClientFailure
	}

	if len(accounts) == 0 {
		clearLastError()
		return C.FioClientSuccess
	}

	items := C.calloc(C.size_t(len(accounts)), C.size_t(unsafe.Sizeof(C.FioClientAccount{})))
	if items == nil {
		setLastErrorMessage("failed allocating accounts")
		return C.FioClientFailure
	}

	slice := unsafe.Slice((*C.FioClientAccount)(items), len(accounts))
	for index := range accounts {
		slice[index] = accountToC(accounts[index])
	}

	out.items = (*C.FioClientAccount)(items)
	out.length = C.size_t(len(accounts))

	clearLastError()

	return C.FioClientSuccess
}

// FioClientRegisterAccount registers an API token and returns account metadata.
//
//export FioClientRegisterAccount
func FioClientRegisterAccount(clientID C.FioClientHandle, contextID C.FioClientContextHandle, apiKey *C.char, longAccessToken C.bool, out *C.FioClientAccount) C.FioClientResult {
	if out == nil {
		setLastError(nullPointerError("out"))
		return C.FioClientFailure
	}

	*out = C.FioClientAccount{}

	if apiKey == nil {
		setLastError(nullPointerError("api_key"))
		return C.FioClientFailure
	}

	client, ctx, err := commonHandles(clientID, contextID)
	if err != nil {
		setLastError(err)
		return C.FioClientFailure
	}

	account, err := client.RegisterAccount(ctx, C.GoString(apiKey), bool(longAccessToken))
	if err != nil {
		setLastError(err)
		return C.FioClientFailure
	}

	*out = accountToC(account)

	clearLastError()

	return C.FioClientSuccess
}

//export FioClientRemoveAccount
func FioClientRemoveAccount(clientID C.FioClientHandle, contextID C.FioClientContextHandle, accountNumber *C.char) C.FioClientResult {
	if accountNumber == nil {
		setLastError(nullPointerError("account_number"))
		return C.FioClientFailure
	}

	client, ctx, err := commonHandles(clientID, contextID)
	if err != nil {
		setLastError(err)
		return C.FioClientFailure
	}

	if err = client.RemoveAccount(ctx, C.GoString(accountNumber)); err != nil {
		setLastError(err)
		return C.FioClientFailure
	}

	clearLastError()

	return C.FioClientSuccess
}

// FioClientOpenAccount creates a handle for operations on one account.
//
//export FioClientOpenAccount
func FioClientOpenAccount(clientID C.FioClientHandle, contextID C.FioClientContextHandle, accountNumber *C.char, out *C.FioClientAccountHandle) C.FioClientResult {
	if out == nil {
		setLastError(nullPointerError("out"))
		return C.FioClientFailure
	}

	*out = 0

	if accountNumber == nil {
		setLastError(nullPointerError("account_number"))
		return C.FioClientFailure
	}

	client, ctx, err := commonHandles(clientID, contextID)
	if err != nil {
		setLastError(err)
		return C.FioClientFailure
	}

	account, err := client.Account(ctx, C.GoString(accountNumber))
	if err != nil {
		setLastError(err)
		return C.FioClientFailure
	}

	*out = C.FioClientAccountHandle(registerHandle(account))

	clearLastError()

	return C.FioClientSuccess
}

//export FioClientFreeAccount
func FioClientFreeAccount(account *C.FioClientAccount) {
	freeAccount(account)
}

//export FioClientFreeAccounts
func FioClientFreeAccounts(accounts *C.FioClientAccounts) {
	if accounts == nil {
		return
	}

	if accounts.items != nil {
		items := unsafe.Slice(accounts.items, int(accounts.length))
		for index := range items {
			freeAccount(&items[index])
		}
		C.free(unsafe.Pointer(accounts.items))
	}

	*accounts = C.FioClientAccounts{}
}
