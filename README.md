# Subjects - Observers

The observer pattern is very useful in a lot of situation, especially when we want to make events driven programs.

A subject is used to distribute the same information to multiple subscribers, here called observers. An observer can
listen on the evolution of the subject and will be notified when an update appears.

## Usage

#### Create a Subject

First, we need to create a subject:

```go
subject := New[string]()
```

#### Subscribe / Unsubscribe

Then, we can subscribe on this subject:

```go
subscription := subject.Subscribe(
    func(val string)
	    fmt.Println(val)
	},
)
```

When we subscribe, we get a subscription with which we can unsubscribe:

```go
subscription.Unsubscribe()
```

#### Publish data

For publishing data, there is two options:
* Async: The data will be delivered one day. The asynchronous implements `Lazy Loading`, so it will not execute
the function until there is at least one subscriber.
* Synchronous: The sender will wait for all observers to receive and treats the data before continuing.

##### Asynchronous push
```go
subject.PubAsync(
	func() string {
        time.Sleep(time.Second)
        return "Test - 1"
    },
)
```

The asynchronous push can be useful for `http` requests. The `Async` call is non-blocking.

##### Synchronous push
```go
subject.Pub("Test - 2")
```

#### Full example

```go
package subject

import (
	"fmt"
	"time"
)

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
	// MAIN 2 (or LAZY EXEC)
	// LAZY EXEC (or MAIN 2)
	// TEST
}
```

## Warranty

This implementation is not complete and may have bugs. It has been created to play a bit with the new Go Generics
implementation using Go 1.18beta1 (which I enjoy generally).