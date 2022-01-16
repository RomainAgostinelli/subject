/*
Package subject provides data structure and method to implement simple Subject-Observer pattern.

A subject is used to distribute the same information to multiple subscribers, here called observers.
An observer can listen on the evolution of the subject and will be notified when an update appears.

This implementation is not complete and may have bugs. It has been created to play a bit with the new Go Generics
implementation using Go 1.18beta1 (which I enjoy generally).

*/
package subject

import (
	"observer/datastore"
	"observer/obs"
	"sync"
)

// New creates a new Subject publishing elements of type T to observers of type Observer[T].
func New[O obs.IObserver[T], T any]() *Subject[O, T] {
	safeObser := &safeObservers[O, T]{
		mutex:     sync.Mutex{},
		observers: make([]O, 0, 1),
	}
	s := &Subject[O, T]{
		safeObs:     safeObser,
		buff:        make(chan T, 100), // do not overheat, so 100 seems good, and for memory it is acceptable
		datastore:   datastore.New[T](10),
		distributor: &publisher[O, T]{},
	}
	s.distributor.To(safeObser)
	s.distributor.From(s.buff)
	s.distributor.Start()
	return s
}

type safeObservers[O obs.IObserver[T], T any] struct {
	mutex     sync.Mutex
	observers []O // observers is a list of observers.
}

// Subject represents the subject (or publisher) in the observer pattern. It delivers elements of type T
// to all observers of type Observer
type Subject[O obs.IObserver[T], T any] struct {
	safeObs     *safeObservers[O, T]
	buff        chan T                  // buff is the buffer of incoming events that are not already distributed.
	datastore   *datastore.Datastore[T] // datastore is used to store already shared events.
	distributor Distributor[O, T]       // distributor is an element which manage the distribution.
	stop        chan chan struct{}      // stop is the channel indicating the managing goroutine to stop.
}

type Distributor[O obs.IObserver[T], T any] interface {
	Start()
	Stop()
	From(chan T)
	To(observers *safeObservers[O, T])
}

type publisher[O obs.IObserver[T], T any] struct {
	stop         chan chan struct{}
	distributing bool
	from         chan T
	to           *safeObservers[O, T]
}

func (p *publisher[O, T]) From(from chan T) {
	p.from = from
}

func (p *publisher[O, T]) To(obs *safeObservers[O, T]) {
	p.to = obs
}

func (p *publisher[O, T]) Start() {
	if !p.distributing {
		go func() {
			for {
				select {
				case val := <-p.from:
					p.to.mutex.Lock()
					for _, o := range p.to.observers {
						if o.IsListening() {
							o.Chan() <- val
						}
					}
					p.to.mutex.Unlock()
				case finished := <-p.stop:
					finished <- struct{}{}
					return
				}
			}
		}()
	}
}

func (p *publisher[O, T]) Stop() {
	if p.distributing {
		finished := make(chan struct{})
		p.stop <- finished
		<-finished
		close(finished)
	}
}

// Pub publishes the data given in parameter to all Observer who are listening on this Subject.
func (s *Subject[O, T]) Pub(data T) {
	s.buff <- data
}

// Subscribe subscribes a list of Observer to this Subject.
func (s *Subject[O, T]) Subscribe(observer ...O) {
	s.safeObs.mutex.Lock()
	defer s.safeObs.mutex.Unlock()
	s.safeObs.observers = append(s.safeObs.observers, observer...)
}

// Unsubscribe unsubscribes the observer given in parameter from the observer list (stop listening). The implementation
// of comparable by the Observer will be used to determine if the observer is in the observers list of the Subject.
func (s *Subject[O, T]) Unsubscribe(observer O) {
	s.safeObs.mutex.Lock()
	defer s.safeObs.mutex.Unlock()
	idx := -1
	for i, ob := range s.safeObs.observers {
		if ob == observer {
			idx = i
			break
		}
	}
	if idx >= 0 {
		s.safeObs.observers[idx] = s.safeObs.observers[len(s.safeObs.observers)-1]
		s.safeObs.observers[len(s.safeObs.observers)-1] = *new(O)
		s.safeObs.observers = s.safeObs.observers[:len(s.safeObs.observers)-1]
	}
}
