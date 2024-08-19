package domain

import (
	"log/slog"

	"moj/judgement/listener"
	mq "moj/judgement/mq/producer"
	"moj/domain/pkg/queue"
)

type EventDispatcher interface {
	Dispatch(localQueue queue.EventQueue)
}

type SyncAndAsyncEventDispatcher struct {
	listener []listener.Listener
	producer []mq.Producer
}

func NewSyncAndAsyncEventDispatcher(listener []listener.Listener, producer []mq.Producer) EventDispatcher {
	return &SyncAndAsyncEventDispatcher{
		listener: listener,
		producer: producer,
	}
}

func (s *SyncAndAsyncEventDispatcher) Dispatch(queue queue.EventQueue) {
	for _, event := range queue.Queue() {
		// deal with remote queue
		for _, p := range s.producer {
			err := p.Send(event)
			if err != nil {
				slog.Error("dispatch remote event error", "err", err)
			}
		}

		// deal with local queue
		for _, l := range s.listener {
			err := l.OnEvent(event)
			if err != nil {
				slog.Error("dispatch local event error", "err", err)
			}
		}
	}
}
