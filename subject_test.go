package subject

import (
	"fmt"
	"observer/obs"
	"testing"
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

func TestNewSubject(t *testing.T) {
	// give nil, must panic
	sub := New[*obs.Observer[string], string]()
	sub.Pub("25")
	if len(sub.safeObs.observers) != 0 {
		t.Fatal("observer list too big")
	}
}

func TestSubject_Subscribe_Unsubscribe(t *testing.T) {
	sub := New[*obs.Observer[string], string]()
	obs1 := obs.FromFunc(func(t string) {})
	obs2 := obs.FromFunc(func(t string) {})
	obs3 := obs.FromFunc(func(t string) {})
	sub.Subscribe(obs1, obs2, obs3)
	if len(sub.safeObs.observers) != 3 {
		t.Fatal("observer list must be 3")
	}
	sub.Unsubscribe(obs1)
	if len(sub.safeObs.observers) != 2 {
		t.Fatal("observer list must be 2")
	}
	sub.Unsubscribe(obs2)
	sub.Unsubscribe(obs1)
	if len(sub.safeObs.observers) != 1 {
		t.Fatal("observer list must be 1")
	}
	sub.Unsubscribe(obs3)
	if len(sub.safeObs.observers) != 0 {
		t.Fatal("observer list must be empty")
	}
}

func TestSubject_Publish(t *testing.T) {
	sub := New[*obs.Observer[string], string]()
	res := make([]string, 3)
	obs1 := obs.FromFunc(
		func(t string) {
			res[0] = t
		},
	)
	obs2 := obs.FromFunc(
		func(t string) {
			res[1] = t
		},
	)
	obs3 := obs.FromFunc(
		func(t string) {
			res[2] = t
		},
	)
	obss := []*obs.Observer[string]{obs1, obs2, obs3}
	sub.Subscribe(obss...)
	for _, o := range obss {
		o.Listen()
	}
	testMsg := "TEST"
	sub.Pub(testMsg)
	time.Sleep(time.Second * 3)
	for _, val := range res {
		if val != testMsg {
			t.Fatal("the observer must have received the message")
		}
	}
}
