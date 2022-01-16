/*
Package subject provides data structure and method to implement simple Subject-Observer pattern.

A subject is used to distribute the same information to multiple subscribers, here called observers.
An observer can listen on the evolution of the subject and will be notified when an update appears.

This implementation is not complete and may have bugs. It has been created to play a bit with the new Go Generics
implementation using Go 1.18beta1 (which I enjoy generally).

*/
package subject

import (
	"sync"
)

// New creates a new Subject publishing elements of type T to observers of type Observer[T].
func New[T any]() *Subject[T] {
	s := &Subject[T]{
		lazyTasks: make([]func() T, 0, 5),
		datastore: newDatastore[T](10),
		distributor: &publisher[T]{
			mutex: sync.Mutex{},
		},
	}
	return s
}

func Of[T any](fn func() T) *Subject[T] {
	s := New[T]()
	s.PubAsync(fn)
	return s
}

// Subject represents the subject (or publisher) in the observer pattern. It delivers elements of type T
// to all observers of type Observer
type Subject[T any] struct {
	lazyTasks   []func() T     // buff is the buffer of incoming events that are not already distributed.
	datastore   *datastore[T]  // datastore is used to store already shared events.
	distributor distributor[T] // distributor is an element which manage the distrib.
}

// Pub publishes the data given in parameter to all Observer who are listening on this Subject.
func (s *Subject[T]) Pub(data T) {
	s.distributor.publish(data)
}

// PubAsync implements lazy loading, it means that it will not execute the task until there is at minima 1 subscriber.
func (s *Subject[T]) PubAsync(fn func() T) *Subject[T] {
	if !s.distributor.hasObserver() {
		// Add to tasks
		s.lazyTasks = append(s.lazyTasks, fn)
		return s
	}
	// make the publishing in another goroutine
	go func() {
		val := fn()
		s.Pub(val)
	}()
	return s
}

// Subscribe subscribes a list of Observer to this Subject.
// create a new subscription
func (s *Subject[T]) Subscribe(fn func(val T)) Subscription[T] {
	sub := Subscription[T]{
		obs:   fromFunc(fn),
		unsub: s.distributor.remove,
	}
	s.distributor.add(sub.obs)

	// if contains lazy tasks, do them now
	for len(s.lazyTasks) > 0 {
		var task func() T
		task, s.lazyTasks = s.lazyTasks[0], s.lazyTasks[1:]
		s.PubAsync(task)
	}
	return sub
}

type Subscription[T any] struct {
	obs   IObserver[T]
	unsub func(observer IObserver[T])
}

func (s *Subscription[T]) Unsubscribe() {
	s.unsub(s.obs)
}
