package main

/*
#include "headers/errors.h"
*/
import "C"

import (
	"errors"
	"unsafe"
)

func clearLastError() {
	C.fioclient_clear_last_error()
}

func setLastError(err error) {
	if err == nil {
		clearLastError()
		return
	}

	message := C.CString(err.Error())
	defer C.free(unsafe.Pointer(message))

	C.fioclient_set_last_error_copy(message)
}

func nullPointerError(name string) error {
	return errors.New(name + " is NULL")
}

//export FioClientGetLastError
func FioClientGetLastError(buffer *C.char, bufferLength C.size_t) C.size_t {
	return C.fioclient_get_last_error(buffer, bufferLength)
}
