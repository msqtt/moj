package db

import (
	"context"

	"moj/apps/record/etc"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	conf   *etc.Config
	client *mongo.Client
}

func NewMongoDB(conf *etc.Config) *MongoDB {
	clientOpt := options.Client().ApplyURI(conf.MongoHost)
	client, err := mongo.Connect(context.Background(), clientOpt)
	if err != nil {
		panic(err)
	}
	return &MongoDB{client: client, conf: conf}
}

func (m *MongoDB) Close() {
	if err := m.client.Disconnect(context.Background()); err != nil {
		panic(err)
	}
}

func (m *MongoDB) Client() *mongo.Client {
	return m.client
}

func (m *MongoDB) Database() *mongo.Database {
	return m.client.Database(m.conf.MongoDBName)
}
