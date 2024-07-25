package domain

import (
	"context"

	"moj/apps/user/db"
	"moj/domain/pkg/queue"
)

type CommandInvoker interface {
	Invoke(run func(queue.EventQueue) error) error
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

func (t *TransactionCommandInvoker) Invoke(run func(queue.EventQueue) error) error {
	return t.transactionManager.Do(context.Background(), func(ctx context.Context) error {
		queue := NewSimpleEventQueue()
		err1 := run(queue)
		t.eventDispatcher.Dispatch(queue)
		return err1
	})
}
