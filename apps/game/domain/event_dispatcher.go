package domain

import (
	"log/slog"

	"moj/user/listener"
	"moj/domain/pkg/queue"
)

type EventDispatcher interface {
	Dispatch(queue queue.EventQueue)
}

type SyncEventDispatcher struct {
	listener []listener.Listener
}

func NewSyncEventDispatcher(listener ...listener.Listener) EventDispatcher {
	return &SyncEventDispatcher{
		listener: listener,
	}
}

func (s *SyncEventDispatcher) Dispatch(queue queue.EventQueue) {
	for _, event := range queue.Queue() {
		for _, l := range s.listener {
			err := l.OnEvent(event)
			if err != nil {
				slog.Error("dispatch event error", "err", err)
			}
		}
	}
}
