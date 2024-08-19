package domain

import (
	"context"
	"errors"
	"log/slog"
	"moj/question/db"
	"moj/question/pkg/app_err"
	"moj/domain/question"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoDBQuestionRepository struct {
	questionCollection *mongo.Collection
}

func NewMongoDBQuestionRepository(mongodb *db.MongoDB) question.QuestionRepository {
	coll := mongodb.Database().Collection("question")
	return &MongoDBQuestionRepository{
		questionCollection: coll,
	}
}

// FindQuestionByID implements question.QuestionRepository.
func (m *MongoDBQuestionRepository) FindQuestionByID(ctx context.Context, questionID string) (*question.Question, error) {
	id, _ := primitive.ObjectIDFromHex(questionID)
	questionModel := db.QuestionModel{}
	err := m.questionCollection.
		FindOne(ctx, bson.M{"_id": id}).
		Decode(&questionModel)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			err = errors.Join(
				app_err.ErrModelNotFound,
				question.ErrQuestionNotFound, err)
			return nil, err
		}
		err = errors.Join(app_err.ErrServerInternal,
			errors.New("failed to find question"),
			err)
	}
	slog.Debug("find question result", "result", questionModel)
	return questionModel.ToAggregate(), err
}

// Save implements question.QuestionRepository.
func (m *MongoDBQuestionRepository) Save(ctx context.Context, ques *question.Question) (id string, err error) {
	model := db.NewFromAggreation(ques)
	slog.Debug("save question model", "model", ques)

	var result any
	if model.ID.IsZero() {
		// insert
		result1, err1 := m.questionCollection.InsertOne(ctx, model)
		if err1 != nil {
			err = errors.Join(errors.New("failed to insert question"), err1)
		}
		id = result1.InsertedID.(primitive.ObjectID).Hex()
		result = result1
	} else {
		// update
		result, err = m.questionCollection.
			UpdateByID(ctx, model.ID, bson.M{"$set": model})
		if err != nil {
			err = errors.Join(errors.New("failed to update question"), err)
		}
		id = model.ID.Hex()
	}
	slog.Debug("save question result", "result", result)
	if err != nil {
		err = errors.Join(app_err.ErrServerInternal, err)
	}
	return
}
