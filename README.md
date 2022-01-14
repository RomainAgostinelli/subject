# Subjects - Observers
The observer pattern is very useful in a lot of situation, especially when we want
to make events driven programs. 

A subject is used to distribute the same information to multiple subscribers, here called observers.
An observer can listen on the evolution of the subject and will be notified when an update appears.

## Usage

#### Create an Observer
First, we must create an implementation of an Observer (no need of subject if nobody's here to observe).
To do so, you must implement the interface "Observer" defined in the package. An observer
can observe on type of data. Below a simple observer implementation capable to listen on subject publishing
strings:
```go
type MyObs struct {
}
func (m *MyObs) Notify(val string) {
	fmt.Println(val)
}
```

#### Create a Subject
Now we can create a simple Subject that manage the kind of observer we declared:
```go
subject := NewSubject[*MyObs, string]()
```

#### Subscribe / Unsubscribe
And finally we can subscribe observers to the subject we created:
```go
obs1 := &MyObs{}
obs2 := &MyObs{}
obs3 := &MyObs{}
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
// MyObs is an example implementation of an Observer
type MyObs struct {
}

func (m *MyObs) Notify(val string) {
	fmt.Println(val)
}

func ExampleNewSubject() {
	myObs := &MyObs{}
	subject := NewSubject[*MyObs, string]()
	subject.Subscribe(myObs)
	subject.Publish("Test - 1")
	subject.Unsubscribe(myObs)
}
```

## Warranty
This implementation is not complete and may have bugs. It has been created to play a bit with the new Go Generics
implementation using Go 1.18beta1 (which I enjoy generally).