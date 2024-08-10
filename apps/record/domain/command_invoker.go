package domain

import (
	"context"

	"moj/apps/record/db"
	"moj/domain/pkg/queue"
)

type CommandInvoker interface {
	Invoke(ctx context.Context, run func(context.Context, queue.EventQueue) (any, error)) error
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

func (t *TransactionCommandInvoker) Invoke(ctx context.Context, run func(context.Context, queue.EventQueue) (any, error)) error {
	return t.transactionManager.Do(ctx, func(ctx context.Context) (any, error) {
		queue := NewSimpleEventQueue()
		res, err1 := run(ctx, queue)
		go t.eventDispatcher.Dispatch(queue)
		return res, err1
	})
}
