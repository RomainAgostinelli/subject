package obs

// IObserver represents an observer listening of a Subject publishing data of type T.
type IObserver[T any] interface {
	// constraints
	comparable
	// Chan returns a channel on which to listen
	Chan() chan T
	// Listen listens on the channel and execute a function
	Listen()
	// StopListening stops listening on the inner channel
	StopListening()
	// IsListening tells if this observer is currently listening on new events
	IsListening() bool
}

type Observer[T any] struct {
	channel   chan T
	stopChan  chan chan struct{}
	ready     chan struct{}
	listening bool
	Do        func(T)
}

func (o *Observer[T]) Chan() chan T {
	return o.channel
}

// Listen 0
func (o *Observer[T]) Listen() {
	if !o.listening {
		go func() {
			o.listening = true
			defer func() { o.listening = false }()
			// empty the list
			for len(o.Chan()) > 0 {
				<-o.Chan()
			}
			// start listening on new events
			o.ready <- struct{}{}
			for {
				select {
				case finished := <-o.stopChan:
					finished <- struct{}{}
					return
				case val := <-o.Chan():
					o.Do(val)
				}
			}
		}()
		<-o.ready
	}
}

func (o *Observer[T]) IsListening() bool {
	return o.listening
}

func (o *Observer[T]) StopListening() {
	if o.listening {
		finished := make(chan struct{})
		o.stopChan <- finished
		close(finished)
	}
}

func FromFunc[T any](fn func(T)) *Observer[T] {
	return &Observer[T]{
		channel:   make(chan T, 10),
		stopChan:  make(chan chan struct{}),
		ready:     make(chan struct{}),
		listening: false,
		Do:        fn,
	}
}
