package sync

import "errors"

const (
	Success = true
	Error   = false
)

type Waiter struct {
	Msgch chan<- string
	prev  *chan bool
	next  *chan bool
}

func NewClosedWaiter(msgch chan<- string) *Waiter {
	w := NewWaiter(msgch)
	w.Close()
	w.Destroy()
	return w
}

func NewWaiter(msgch chan<- string) *Waiter {
	next := make(chan bool, 1)
	return &Waiter{
		Msgch: msgch,
		next:  &next,
	}
}

func (w *Waiter) Destroy() {
	if w.next == nil {
		return
	}
	close(*w.next)
}

func (w *Waiter) Next() *Waiter {
	next := make(chan bool, 1)

	return &Waiter{
		Msgch: w.Msgch,
		prev:  w.next,
		next:  &next,
	}
}

func (w *Waiter) Wait() error {
	if w.prev == nil {
		return nil
	}

	s := <-*w.prev
	w.prev = nil //一回読み込んだら破棄

	if s {
		return nil
	}

	return errors.New("an error occurred in previous processing")
}

func (w *Waiter) Close() {
	if w.next != nil {
		*w.next <- Success
	}
}

func (w *Waiter) Error(err error) {
	w.Msgch <- err.Error()
	if w.next != nil {
		*w.next <- Error
	}
}

func (w *Waiter) Log(msg string) {
	w.Msgch <- msg
}
