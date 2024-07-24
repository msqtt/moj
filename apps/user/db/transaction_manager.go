package db

import (
	"context"
	"log/slog"

	"go.mongodb.org/mongo-driver/mongo"
)

type TransactionManager interface {
	Do(ctx context.Context, callback func(ctx context.Context) error) error
}

type MongoDBTransactionManager struct {
	mongoDB *MongoDB
}

func NewMongoDBTransactionManager(mongoDB *MongoDB) TransactionManager {
	return &MongoDBTransactionManager{mongoDB: mongoDB}
}

// Do implements TransactionManager.
func (m *MongoDBTransactionManager) Do(ctx context.Context,
	callback func(ctx context.Context) error) error {
	session, err := m.mongoDB.Client().StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)

	slog.Info("start transaction", "session_id", session.ID())

	res, err := session.WithTransaction(ctx, func(ctx mongo.SessionContext) (interface{}, error) {
		return nil, callback(ctx)
	})

	if err != nil {
		slog.Error("transaction failed, rollback", "session_id", session.ID(), "error", err)
	} else {
		slog.Info("end transaction", "session_id", session.ID(), "result", res)
	}

	return err
}

var _ TransactionManager = (*MongoDBTransactionManager)(nil)
