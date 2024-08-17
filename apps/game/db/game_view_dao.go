package db

import (
	"context"
	"errors"
	"log/slog"
	"moj/apps/game/pkg/app_err"
	"moj/domain/game"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type GameViewDao interface {
	FindGamePage(ctx context.Context, cursor string, pageSize int, f map[string]any) (games []*GameModel, err error)
	DeleteGame(ctx context.Context, gid string) error
}

type MongoDBGameViewDao struct {
	collection *mongo.Collection
	mongodb    *MongoDB
}

// DeleteGame implements GameViewDao.
func (m *MongoDBGameViewDao) DeleteGame(ctx context.Context, gid string) error {
	id, _ := primitive.ObjectIDFromHex(gid)
	slog.Debug("delete game", "id", id)
	result, err := m.collection.DeleteOne(ctx, bson.M{"game_id": id})
	if err != nil {
		err = errors.Join(app_err.ErrServerInternal,
			errors.New("cannot delete game"), err)
		return err
	}
	slog.Debug("delete game", "id", id, "result", result)
	return err
}

func NewMongoDBGameViewDao(m *MongoDB) GameViewDao {
	return &MongoDBGameViewDao{
		collection: m.Database().Collection("game"),
		mongodb:    m,
	}
}

// FindGamePage implements GameDao.
func (m *MongoDBGameViewDao) FindGamePage(ctx context.Context, cursor string, pageSize int, f map[string]any) (games []*GameModel, err error) {
	filter := bson.D{}

	id, _ := primitive.ObjectIDFromHex(cursor)

	if cursor != "" {
		filter = append(filter, bson.E{Key: "_id", Value: bson.M{"$gt": id}})
	}

	for k, v := range f {
		if k == "word" {
			reg := bson.M{"$regex": f["word"], "$options": "i"}
			filter = append(filter, bson.E{Key: "title", Value: reg})
			continue
		}
		if k == "time" {
			filter = append(filter, bson.E{Key: "start_time", Value: bson.M{"$gte": v}},
				bson.E{Key: "end_time", Value: bson.M{"$lt": v}},
			)
			continue
		}
		filter = append(filter, bson.E{Key: k, Value: v})
	}

	opts := options.Find().SetLimit(int64(pageSize)).SetSort(bson.M{"game_id": 1})

	slog.Debug("find game view by page", "filter", filter, "opts", opts)

	games = []*GameModel{}
	cur, err := m.collection.Find(ctx, filter, opts)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			err = errors.Join(
				app_err.ErrModelNotFound, game.ErrGameNotFound, err)
			return nil, err
		}
		err = errors.Join(app_err.ErrServerInternal,
			errors.New("failed to find question view"), err)
		return nil, err
	}
	defer cur.Close(ctx)

	err = cur.All(ctx, &games)
	if err != nil {
		err = errors.Join(app_err.ErrServerInternal,
			errors.New("failed to get all game view"), err)
		return nil, err
	}
	return
}

func NewmongoDBGameDao() GameViewDao {
	return &MongoDBGameViewDao{}
}
