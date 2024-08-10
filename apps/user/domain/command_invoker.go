package domain

import (
	"context"

	"moj/apps/user/db"
	"moj/domain/pkg/queue"
)

type CommandInvoker interface {
	Invoke(ctx context.Context, run func(ctx context.Context, queue queue.EventQueue) error) error
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

func (t *TransactionCommandInvoker) Invoke(ctx context.Context, run func(context.Context, queue.EventQueue) error) error {
	return t.transactionManager.Do(ctx, func(ctx context.Context) error {
		queue := NewSimpleEventQueue()
		err1 := run(ctx, queue)
		t.eventDispatcher.Dispatch(queue)
		return err1
	})
}
