package main

/*
#include "headers/common.h"
*/
import "C"

import "context"

type contextHandle struct {
	context.Context
	cancel context.CancelFunc
}

func (c *contextHandle) Close() error {
	c.cancel()

	return nil
}

//export FioClientNewContext
func FioClientNewContext(out *C.FioClientContextHandle) C.FioClientResult {
	if out == nil {
		setLastError(nullPointerError("out"))
		return C.FioClientFailure
	}

	ctx, cancel := context.WithCancel(context.Background())
	*out = C.FioClientContextHandle(registerHandle(&contextHandle{Context: ctx, cancel: cancel}))

	clearLastError()

	return C.FioClientSuccess
}

//export FioClientCloseHandle
func FioClientCloseHandle(id C.FioClientHandle) C.FioClientResult {
	if err := closeHandle(handle(id)); err != nil {
		setLastError(err)
		return C.FioClientFailure
	}

	clearLastError()

	return C.FioClientSuccess
}
