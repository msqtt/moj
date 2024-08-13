package db

import (
	"context"
	"errors"
	"log/slog"

	inter_error "moj/apps/user/pkg/app_err"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ErrAccountViewNotFound = errors.New("account view not found")

type AccountViewDAO interface {
	FindByAccountID(ctx context.Context, id string) (*AccountViewModel, error)
	FindLatestByEmail(ctx context.Context, email string) (*AccountViewModel, error)
	FindByPage(ctx context.Context, pageSize int, cursor string, filter map[string]any) ([]*AccountViewModel, error)
	Insert(ctx context.Context, accountView *AccountViewModel) error
	Update(ctx context.Context, accountView *AccountViewModel) error
}

type MongoDBAccountViewDAO struct {
	mongodb               *MongoDB
	accountViewCollection *mongo.Collection
}

func NewMongoDBAccountViewDAO(
	mongodb *MongoDB,
) AccountViewDAO {
	accountVieweCollection := mongodb.
		Database().
		Collection("account_view")
	return &MongoDBAccountViewDAO{
		mongodb:               mongodb,
		accountViewCollection: accountVieweCollection,
	}
}

// FindByAccountID implements AccountViewDAO.
func (a *MongoDBAccountViewDAO) FindByAccountID(ctx context.Context, id string) (*AccountViewModel, error) {
	view := &AccountViewModel{}
	slog.Debug("find account view by account id", "id", id)
	err := a.accountViewCollection.FindOne(ctx, bson.M{"account_id": id}).Decode(view)
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
func (a *MongoDBAccountViewDAO) Insert(ctx context.Context, accountView *AccountViewModel) error {
	slog.Debug("insert account view", "accountView", accountView)
	result, err := a.accountViewCollection.InsertOne(ctx, accountView)
	if err != nil {
		err = errors.Join(inter_error.ErrServerInternal,
			errors.New("failed to insert account view"), err)
	}
	slog.Debug("insert account view", "result", result)
	return err
}

// Update implements AccountViewDAO.
func (a *MongoDBAccountViewDAO) Update(ctx context.Context, accountView *AccountViewModel) error {
	slog.Debug("update account view", "accountView", accountView)
	result, err := a.accountViewCollection.UpdateOne(ctx,
		bson.M{"account_id": accountView.AccountID}, bson.M{"$set": accountView})
	if err != nil {
		err = errors.Join(inter_error.ErrServerInternal,
			errors.New("failed to update account view"), err)
	}
	slog.Debug("update account view", "result", result)
	return err
}

// FindLatestByEmail implements AccountViewDAO.
func (a *MongoDBAccountViewDAO) FindLatestByEmail(ctx context.Context, email string) (*AccountViewModel, error) {
	var ret AccountViewModel

	filter := bson.D{
		{Key: "email", Value: email},
		{Key: "enabled", Value: true},
	}
	opts := &options.FindOneOptions{Sort: bson.M{"create_time": -1}}

	slog.Debug("find latest account view by email", "email", email, "filter", filter, "options", opts)

	err := a.accountViewCollection.FindOne(ctx, filter, opts).Decode(&ret)
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

// FindByPage implements AccountViewDAO.
func (a *MongoDBAccountViewDAO) FindByPage(ctx context.Context, pageSize int,
	cursor string, f map[string]any) ([]*AccountViewModel, error) {
	filter := bson.D{}

	if cursor != "" {
		filter = append(filter, bson.E{Key: "account_id", Value: bson.M{"$gt": cursor}})
	}

	for k, v := range f {
		if k == "word" {
			reg := bson.M{"$regex": f["word"], "$options": "i"}
			filter = append(filter, bson.E{Key: "$or", Value: bson.A{
				bson.M{"account_id": reg},
				bson.M{"nick_name": reg},
				bson.M{"email": reg},
			}})
			continue
		}
		filter = append(filter, bson.E{Key: k, Value: v})
	}

	opts := options.Find().SetSort(bson.M{"account_id": 1}).SetLimit(int64(pageSize))
	slog.Debug("find account view by page", "filter", filter, "options", opts)

	cur, err := a.accountViewCollection.Find(ctx, filter, opts)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			err = errors.Join(
				inter_error.ErrModelNotFound, ErrAccountViewNotFound, err)
			return nil, err
		}
		err = errors.Join(inter_error.ErrServerInternal,
			errors.New("failed to find account view"), err)
		return nil, err
	}
	defer cur.Close(ctx)

	var ret []*AccountViewModel
	if err = cur.All(ctx, &ret); err != nil {
		err = errors.Join(inter_error.ErrServerInternal,
			errors.New("failed to get all account view"), err)
		return nil, err
	}
	slog.Debug("find account view by page", "result", ret)

	return ret, nil
}

var _ AccountViewDAO = (*MongoDBAccountViewDAO)(nil)
