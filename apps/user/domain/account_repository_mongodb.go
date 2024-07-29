package domain

import (
	"context"
	"errors"
	"log/slog"

	"moj/apps/user/db"
	"moj/apps/user/etc"
	inter_error "moj/apps/user/pkg/app_err"
	"moj/domain/account"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDBAccountRepository struct {
	conf              *etc.Config
	mongodb           *db.MongoDB
	accountCollection *mongo.Collection
}

func NewMongoDBAccountRepository(conf *etc.Config, mongodb *db.MongoDB) account.AccountRepository {
	accountCollection := mongodb.Client().Database(conf.DatabaseName).Collection("account")
	return &MongoDBAccountRepository{
		conf:              conf,
		mongodb:           mongodb,
		accountCollection: accountCollection,
	}
}

// FindAccountByID implements account.AccountRepository.
func (m *MongoDBAccountRepository) FindAccountByID(accountID string) (*account.Account, error) {
	var ret db.AccountModel
	id, _ := primitive.ObjectIDFromHex(accountID)
	err := m.accountCollection.FindOne(context.TODO(),
		bson.M{"_id": id}).Decode(&ret)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			err = errors.Join(
				inter_error.ErrModelNotFound,
				account.ErrAccountNotFound, err)
			return nil, err
		}
		err = errors.Join(inter_error.ErrServerInternal,
			errors.New("failed to find account"),
			err)
		return nil, err
	}
	slog.Debug("find account result", "result", ret)
	return ret.ToAggregate(), err
}

// Save implements account.AccountRepository.
func (m *MongoDBAccountRepository) Save(acc *account.Account) (err error) {
	model := db.NewAccountModelFromAggregate(acc)
	slog.Debug("save account aggreation", "acc", acc)

	var result any
	if model.ID.IsZero() {
		// insert
		result, err = m.accountCollection.InsertOne(context.TODO(), model)
		if err != nil {
			err = errors.Join(errors.New("failed to insert account"), err)
		}
	} else {
		// update
		result, err = m.accountCollection.UpdateOne(context.TODO(),
			bson.M{"_id": model.ID}, bson.M{"$set": model})
		if err != nil {
			err = errors.Join(errors.New("failed to update account"), err)
		}
	}
	slog.Debug("save account result", "result", result)
	if err != nil {
		err = errors.Join(inter_error.ErrServerInternal, err)
	}
	return
}

// FindAccountByEmail implements account.AccountRepository.
func (m *MongoDBAccountRepository) FindAccountByEmail(email string) (*account.Account, error) {
	var ret db.AccountModel

	filter := bson.D{
		{Key: "email", Value: email},
		{Key: "enabled", Value: true},
	}
	opts := &options.FindOneOptions{Sort: bson.M{"create_time": -1}}

	slog.Debug("find latest account by email", "email", email, "filter", filter, "options", opts)

	err := m.accountCollection.FindOne(context.TODO(), filter, opts).Decode(&ret)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			err = errors.Join(
				inter_error.ErrModelNotFound, account.ErrAccountNotFound, err)
			return nil, err
		}
		err = errors.Join(inter_error.ErrServerInternal,
			errors.New("failed to find account"), err)
	}
	slog.Debug("find latest account view by email", "result", ret)
	return ret.ToAggregate(), err
}

var _ account.AccountRepository = (*MongoDBAccountRepository)(nil)
