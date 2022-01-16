package subject

import (
	"sync"
)

type IObserver[T any] interface {
	Do(T)
}

type distributor[T any] interface {
	add(IObserver[T])
	remove(IObserver[T])
	publish(T)
	hasObserver() bool
}

type basicObserver[T any] struct {
	update func(T)
}

func (b *basicObserver[T]) Do(val T) {
	b.update(val)
}

func fromFunc[T any](fn func(T)) IObserver[T] {
	return &basicObserver[T]{
		update: fn,
	}
}

type publisher[T any] struct {
	mutex sync.Mutex
	to    []IObserver[T]
}

func (p *publisher[T]) withMutex(fn func()) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	fn()
}

func (p *publisher[T]) add(observer IObserver[T]) {
	p.withMutex(
		func() {
			p.to = append(p.to, observer)
		},
	)
}

func (p *publisher[T]) remove(observer IObserver[T]) {
	p.withMutex(
		func() {
			idx := -1
			for i, ob := range p.to {
				if ob == observer {
					idx = i
					break
				}
			}
			if idx >= 0 {
				p.to[idx] = p.to[len(p.to)-1]
				p.to[len(p.to)-1] = *new(IObserver[T])
				p.to = p.to[:len(p.to)-1]
			}
		},
	)
}

func (p *publisher[T]) publish(data T) {
	p.withMutex(
		func() {
			for _, observer := range p.to {
				observer.Do(data)
			}
		},
	)
}

func (p *publisher[T]) hasObserver() bool {
	return len(p.to) > 0
}
