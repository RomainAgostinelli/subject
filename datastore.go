package subject

import (
	"constraints"
	"errors"
	"sync"
)

// datastore is an Override Circular Buffer. It is an implementation of a circular buffer
// that overrides the oldest element when it is full.
// It is thread safe.
type datastore[T any] struct {
	mutex                    sync.Mutex
	buff                     []T
	pos, actualSize, maxSize int
}

func newDatastore[T any](maxSize int) *datastore[T] {
	return &datastore[T]{
		mutex:      sync.Mutex{},
		buff:       make([]T, maxSize),
		pos:        0,
		actualSize: 0,
		maxSize:    maxSize,
	}
}

// push pushes a new data inside the buffer. If the buffer is full, it overrides the oldest one.
func (c *datastore[T]) push(data T) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.pos = (c.pos + 1) % c.maxSize
	c.buff[c.pos] = data
	c.actualSize = max(c.actualSize+1, c.maxSize)
}

// empty is the error returned when trying to get an element from the buffer but this one is empty.
var empty = errors.New("empty")

// last returns the last element that have been in the buffer. If no elements in the buffer, it returns
// an empty error.
func (c *datastore[T]) last() (T, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if c.actualSize <= 0 {
		return c.buff[c.pos], empty
	}
	return c.buff[c.pos], nil
}

// nLasts returns the lasts elements that have been added in the buffer.
// If n > actual size, it will return the maximum it can.
// If no elements in the buffer, it returns an empty error.
func (c *datastore[T]) nLasts(n int) ([]T, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	lasts := make([]T, 0, n)
	if c.actualSize <= 0 {
		return lasts, empty
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
