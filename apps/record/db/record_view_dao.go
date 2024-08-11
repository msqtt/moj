package db

import (
	"context"
	"errors"
	"log/slog"
	"moj/apps/record/pkg/app_err"
	"moj/domain/record"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type RecordViewDao interface {
	FindPage(ctx context.Context, questionID, accountID string, page, pageSize int, filter map[string]any) (
		[]*RecordModel, int64, error)
	FindAllUnFinished(context.Context) ([]*RecordModel, error)
	CountAllByID(ctx context.Context, questionID, optianalGameID string) (int64, error)
}

type MongoDBRecordViewDao struct {
	mongodb    *MongoDB
	collection *mongo.Collection
}

// CountAllByID implements RecordViewDao.
func (m *MongoDBRecordViewDao) CountAllByID(ctx context.Context, questionID, optianalGameID string) (int64, error) {
	filter := bson.M{"question_id": questionID}
	if optianalGameID != "" {
		filter["game_id"] = optianalGameID

	}
	number, err := m.collection.CountDocuments(ctx, filter)
	if err != nil {
		err = errors.Join(app_err.ErrServerInternal,
			errors.New("failed to count question by id"), err)
	}
	return number, err
}

// FindAllUnFinished implements RecordViewDao.
func (m *MongoDBRecordViewDao) FindAllUnFinished(ctx context.Context) ([]*RecordModel, error) {
	filter := bson.M{"$or": bson.A{
		bson.M{"judge_status": bson.M{"$exists": false}},
		bson.M{"judge_status": bson.M{"$eq": ""}},
		bson.M{"judge_status": nil},
	}}

	cur, err := m.collection.Find(ctx, filter)
	if err != nil {
		err = errors.Join(app_err.ErrServerInternal,
			errors.New("failed to find unfinished record view"), err)
		return nil, err
	}
	var ret []*RecordModel
	err = cur.All(ctx, &ret)
	if err != nil {
		err = errors.Join(app_err.ErrServerInternal,
			errors.New("failed to get unfinished record view"), err)
		return nil, err
	}
	return ret, nil
}

// FindPage implements RecordViewDao.
func (m *MongoDBRecordViewDao) FindPage(ctx context.Context, questionID string,
	accountID string, page int, pageSize int, f map[string]any) ([]*RecordModel, int64, error) {
	filter := bson.D{
		bson.E{Key: "game_id", Value: questionID},
		bson.E{Key: "account_id", Value: accountID},
	}

	for k, v := range f {
		filter = append(filter, bson.E{Key: k, Value: v})
	}

	if page < 1 {
		page = 1
	}

	total, err := m.collection.CountDocuments(ctx, filter)
	if err != nil {
		err = errors.Join(app_err.ErrServerInternal, errors.New("failed to count record view"), err)
		return nil, total, err
	}

	opts := options.Find().SetSkip(int64(page-1) * int64(pageSize)).
		SetLimit(int64(pageSize)).
		SetSort(bson.D{{Key: "create_time", Value: -1}})

	slog.Debug("find record view page", "filter", filter, "opts", opts)

	var ret []*RecordModel
	cur, err := m.collection.Find(ctx, filter, opts)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			err = errors.Join(
				app_err.ErrModelNotFound, record.ErrRecordNotFound, err)
			return nil, 0, err
		}
		err = errors.Join(app_err.ErrServerInternal,
			errors.New("failed to find question view"), err)
		return nil, 0, err
	}
	defer cur.Close(ctx)

	err = cur.All(ctx, &ret)
	if err != nil {
		err = errors.Join(app_err.ErrServerInternal,
			errors.New("failed to get all game view"), err)
		return nil, 0, err
	}
	return ret, total, err
}

func NewMongoDBRecordViewDao(mongodb *MongoDB) RecordViewDao {
	collection := mongodb.Database().Collection("record")
	return &MongoDBRecordViewDao{
		mongodb:    mongodb,
		collection: collection,
	}
}
