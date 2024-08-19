package domain

import (
	"context"

	"moj/user/db"
	"moj/domain/pkg/queue"
)

type CommandInvoker interface {
	InvokeWithTrans(ctx context.Context, run func(ctx context.Context, queue queue.EventQueue) error) error
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

func (t *TransactionCommandInvoker) InvokeWithTrans(ctx context.Context, run func(context.Context, queue.EventQueue) error) error {
	return t.transactionManager.Do(ctx, func(ctx context.Context) (any, error) {
		queue := NewSimpleEventQueue()
		err1 := run(ctx, queue)
		t.eventDispatcher.Dispatch(queue)
		return nil, err1
	})
}

func (t *TransactionCommandInvoker) Invoke(ctx context.Context, run func(context.Context, queue.EventQueue) error) error {
	queue := NewSimpleEventQueue()
	err := run(ctx, queue)
	t.eventDispatcher.Dispatch(queue)
	return err
}
