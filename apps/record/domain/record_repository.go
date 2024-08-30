package domain

import (
	"context"
	"errors"
	"log/slog"
	"moj/domain/record"
	"moj/record/db"
	"moj/record/pkg/app_err"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDBRecordRepository struct {
	mongodb    *db.MongoDB
	collection *mongo.Collection
}

// FindBestRecord implements record.RecordRepository.
func (m *MongoDBRecordRepository) FindBestRecord(ctx context.Context, uid string, qid string, gid string) (*record.Record, error) {
	filter := bson.M{"account_id": uid, "question_id": qid}
	opts := options.FindOne().SetSort(bson.M{"number_finish_at": -1})
	if gid != "" {
		filter["game_id"] = gid
	}
	slog.Debug("find best record", "filter", filter, "opts", opts)
	var ret record.Record
	err := m.collection.FindOne(ctx, filter, opts).Decode(&ret)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return &ret, nil
		}
		err = errors.Join(app_err.ErrServerInternal,
			errors.New("failed to find best record by id"), err)
		return nil, err
	}
	return &ret, err
}

// FindRecordByID implements record.RecordRepository.
func (m *MongoDBRecordRepository) FindRecordByID(ctx context.Context, recordID string) (*record.Record, error) {
	var ret record.Record
	slog.Debug("find record by id", "recordID", recordID)
	err := m.collection.FindOne(ctx, bson.M{"_id": recordID}).Decode(&ret)
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
func (m *MongoDBRecordRepository) Save(ctx context.Context, r *record.Record) (id string, err error) {
	model := db.NewRecordFromAggregation(r)
	slog.Debug("start to save record", "model", model)
	var result any
	if model.ID.IsZero() {
		result1, err1 := m.collection.InsertOne(ctx, model)
		if err1 != nil {
			err = errors.Join(app_err.ErrServerInternal, errors.New("failed to insert record"), err1)
		}
		id = result1.InsertedID.(primitive.ObjectID).Hex()
		result = result1
	} else {
		result1, err1 := m.collection.UpdateByID(ctx, model.ID, bson.M{"$set": model})
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
