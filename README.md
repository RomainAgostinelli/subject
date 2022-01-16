# Subjects - Observers

The observer pattern is very useful in a lot of situation, especially when we want to make events driven programs.

A subject is used to distribute the same information to multiple subscribers, here called observers. An observer can
listen on the evolution of the subject and will be notified when an update appears.

## Usage

#### Create an Observer

First, we must create an observer. To do so, you can call the method `FromFunc` in the `obs` package.

```go
myObs := obs.FromFunc[string](
func (val string) {
fmt.Println(val)
},
)
```

#### Create a Subject

Now we can create a simple Subject that manage the kind of observer we declared:

```go
subject := New[*obs.Observer[string], string]()
```

#### Subscribe / Unsubscribe

And finally we can subscribe observers to the subject we created:

```go
obs1 := obs.FromFunc(func (t string) {})
obs2 := obs.FromFunc(func (t string) {})
obs3 := obs.FromFunc(func (t string) {})
sub.Subscribe(obs1, obs2, obs3)
```

And when we don't want to listen on the subject anymore:

```go
sub.Unsubscribe(obs2)
```

#### Publish data

To publish data to all observers via our Subject:

```go
sub.Publish("25")
```

#### Full example

```go
package subject

import (
	"fmt"
	"observer/obs"
	"time"
)

func ExampleNew() {
	myObs := obs.FromFunc[string](
		func(val string) {
			fmt.Println(val)
		},
	)
	subject := New[*obs.Observer[string], string]()
	subject.Subscribe(myObs)
	myObs.Listen()
	subject.Pub("Test - 1")
	subject.Pub("Test - 2")
	time.Sleep(time.Second)
	subject.Unsubscribe(myObs)
	// Output:
	// Test - 1
	// Test - 2
}
```

## Warranty

This implementation is not complete and may have bugs. It has been created to play a bit with the new Go Generics
implementation using Go 1.18beta1 (which I enjoy generally).