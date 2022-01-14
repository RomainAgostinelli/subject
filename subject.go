/*
Package subject provides data structure and method to implement simple Subject-Observer pattern.

A subject is used to distribute the same information to multiple subscribers, here called observers.
An observer can listen on the evolution of the subject and will be notified when an update appears.

This implementation is not complete and may have bugs. It has been created to play a bit with the new Go Generics
implementation using Go 1.18beta1 (which I enjoy generally).

*/
package subject

// Observer represents an observer listening of a Subject publishing data of type T.
type Observer[T any] interface {
	// constraints
	comparable
	// Notify accepts a value of type T. It will be called by the Subject when new item is published.
	Notify(val T)
}

// NewSubject creates a new Subject publishing elements of type T to observers of type Observer[T].
func NewSubject[O Observer[T], T any]() *Subject[O, T] {
	s := &Subject[O, T]{obs: make([]O, 0)}
	s.Subscribe()
	return s
}

// Subject represents the subject (or publisher) in the observer pattern. It delivers elements of type T
// to all observers of type Observer
type Subject[O Observer[T], T any] struct {
	obs []O // obs is the list of observers actually observing this Subject.
}

// Publish publishes the data given in parameter to all Observer who are listening on this Subject.
func (s *Subject[O, T]) Publish(data T) {
	for _, ob := range s.obs {
		ob.Notify(data)
	}
}

// Subscribe subscribes a list of Observer to this Subject.
func (s *Subject[O, T]) Subscribe(observer ...O) {
	s.obs = append(s.obs, observer...)
}

// Unsubscribe unsubscribes the observer given in parameter from the observer list (stop listening). The implementation
// of comparable by the Observer will be used to determine if the observer is in the observers list of the Subject.
func (s *Subject[O, T]) Unsubscribe(observer O) {
	idx := -1
	for i, ob := range s.obs {
		if ob == observer {
			idx = i
			break
		}
	}
	if idx >= 0 {
		s.obs[idx] = s.obs[len(s.obs)-1]
		s.obs[len(s.obs)-1] = *new(O)
		s.obs = s.obs[:len(s.obs)-1]
	}
}
