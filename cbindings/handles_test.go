package main

import (
	"errors"
	"testing"
)

type testCloser struct {
	closed bool
	err    error
}

func (value *testCloser) Close() error {
	value.closed = true
	return value.err
}

func TestHandleLifecycle(t *testing.T) {
	value := &testCloser{}
	id := registerHandle(value)

	got, err := getHandle[*testCloser](id)
	if err != nil {
		t.Fatalf("getHandle() error = %v", err)
	}
	if got != value {
		t.Fatalf("getHandle() = %p, want %p", got, value)
	}
	if _, err := getHandle[string](id); err == nil {
		t.Fatal("getHandle() accepted a handle of the wrong type")
	}
	if err := closeHandle(id); err != nil {
		t.Fatalf("closeHandle() error = %v", err)
	}
	if !value.closed {
		t.Error("closeHandle() did not close the object")
	}
	if _, err := getHandle[*testCloser](id); err == nil {
		t.Fatal("getHandle() returned a closed handle")
	}
	if err := closeHandle(id); err == nil {
		t.Fatal("closeHandle() accepted an unknown handle")
	}
}

func TestCloseHandleRemovesObjectWhenCloseFails(t *testing.T) {
	id := registerHandle(&testCloser{err: errors.New("close failed")})
	if err := closeHandle(id); err == nil {
		t.Fatal("closeHandle() did not return the close error")
	}
	if _, err := getHandle[*testCloser](id); err == nil {
		t.Fatal("getHandle() returned a handle after Close failed")
	}
}
