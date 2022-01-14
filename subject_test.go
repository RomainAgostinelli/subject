package subject

import (
	"fmt"
	"testing"
)

// MyObs is an example implementation of an Observer
type MyObs struct {
	rcv string // rcv for test purpose, stores the message received via Notify in continue into this variable
}

func (m *MyObs) Notify(val string) {
	m.rcv = val
	fmt.Println(val)
}

func ExampleNewSubject() {
	myObs := &MyObs{}
	subject := NewSubject[*MyObs, string]()
	subject.Subscribe(myObs)
	subject.Publish("Test - 1")
	subject.Unsubscribe(myObs)
	subject.Publish("Test - 2")
	// Output: Test - 1
}

func TestNewSubject(t *testing.T) {
	// give nil, must panic
	sub := NewSubject[*MyObs, string]()
	sub.Publish("25")
	if len(sub.obs) != 0 {
		t.Fatal("observer list too big")
	}
}

func TestSubject_Subscribe_Unsubscribe(t *testing.T) {
	sub := NewSubject[*MyObs, string]()
	obs1 := &MyObs{}
	obs2 := &MyObs{}
	obs3 := &MyObs{}
	sub.Subscribe(obs1, obs2, obs3)
	if len(sub.obs) != 3 {
		t.Fatal("observer list must be 3")
	}
	sub.Unsubscribe(obs1)
	if len(sub.obs) != 2 {
		t.Fatal("observer list must be 2")
	}
	sub.Unsubscribe(obs2)
	sub.Unsubscribe(obs1)
	if len(sub.obs) != 1 {
		t.Fatal("observer list must be 1")
	}
	sub.Unsubscribe(obs3)
	if len(sub.obs) != 0 {
		t.Fatal("observer list must be empty")
	}
}

func TestSubject_Publish(t *testing.T) {
	sub := NewSubject[*MyObs, string]()
	obs1 := &MyObs{}
	obs2 := &MyObs{}
	obs3 := &MyObs{}
	obss := []*MyObs{obs1, obs2, obs3}
	sub.Subscribe(obss...)
	testMsg := "TEST"
	sub.Publish(testMsg)
	for _, obs := range obss {
		if obs.rcv != testMsg {
			t.Fatal("the observer must have received the message")
		}
	}
}
