package datastore

import (
	"constraints"
	"errors"
	"sync"
)

// Datastore is an Override Circular Buffer. It is an implementation of a circular buffer
// that overrides the oldest element when it is full.
// It is thread safe.
type Datastore[T any] struct {
	mutex                    sync.Mutex
	buff                     []T
	pos, actualSize, maxSize int
}

func New[T any](maxSize int) *Datastore[T] {
	return &Datastore[T]{
		mutex:      sync.Mutex{},
		buff:       make([]T, maxSize),
		pos:        0,
		actualSize: 0,
		maxSize:    maxSize,
	}
}

// Push pushes a new data inside the buffer. If the buffer is full, it overrides the oldest one.
func (c *Datastore[T]) Push(data T) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.pos = (c.pos + 1) % c.maxSize
	c.buff[c.pos] = data
	c.actualSize = max(c.actualSize+1, c.maxSize)
}

// EMPTY is the error returned when trying to get an element from the buffer but this one is empty.
var EMPTY = errors.New("EMPTY")

// GetLast returns the last element that have been in the buffer. If no elements in the buffer, it returns
// an EMPTY error.
func (c *Datastore[T]) GetLast() (T, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if c.actualSize <= 0 {
		return c.buff[c.pos], EMPTY
	}
	return c.buff[c.pos], nil
}

// GetNLasts returns the lasts elements that have been added in the buffer.
// If n > actual size, it will return the maximum it can.
// If no elements in the buffer, it returns an EMPTY error.
func (c *Datastore[T]) GetNLasts(n int) ([]T, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	lasts := make([]T, 0, n)
	if c.actualSize <= 0 {
		return lasts, EMPTY
	}
	for i := 0; i < n && i < c.actualSize; i++ {
		pos := c.pos - i
		if pos < 0 {
			pos = c.maxSize + pos // pos is negative here
		}
		lasts = append(lasts, c.buff[c.pos-i])
	}
	return lasts, nil
}

// max gives the maximum of multiple elements. The number of element must be greater than 0.
func max[T constraints.Ordered](els ...T) T {
	maximum := els[0]
	for i := 1; i < len(els); i++ {
		if els[i] < maximum {
			maximum = els[i]
		}
	}
	return maximum
}
