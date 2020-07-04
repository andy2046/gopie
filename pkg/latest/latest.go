// Package latest provides latest value.
package latest

// New send latest value from returned chan to `receiver` chan.
func New(receiver chan<- interface{}) chan<- interface{} {
	sender := make(chan interface{})

	go func() {
		var (
			latest interface{}
			ok     bool
			temp   chan<- interface{}
		)
		for {
			select {
			case latest, ok = <-sender:
				if !ok {
					return
				}
				if temp == nil {
					temp = receiver
				}
				continue
			case temp <- latest:
				break
			}
		}
	}()

	return sender
}

// NewN send latest `n` values from returned chan to `receiver` chan.
func NewN(receiver chan<- []interface{}, n int) chan<- interface{} {
	if n < 1 {
		panic("n should be positive")
	}
	sender := make(chan interface{})
	store := make([]interface{}, n)

	go func() {
		var (
			latest interface{}
			ok     bool
			temp   chan<- []interface{}
		)
		for {
			select {
			case latest, ok = <-sender:
				if !ok {
					return
				}
				store = append(store[1:], latest)

				if nil == temp {
					temp = receiver
				}
				continue
			case temp <- append(store[:0:0], store...):
				break
			}
		}
	}()

	return sender
}
