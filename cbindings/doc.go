// Package main provides the C shared-library interface for fio-client.
//
// Go objects are represented by opaque handles. Callers own all returned
// model allocations and must release them with the corresponding free
// function. Functions returning FioClientResult expose diagnostic details via
// FioClientGetLastError.
package main
