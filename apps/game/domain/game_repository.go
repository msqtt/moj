package domain

import (
	"context"
	"errors"
	"log/slog"
	"moj/apps/game/db"
	"moj/apps/game/pkg/app_err"
	"moj/domain/game"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoDBGameRepository struct {
	mongodb    *db.MongoDB
	collection *mongo.Collection
}

// DeletSignUpAccount implements game.GameRepository.
func (m *MongoDBGameRepository) DeletSignUpAccount(ctx context.Context, gameID string, accountID string) error {
	id, _ := primitive.ObjectIDFromHex(gameID)
	result, err := m.collection.UpdateByID(ctx, id, bson.M{"$pull": bson.M{"sign_up_account_list": bson.M{"account_id": accountID}}})
	if err != nil {
		err = errors.Join(app_err.ErrServerInternal, errors.New("failed to delete sign up account"), err)
	}
	slog.Debug("delete sign up account result", "result", result)
	return err
}

// FindGameByID implements game.GameRepository.
func (m *MongoDBGameRepository) FindGameByID(ctx context.Context, gameID string) (*game.Game, error) {
	id, _ := primitive.ObjectIDFromHex(gameID)
	slog.Debug("find game by id", "gameID", gameID)
	var model db.GameModel
	err := m.collection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&model)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.Join(app_err.ErrModelNotFound, game.ErrGameNotFound, err)
		}
		return nil, errors.Join(app_err.ErrServerInternal, errors.New("failed to find game by id"), err)
	}
	return model.ToAggreation(), nil
}

// InsertSignUpAccount implements game.GameRepository.
func (m *MongoDBGameRepository) InsertSignUpAccount(ctx context.Context, gameID string, accountID string, ti int64) error {
	id, _ := primitive.ObjectIDFromHex(gameID)
	sTime := time.Unix(ti, 0)
	result, err := m.collection.UpdateByID(context.TODO(), id,
		bson.M{"$push": bson.M{"sign_up_account_list": bson.M{"account_id": accountID, "sign_up_time": sTime}}})

	if err != nil {
		err = errors.Join(app_err.ErrServerInternal, errors.New("failed to insert sign up account"), err)
	}
	slog.Debug("insert sign up account result", "result", result)
	return err
}

// Save implements game.GameRepository.
func (m *MongoDBGameRepository) Save(ctx context.Context, game *game.Game) (id string, err error) {
	model := db.NewFromAggreation(game)

	slog.Debug("save game model", "model", game)

	var result any
	if model.ID.IsZero() {
		result1, err1 := m.collection.InsertOne(context.TODO(), model)
		if err1 != nil {
			err = errors.Join(errors.New("failed to insert game"), err1)
		}
		id = result1.InsertedID.(primitive.ObjectID).Hex()
		result = result1
	} else {
		result, err = m.collection.UpdateByID(context.TODO(), model.ID, bson.M{"$set": model})
		if err != nil {
			err = errors.Join(errors.New("failed to update game"), err)
		}
		id = model.ID.Hex()
	}
	if err != nil {
		err = errors.Join(app_err.ErrServerInternal, err)
	}
	slog.Debug("save game result", "result", result)
	return
}

func NewMongoDBGameRepository(mongodb *db.MongoDB) game.GameRepository {
	colle := mongodb.Database().Collection("game")
	return &MongoDBGameRepository{
		mongodb:    mongodb,
		collection: colle,
	}
}
