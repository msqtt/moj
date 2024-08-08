package domain

import (
	"moj/domain/pkg/queue"
)

type SimpleEventQueue struct {
	queue []any
}

func NewSimpleEventQueue() queue.EventQueue {
	return &SimpleEventQueue{queue: make([]any, 0)}
}

// EnQueue implements queue.EventQueue.
func (s *SimpleEventQueue) EnQueue(event any) error {
	s.queue = append(s.queue, event)
	return nil
}

func (s *SimpleEventQueue) Queue() []any {
	return s.queue
}

var _ queue.EventQueue = (*SimpleEventQueue)(nil)
