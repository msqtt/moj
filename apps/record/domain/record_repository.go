package domain

import (
	"context"
	"errors"
	"log/slog"
	"moj/apps/record/db"
	"moj/apps/record/pkg/app_err"
	"moj/domain/record"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoDBRecordRepository struct {
	mongodb    *db.MongoDB
	collection *mongo.Collection
}

// FindRecordByID implements record.RecordRepository.
func (m *MongoDBRecordRepository) FindRecordByID(recordID string) (*record.Record, error) {
	var ret record.Record
	slog.Debug("find record by id", "recordID", recordID)
	err := m.collection.FindOne(context.TODO(), bson.M{"_id": recordID}).Decode(&ret)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			err = errors.Join(app_err.ErrModelNotFound, record.ErrRecordNotFound, err)
			return nil, err
		}
		err = errors.Join(app_err.ErrServerInternal, errors.New("failed to find record by id"), err)
	}
	return &ret, err
}

// Save implements record.RecordRepository.
func (m *MongoDBRecordRepository) Save(r *record.Record) (id string, err error) {
	model := db.NewRecordFromAggregation(r)
	slog.Debug("start to save record", "model", model)
	var result any
	if model.ID.IsZero() {
		result1, err1 := m.collection.InsertOne(context.TODO(), model)
		if err1 != nil {
			err = errors.Join(app_err.ErrServerInternal, errors.New("failed to insert record"), err1)
		}
		id = result1.InsertedID.(primitive.ObjectID).Hex()
		result = result1
	} else {
		result1, err1 := m.collection.UpdateByID(context.TODO(), model.ID, bson.M{"$set": model})
		if err1 != nil {
			err = errors.Join(app_err.ErrServerInternal, errors.New("failed to update record"), err1)
		}
		id = model.ID.Hex()
		result = result1
	}
	slog.Debug("save record", "result", result)
	return
}

func NewMongoDBRecordRepository(
	mongoDB *db.MongoDB,
) record.RecordRepository {
	collection := mongoDB.Database().Collection("record")
	return &MongoDBRecordRepository{
		mongodb:    mongoDB,
		collection: collection,
	}
}