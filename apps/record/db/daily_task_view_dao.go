package db

import (
	"context"
	"errors"
	"log/slog"
	"moj/apps/record/pkg"
	"moj/apps/record/pkg/app_err"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type DailyTaskViewDao interface {
	FindByDate(ctx context.Context, dateTime time.Time) (*DailyTaskViewModel, error)
	Save(ctx context.Context, model *DailyTaskViewModel) (id string, err error)
	SumOneFinishByDate(ctx context.Context, dateTime time.Time) error
	SumOneSubmitByDate(ctx context.Context, dateTime time.Time) error
}

type MongoDBDayTaskViewDao struct {
	mongodb    *MongoDB
	collection *mongo.Collection
}

// SumOneFinishByDate implements DayTaskViewDao.
func (m *MongoDBDayTaskViewDao) SumOneFinishByDate(ctx context.Context, dateTime time.Time) error {
	slog.Debug("count finish record day task view model by date", "dateTime", dateTime)
	t := pkg.GetDateTime(dateTime)
	result := m.collection.FindOneAndUpdate(ctx, bson.M{"time": t}, bson.M{"": bson.M{"submit_number": 1}})

	err := result.Err()
	if err != nil {
		err = errors.Join(app_err.ErrServerInternal,
			errors.New("failed to count finish record day task view by date"), err)
	}
	slog.Debug("count finish record day task view model result", "result", result)
	return err
}

// SumOneSubmitByDate implements DayTaskViewDao.
func (m *MongoDBDayTaskViewDao) SumOneSubmitByDate(ctx context.Context, dateTime time.Time) error {
	slog.Debug("count submit record day task view model by date", "dateTime", dateTime)
	t := pkg.GetDateTime(dateTime)
	result := m.collection.FindOneAndUpdate(ctx, bson.M{"time": t}, bson.M{"": bson.M{"submit_number": 1}})

	err := result.Err()
	if err != nil {
		err = errors.Join(app_err.ErrServerInternal,
			errors.New("failed to count submit record day task view by date"), err)
	}
	slog.Debug("count submit record day task view model result", "result", result)
	return err
}

// FindByDate implements DayTaskViewDao.
func (m *MongoDBDayTaskViewDao) FindByDate(ctx context.Context, dateTime time.Time) (*DailyTaskViewModel, error) {
	slog.Debug("find record day task view model by date", "dateTime", dateTime)

	t := pkg.GetDateTime(dateTime)

	var model DailyTaskViewModel
	err := m.collection.FindOne(ctx, bson.M{"time": t}).Decode(&model)
	if err != nil {
		err = errors.Join(app_err.ErrServerInternal,
			errors.New("failed to find record day task view by date"), err)
	}
	return &model, err
}

// Save implements DayTaskViewDao.
func (m *MongoDBDayTaskViewDao) Save(ctx context.Context, model *DailyTaskViewModel) (id string, err error) {
	slog.Debug("insert record day task view model", "model", model)

	result, err := m.collection.InsertOne(ctx, model)
	if err != nil {
		err = errors.Join(app_err.ErrServerInternal,
			errors.New("insert record day task view model failed"), err)
	}
	slog.Debug("insert record day task view model result", "result", result)
	id = result.InsertedID.(primitive.ObjectID).Hex()
	return
}

func NewMongoDBDayTaskViewDao(mongodb *MongoDB) DailyTaskViewDao {
	collection := mongodb.Database().Collection("record_day_task")
	return &MongoDBDayTaskViewDao{
		mongodb:    mongodb,
		collection: collection,
	}
}
