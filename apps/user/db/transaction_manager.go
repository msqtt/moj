package db

import (
	"context"
	"log/slog"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

type TransactionManager interface {
	Do(ctx context.Context, callback func(ctx context.Context) (any, error)) error
}

type MongoDBTransactionManager struct {
	mongoDB *MongoDB
}

func NewMongoDBTransactionManager(mongoDB *MongoDB) TransactionManager {
	return &MongoDBTransactionManager{mongoDB: mongoDB}
}

// Do implements TransactionManager.
func (m *MongoDBTransactionManager) Do(ctx context.Context,
	callback func(ctx context.Context) (any, error)) error {

	wc := writeconcern.Majority()
	rc := readconcern.Snapshot()
	txnOptions := options.Transaction().SetWriteConcern(wc).SetReadConcern(rc)

	session, err := m.mongoDB.Client().StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)

	slog.Info("start transaction", "session_id", session.ID())

	res, err := session.WithTransaction(ctx, func(ctx mongo.SessionContext) (interface{}, error) {
		return callback(ctx)
	}, txnOptions)

	if err != nil {
		slog.Error("transaction failed, rollback", "session_id", session.ID(), "error", err)
	} else {
		slog.Info("end transaction", "session_id", session.ID(), "result", res)
	}

	return err
}

var _ TransactionManager = (*MongoDBTransactionManager)(nil)
