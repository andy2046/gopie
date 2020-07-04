package latest

import "fmt"

func ExampleNew() {
	receiver := make(chan interface{})
	sender := New(receiver)
	defer close(sender)

	sender <- 1

	fmt.Println("got", <-receiver)

	sender <- 2
	sender <- 3

	fmt.Println("got", <-receiver)
	fmt.Println("got", <-receiver)

	sender <- 4
	sender <- 5

	fmt.Println("got", <-receiver)
	fmt.Println("got", <-receiver)

	// Output:
	// got 1
	// got 3
	// got 3
	// got 5
	// got 5
}

func ExampleNewN() {
	receiver := make(chan []interface{})
	n := 3
	sender := NewN(receiver, n)
	defer close(sender)

	sender <- 11

	fmt.Println("got", <-receiver)

	sender <- 12
	sender <- 13

	fmt.Println("got", <-receiver)
	fmt.Println("got", <-receiver)

	sender <- 14
	sender <- 15

	fmt.Println("got", <-receiver)
	fmt.Println("got", <-receiver)

	// Output:
	// got [<nil> <nil> 11]
	// got [11 12 13]
	// got [11 12 13]
	// got [13 14 15]
	// got [13 14 15]
}
