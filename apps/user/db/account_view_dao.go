package db

import (
	"context"
	"errors"
	"log/slog"

	"moj/apps/user/etc"
	inter_error "moj/apps/user/pkg/app_err"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ErrAccountViewNotFound = errors.New("account view not found")

type AccountViewDAO interface {
	FindByAccountID(id string) (*AccountViewModel, error)
	FindLatestByEmail(email string) (*AccountViewModel, error)
	Insert(accountView *AccountViewModel) error
	Update(accountView *AccountViewModel) error
}

type MongoDBAccountViewDAO struct {
	conf                   *etc.Config
	mongodb                *MongoDB
	accountVieweCollection *mongo.Collection
}

func NewMongoDBAccountViewDAO(
	conf *etc.Config,
	mongodb *MongoDB,
) AccountViewDAO {
	accountVieweCollection := mongodb.
		Client().
		Database(conf.DatabaseName).
		Collection("view_account")
	return &MongoDBAccountViewDAO{
		conf:                   conf,
		mongodb:                mongodb,
		accountVieweCollection: accountVieweCollection,
	}
}

// FindByAccountID implements AccountViewDAO.
func (a *MongoDBAccountViewDAO) FindByAccountID(id string) (*AccountViewModel, error) {
	view := &AccountViewModel{}
	slog.Debug("find account view by account id", "id", id)
	err := a.accountVieweCollection.FindOne(context.TODO(), bson.M{"account_id": id}).Decode(view)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			err = errors.Join(
				inter_error.ErrModelNotFound,
				ErrAccountViewNotFound, err)
			return nil, err
		}
		err = errors.Join(inter_error.ErrServerInternal,
			errors.New("failed to find account view"), err)
	}
	slog.Debug("find account view result", "result", view)
	return view, err
}

// Insert implements AccountViewDAO.
func (a *MongoDBAccountViewDAO) Insert(accountView *AccountViewModel) error {
	slog.Debug("insert account view", "accountView", accountView)
	result, err := a.accountVieweCollection.InsertOne(context.TODO(), accountView)
	if err != nil {
		err = errors.Join(inter_error.ErrServerInternal,
			errors.New("failed to insert account view"), err)
	}
	slog.Debug("insert account view", "result", result)
	return err
}

// Update implements AccountViewDAO.
func (a *MongoDBAccountViewDAO) Update(accountView *AccountViewModel) error {
	slog.Debug("update account view", "accountView", accountView)
	result, err := a.accountVieweCollection.UpdateOne(context.TODO(),
		bson.M{"account_id": accountView.AccountID}, bson.M{"$set": accountView})
	if err != nil {
		err = errors.Join(inter_error.ErrServerInternal,
			errors.New("failed to update account view"), err)
	}
	slog.Debug("update account view", "result", result)
	return err
}

// FindLatestByEmail implements AccountViewDAO.
func (a *MongoDBAccountViewDAO) FindLatestByEmail(email string) (*AccountViewModel, error) {
	var ret AccountViewModel

	filter := bson.D{
		{Key: "email", Value: email},
		{Key: "enabled", Value: true},
	}
	opts := &options.FindOneOptions{Sort: bson.M{"create_time": -1}}

	slog.Debug("find latest account view by email", "email", email, "filter", filter, "options", opts)

	err := a.accountVieweCollection.FindOne(context.TODO(), filter, opts).Decode(&ret)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			err = errors.Join(
				inter_error.ErrModelNotFound, ErrAccountViewNotFound, err)
			return nil, err
		}
		err = errors.Join(inter_error.ErrServerInternal,
			errors.New("failed to find account view"), err)
	}
	slog.Debug("find latest account view by email", "result", ret)
	return &ret, err
}

var _ AccountViewDAO = (*MongoDBAccountViewDAO)(nil)
