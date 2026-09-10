package main

import (
	"fmt"
	"io"
	"sync"
)

type handle uint64

var (
	handlesMu  sync.RWMutex
	handles           = make(map[handle]any)
	nextHandle handle = 1
)

func registerHandle(value any) handle {
	handlesMu.Lock()
	defer handlesMu.Unlock()

	id := nextHandle
	nextHandle++
	handles[id] = value

	return id
}

func getHandle[T any](id handle) (T, error) {
	handlesMu.RLock()
	defer handlesMu.RUnlock()

	var zero T

	value, ok := handles[id]
	if !ok {
		return zero, fmt.Errorf("handle %d not found", id)
	}

	typed, ok := value.(T)
	if !ok {
		return zero, fmt.Errorf("handle %d has the wrong type", id)
	}

	return typed, nil
}

func closeHandle(id handle) error {
	handlesMu.Lock()
	value, ok := handles[id]
	if ok {
		delete(handles, id)
	}
	handlesMu.Unlock()

	if !ok {
		return fmt.Errorf("handle %d is not registered", id)
	}
	if closer, ok := value.(io.Closer); ok {
		return closer.Close()
	}

	return nil
}
