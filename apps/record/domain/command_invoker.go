package domain

import (
	"context"

	"moj/apps/record/db"
	"moj/domain/pkg/queue"
)

type CommandInvoker interface {
	Invoke(run func(queue.EventQueue) (any, error)) error
}

type TransactionCommandInvoker struct {
	transactionManager db.TransactionManager
	eventDispatcher    EventDispatcher
}

func NewTransactionCommandInvoker(tm db.TransactionManager,
	ed EventDispatcher) CommandInvoker {
	return &TransactionCommandInvoker{
		transactionManager: tm,
		eventDispatcher:    ed,
	}
}

func (t *TransactionCommandInvoker) Invoke(run func(queue.EventQueue) (any, error)) error {
	return t.transactionManager.Do(context.Background(), func(ctx context.Context) (any, error) {
		queue := NewSimpleEventQueue()
		res, err1 := run(queue)
		go t.eventDispatcher.Dispatch(queue)
		return res, err1
	})
}
