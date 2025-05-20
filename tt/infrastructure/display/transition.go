package display

import "fmt"

type pair[T ~int] struct {
	from T
	to   T
}

type Handler[T ~int] func(from T, to T) error

type StatusHandler[T ~int] struct {
	items map[pair[T]]Handler[T]
	from  T
}

func NewStatusHandler[T ~int]() *StatusHandler[T] {
	return &StatusHandler[T]{
		items: make(map[pair[T]]Handler[T]),
	}
}

func (t *StatusHandler[T]) Register(from T, to T, handler Handler[T]) {
	t.items[pair[T]{from: from, to: to}] = handler
}

func (t *StatusHandler[T]) UpdateStatus(to T) error {
	if handler, ok := t.items[pair[T]{from: t.from, to: to}]; ok {
		err := handler(t.from, to)
		if err != nil {
			return err
		}

		t.from = to
	}

	return fmt.Errorf("invalid transition from %d to %d", t.from, to)
}
