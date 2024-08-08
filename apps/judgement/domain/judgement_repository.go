package domain

import (
	"context"
	"errors"
	"log/slog"
	"moj/apps/judgement/db"
	"moj/apps/judgement/pkg/app_err"
	"moj/domain/judgement"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoDBJudementRepository struct {
	mongodb    *db.MongoDB
	collection *mongo.Collection
}

// FindJudgementByHash implements judgement.JudgementRepository.
func (m *MongoDBJudementRepository) FindJudgementByHash(questionID string, hash string, questionTime int64) (*judgement.Judgement, error) {
	filter := bson.M{
		"question_id":  questionID,
		"code_hash":    hash,
		"execute_time": bson.M{"$gt": time.Unix(questionTime, 0)},
	}
	slog.Debug("find judgement by hash", "filter", filter)

	var ret judgement.Judgement
	err := m.collection.FindOne(context.TODO(), filter).Decode(&ret)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			err = errors.Join(app_err.ErrModelNotFound, judgement.ErrJudgementNotFound, err)
			return nil, err
		}
		err = errors.Join(app_err.ErrServerInternal, errors.New("failed to find judgement by hash"), err)
	}
	return &ret, err
}

// FindJudgementByID implements judgement.JudgementRepository.
func (m *MongoDBJudementRepository) FindJudgementByID(jid string) (*judgement.Judgement, error) {
	id, _ := primitive.ObjectIDFromHex(jid)
	var ret judgement.Judgement
	slog.Debug("find judgement by id", "judgmentID", jid)
	err := m.collection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&ret)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			err = errors.Join(app_err.ErrModelNotFound, judgement.ErrJudgementNotFound, err)
			return nil, err
		}
		err = errors.Join(app_err.ErrServerInternal, errors.New("failed to find judgement by id"), err)
	}
	return &ret, err
}

// Save implements judgement.JudgementRepository.
func (m *MongoDBJudementRepository) Save(judgement *judgement.Judgement) (err error) {
	model := db.NewJudgementModelFromAggreation(judgement)
	var result any
	if model.ID.IsZero() {
		result, err = m.collection.InsertOne(context.TODO(), model)
		if err != nil {
			err = errors.Join(errors.New("failed to insert judgement"), err)
		}
	} else {
		result, err = m.collection.UpdateByID(context.TODO(), model.ID, bson.M{"$set": model})
		if err != nil {
			err = errors.Join(errors.New("failed to update judgement"), err)
		}
	}

	if err != nil {
		err = errors.Join(app_err.ErrServerInternal, err)
	}
	slog.Debug("save judgement", "result", result)
	return
}

func NewMongoDBJudementRepository(mongodb *db.MongoDB) judgement.JudgementRepository {
	collection := mongodb.Database().Collection("judgement")
	return &MongoDBJudementRepository{
		mongodb:    mongodb,
		collection: collection,
	}
}
