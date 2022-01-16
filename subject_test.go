package subject

import (
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"
)

func ExampleNew() {
	subject := New[string]()
	subscription := subject.Subscribe(
		func(val string) {
			fmt.Println(val)
		},
	)

	subject.PubAsync(
		func() string {
			time.Sleep(time.Second)
			return "Test - 1"
		},
	)
	subject.Pub("Test - 2")
	time.Sleep(time.Second)
	subscription.Unsubscribe()
	// Output:
	// Test - 2
	// Test - 1
}

func ExampleOf() {
	// just create a simple web service
	go func() {
		http.HandleFunc(
			"/test", func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("TEST"))
			},
		)
		http.ListenAndServe(":64999", nil)
	}()

	// wait a bit for the endpoint to be served
	time.Sleep(time.Second * 2)

	// Demonstration
	subject := Of(
		func() string {
			fmt.Println("LAZY EXEC")
			get, _ := http.Get("http://localhost:64999/test")
			defer get.Body.Close()
			all, _ := io.ReadAll(get.Body)
			return string(all)
		},
	)

	fmt.Println("MAIN 1")
	subject.Subscribe(
		func(val string) {
			fmt.Println(val)
		},
	)
	fmt.Println("MAIN 2")
	time.Sleep(time.Second * 3)
	// Output:
	// MAIN 1
	// MAIN 2
	// LAZY EXEC
	// TEST
}

func TestNewSubject(t *testing.T) {
	// give nil, must panic
	sub := New[string]()
	sub.Pub("25")
	if sub.distributor.hasObserver() {
		t.Fatal("observer list too big")
	}
}

func TestSubject_Subscribe_Unsubscribe(t *testing.T) {
	sub := New[string]()
	subs1 := sub.Subscribe(func(t string) {})
	subs2 := sub.Subscribe(func(t string) {})
	subs3 := sub.Subscribe(func(t string) {})

	if !sub.distributor.hasObserver() {
		t.Fatal("It must have observer")
	}
	subs1.Unsubscribe()
	if !sub.distributor.hasObserver() {
		t.Fatal("It must have observer")
	}
	subs2.Unsubscribe()
	subs3.Unsubscribe()

	if sub.distributor.hasObserver() {
		t.Fatal("Subject must not have observers anymore")
	}
}

func TestSubject_Pub(t *testing.T) {
	sub := New[string]()
	res := make([]string, 3)
	sub.Subscribe(
		func(t string) {
			res[0] = t
		},
	)
	sub.Subscribe(
		func(t string) {
			res[1] = t
		},
	)
	sub.Subscribe(
		func(t string) {
			res[2] = t
		},
	)
	testMsg := "TEST"
	sub.Pub(testMsg)
	time.Sleep(time.Second * 3)
	for _, val := range res {
		if val != testMsg {
			t.Fatal("the observer must have received the message")
		}
	}
}
