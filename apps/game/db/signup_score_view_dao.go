package db

import (
	"context"
	"errors"
	"log/slog"
	"moj/domain/game"
	"moj/game/pkg/app_err"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SignUpScoreDao interface {
	FindByID(ctx context.Context, gameID, accountID string) (*SignUpScoreViewModel, error)
	FindPage(ctx context.Context, gameID string, page, pageSize int) ([]*SignUpScoreViewModel, int64, error)
	UpdateScore(ctx context.Context, gameID, accountID string, score int) error
	Save(ctx context.Context, model *SignUpScoreViewModel) error
	Delete(ctx context.Context, gameID, accountID string) error
}

type MongoDBSignUpScoreDao struct {
	mongodb    *MongoDB
	collection *mongo.Collection
}

// UpdateScore implements SignUpScoreDao.
func (m *MongoDBSignUpScoreDao) UpdateScore(ctx context.Context, gameID string, accountID string, score int) error {
	slog.Debug("update signup score", "gameID", gameID, "accountID", accountID)
	filter := bson.M{"game_id": gameID, "account_id": accountID}
	result, err := m.collection.UpdateOne(ctx, filter, bson.M{"": bson.M{"score": score}})
	if err != nil {
		err = errors.Join(app_err.ErrServerInternal,
			errors.New("failed to update signup score"), err)
	}

	slog.Debug("update signup user score", "result", result)
	return err
}

// Delete implements SignUpScoreDao.
func (m *MongoDBSignUpScoreDao) Delete(ctx context.Context, gameID, accountID string) error {
	slog.Debug("delete signup score", "gameID", gameID, "accountID", accountID)
	_, err := m.collection.DeleteOne(ctx, bson.M{"game_id": gameID, "account_id": accountID})
	if err != nil {
		err = errors.Join(app_err.ErrServerInternal, errors.New("failed to delete signup score"), err)
	}
	return err
}

// Save implements SignUpScoreDao.
func (m *MongoDBSignUpScoreDao) Save(ctx context.Context, model *SignUpScoreViewModel) (err error) {
	slog.Debug("save signup score", "model", model)
	if model.ID.IsZero() {
		_, err = m.collection.InsertOne(ctx, model)
		if err != nil {
			err = errors.Join(app_err.ErrServerInternal, errors.New("failed to insert signup score"), err)
		}
	} else {
		_, err = m.collection.UpdateByID(ctx, model.ID, bson.M{"$set": model})
		if err != nil {
			err = errors.Join(app_err.ErrServerInternal, errors.New("failed to update signup score"), err)
		}
	}
	return
}

// FindByGameID implements SignUpScoreDao.
func (m *MongoDBSignUpScoreDao) FindByID(ctx context.Context, gameID, accountID string) (*SignUpScoreViewModel, error) {
	slog.Debug("find signup score", "gameID", gameID, "accountID", accountID)
	ret := SignUpScoreViewModel{}
	err := m.collection.FindOne(ctx, bson.M{"game_id": gameID, "account_id": accountID}).Decode(&ret)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			err = errors.Join(app_err.ErrModelNotFound, errors.New("signup score not found"), err)
			return nil, err
		}
		err = errors.Join(app_err.ErrServerInternal, errors.New("find signup score error"), err)
	}
	return &ret, err
}

// FindPage implements SignUpScoreDao.
func (m *MongoDBSignUpScoreDao) FindPage(ctx context.Context, gameID string, page int, pageSize int) ([]*SignUpScoreViewModel, int64, error) {
	filter := bson.M{"game_id": gameID}
	if page < 1 {
		page = 1
	}

	opts1 := options.Count().
		SetSkip(int64((page - 1) * pageSize)).
		SetLimit(int64(pageSize))

	total, err := m.collection.CountDocuments(ctx, filter, opts1)
	if err != nil {
		err = errors.Join(app_err.ErrServerInternal, errors.New("failed to count signup score"), err)
		return nil, 0, err
	}

	opts2 := options.Find().
		SetSkip(int64((page - 1) * pageSize)).
		SetLimit(int64(pageSize)).
		SetSort(bson.M{"sign_up_time": 1})

	cur, err := m.collection.Find(ctx, filter, opts2)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			err = errors.Join(
				app_err.ErrModelNotFound, game.ErrGameNotFound, err)
			return nil, 0, err
		}
		err = errors.Join(app_err.ErrServerInternal,
			errors.New("failed to find question view"), err)
		return nil, 0, err
	}
	defer cur.Close(ctx)

	var games []*SignUpScoreViewModel
	err = cur.All(ctx, &games)
	if err != nil {
		err = errors.Join(app_err.ErrServerInternal,
			errors.New("failed to get all game view"), err)
		return nil, 0, err
	}
	return games, total, nil
}

func NewMongoDBSignUpScoreDao(
	mongodb *MongoDB,
) SignUpScoreDao {
	collection := mongodb.Database().Collection("game_score")
	return &MongoDBSignUpScoreDao{
		mongodb:    mongodb,
		collection: collection,
	}
}
