package sync

import "errors"

const (
	Success = true
	Error   = false
)

type Waiter struct {
	Erch chan<- error
	prev *chan bool
	next *chan bool
}

func NewClosedWaiter(erch chan<- error) *Waiter {
	w := NewWaiter(erch)
	w.Close()
	w.Destroy()
	return w
}

func NewWaiter(erch chan<- error) *Waiter {
	next := make(chan bool, 1)
	return &Waiter{
		Erch: erch,
		next: &next,
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
		Erch: w.Erch,
		prev: w.next,
		next: &next,
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
	w.Erch <- err
	if w.next != nil {
		*w.next <- Error
	}
}
