package db

import (
	"context"
	"errors"
	"log/slog"
	"moj/apps/question/pkg/app_err"
	"moj/domain/question"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type QuestionDao interface {
	FindQuestionPage(cursor string, pageSize int, filter map[string]any) ([]*QuestionViewModel, error)
}

type MongoDBQuestionDAO struct {
	questionCollection *mongo.Collection
}

// FindQuestionPage implements QuestionDao.
func (m *MongoDBQuestionDAO) FindQuestionPage(cursor string, pageSize int,
	f map[string]any) ([]*QuestionViewModel, error) {
	filter := bson.D{}

	if cursor != "" {
		id, _ := primitive.ObjectIDFromHex(cursor)
		filter = append(filter, bson.E{Key: "_id", Value: bson.M{"$gt": id}})
	}

	for k, v := range f {
		if k == "word" {
			reg := bson.M{"$regex": f["word"], "$options": "i"}
			filter = append(filter, bson.E{Key: "$or", Value: bson.A{
				bson.M{"question_id": reg},
				bson.M{"title": reg},
				bson.M{"tags": bson.M{"$elemMatch": reg}},
			}})
			continue
		}
		if k == "language" {
			filter = append(filter, bson.E{Key: "allowed_languages", Value: bson.M{"$in": v}})
			continue
		}
		filter = append(filter, bson.E{Key: k, Value: v})
	}

	pipline := mongo.Pipeline{
		{{Key: "$addFields", Value: bson.M{"question_id": bson.M{"$toString": "$_id"}}}},
		{{Key: "$match", Value: filter}},
		{{Key: "$sort", Value: bson.M{"_id": 1}}},
		{{Key: "$limit", Value: pageSize}},
		{{Key: "$project", Value: bson.M{
			"question_id":       1,
			"account_id":        1,
			"enabled":           1,
			"title":             1,
			"level":             1,
			"tags":              1,
			"total_case_number": bson.M{"$size": "$cases"},
			"create_time":       1,
		}}},
	}

	slog.Debug("find question view by page", "pipline", pipline)

	cur, err := m.questionCollection.Aggregate(context.TODO(), pipline)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			err = errors.Join(
				app_err.ErrModelNotFound, question.ErrQuestionNotFound, err)
			return nil, err
		}
		err = errors.Join(app_err.ErrServerInternal,
			errors.New("failed to find question view"), err)
		return nil, err
	}
	defer cur.Close(context.TODO())

	var ret []*QuestionViewModel
	if err := cur.All(context.TODO(), &ret); err != nil {
		err = errors.Join(app_err.ErrServerInternal,
			errors.New("failed to get all question view"), err)
		return nil, err
	}
	slog.Debug("find question view by page", "result", ret)

	return ret, nil
}

func NewMongoDBQuestionDAO(db *MongoDB) QuestionDao {
	questionCollection := db.Database().Collection("question")
	return &MongoDBQuestionDAO{
		questionCollection: questionCollection,
	}
}
