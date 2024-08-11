package db

import (
	"context"
	"errors"
	"log/slog"
	"moj/apps/record/pkg/app_err"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type PassedQuestionViewDao interface {
	FindByAccountIDAndQuestionID(ctx context.Context, accountID string, questionID string) (*PassedQuestionViewModel, error)
	CountByQuestionID(ctx context.Context, questionID string) (passNumber int64, err error)
	CountLevelByAccountID(ctx context.Context, accountID string) (eazy, normal, hard int64, err error)
	Save(ctx context.Context, model *PassedQuestionViewModel) (id string, err error)
}

type MongoDBPassedQuestionViewDao struct {
	mondodb    *MongoDB
	collection *mongo.Collection
}

// CountByQuestionID implements PassedQuestionViewDao.
func (m *MongoDBPassedQuestionViewDao) CountByQuestionID(ctx context.Context, questionID string) (passNumber int64, err error) {
	passNumber, err = m.collection.CountDocuments(ctx, bson.M{"question_id": questionID})
	if err != nil {
		err = errors.Join(app_err.ErrServerInternal,
			errors.New("failed to count passed question"), err)
	}
	return
}

// CountLevelByAccountID implements PassedQuestionViewDao.
func (m *MongoDBPassedQuestionViewDao) CountLevelByAccountID(ctx context.Context, accountID string) (eazy int64, normal int64, hard int64, err error) {
	pipeline := mongo.Pipeline{
		{{Key: "$group", Value: bson.M{
			"_id":   "$level",
			"count": bson.M{"$sum": 1},
		}}},
		{{Key: "$match", Value: bson.M{"account_id": accountID,
			"level": bson.M{"$in": bson.A{"eazy", "normal", "hard"}}}}},
	}

	slog.Debug("count passed question", "accountID", accountID, "pipeline", pipeline)

	cur, err := m.collection.Aggregate(ctx, pipeline)
	if err != nil {
		err = errors.Join(app_err.ErrServerInternal,
			errors.New("failed to count passed question"), err)
		return
	}
	ret := []struct {
		ID    string `bson:"_id"`
		Count int64  `bson:"count"`
	}{}
	err = cur.All(ctx, &ret)
	if err != nil {
		err = errors.Join(app_err.ErrServerInternal,
			errors.New("failed to count passed question"), err)
		return
	}

	for _, r := range ret {
		switch r.ID {
		case "eazy":
			eazy = r.Count
		case "normal":
			normal = r.Count
		case "hard":
			hard = r.Count
		}
	}
	return
}

// FindByAccountIDAndQuestionID implements PassedQuestionViewDao.
func (m *MongoDBPassedQuestionViewDao) FindByAccountIDAndQuestionID(ctx context.Context,
	accountID string, questionID string) (*PassedQuestionViewModel, error) {
	var ret PassedQuestionViewModel
	err := m.collection.FindOne(ctx, bson.M{"account_id": accountID, "question_id": questionID}).
		Decode(&ret)
	if err != nil {
		err = errors.Join(app_err.ErrModelNotFound,
			errors.New("failed to find passed question"), err)
	}
	return &ret, err
}

// Save implements PassedQuestionViewDao.
func (m *MongoDBPassedQuestionViewDao) Save(ctx context.Context, model *PassedQuestionViewModel) (id string, err error) {
	slog.Debug("save passed question", "model", model)
	result, err := m.collection.InsertOne(ctx, model)
	if err != nil {
		err = errors.Join(app_err.ErrServerInternal,
			errors.New("failed to save passed question"), err)
	}
	slog.Debug("save passed question result", "result", result)
	id = result.InsertedID.(primitive.ObjectID).Hex()
	return
}

func NewMongoDBPassedQuestionViewDao(mondodb *MongoDB) PassedQuestionViewDao {
	return &MongoDBPassedQuestionViewDao{
		mondodb:    mondodb,
		collection: mondodb.Database().Collection("record_passed_question"),
	}
}
