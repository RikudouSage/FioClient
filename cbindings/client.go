package main

/*
#include <stdbool.h>
#include "headers/common.h"

typedef struct {
	const char* database_path;
	const char* database_key;
} FioClientOptions;
*/
import "C"

import (
	fioclient "go.chrastecky.dev/fio-client/fioclient"
	"go.chrastecky.dev/fio-client/fioclient/types"
)

func commonHandles(clientID C.FioClientHandle, contextID C.FioClientContextHandle) (fioclient.Client, *contextHandle, error) {
	client, err := getHandle[fioclient.Client](handle(clientID))
	if err != nil {
		return nil, nil, err
	}

	ctx, err := getHandle[*contextHandle](handle(contextID))
	if err != nil {
		return nil, nil, err
	}

	return client, ctx, nil
}

// FioClientNew creates a client backed by an encrypted SQLite database.
// The option strings are borrowed for the duration of this call.
//
//export FioClientNew
func FioClientNew(out *C.FioClientHandle, options C.FioClientOptions) C.FioClientResult {
	if out == nil {
		setLastError(nullPointerError("out"))
		return C.FioClientFailure
	}

	*out = 0

	if options.database_path == nil {
		setLastError(nullPointerError("options.database_path"))
		return C.FioClientFailure
	}

	if options.database_key == nil {
		setLastError(nullPointerError("options.database_key"))
		return C.FioClientFailure
	}

	client, err := fioclient.New(fioclient.WithDatabase(
		C.GoString(options.database_path),
		types.NewStringSecretKey(C.GoString(options.database_key)),
	))
	if err != nil {
		setLastError(err)
		return C.FioClientFailure
	}

	*out = C.FioClientHandle(registerHandle(client))

	clearLastError()

	return C.FioClientSuccess
}
