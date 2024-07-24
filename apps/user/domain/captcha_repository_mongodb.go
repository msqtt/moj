package domain

import (
	"context"
	"errors"
	"log/slog"

	"moj/apps/user/db"
	"moj/apps/user/etc"
	inter_error "moj/apps/user/pkg/app_err"
	"moj/domain/captcha"
	svc_account "moj/domain/service/account"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ErrFaieldToFindCaptcha = errors.New("failed to find captcha")
	ErrFailedToSaveCaptcha = errors.New("failed to save captcha")
)

type MongoDBCaptchaRepository struct {
	conf              *etc.Config
	mongodb           *db.MongoDB
	captchaCollection *mongo.Collection
}

func NewMongoDBCaptchaRepository(
	conf *etc.Config,
	mongodb *db.MongoDB,
) captcha.CaptchaRepository {
	captchaCollection := mongodb.Client().Database(conf.DatabaseName).Collection("captcha")
	return &MongoDBCaptchaRepository{
		conf:              conf,
		mongodb:           mongodb,
		captchaCollection: captchaCollection,
	}
}

// FindLatestCaptcha implements captcha.CaptchaRepository.
func (m *MongoDBCaptchaRepository) FindLatestCaptcha(email string, code string, captchaType captcha.CaptchaType) (*captcha.Captcha, error) {
	var ret db.CapthcaModel

	filter := bson.D{
		{Key: "email", Value: email},
		{Key: "code", Value: code},
		{Key: "type", Value: captchaType},
		{Key: "enabled", Value: true},
	}

	// 找到最新创建的一条
	option := options.FindOneOptions{Sort: bson.D{{Key: "create_time", Value: -1}}}

	slog.Debug("find latest captcha", "filter", filter, "option", option)
	err := m.captchaCollection.FindOne(context.TODO(), filter, &option).Decode(&ret)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			err = errors.Join(
				inter_error.ErrModelNotFound,
				svc_account.ErrCaptchaNotFound, err)
			return nil, err
		}
		err = errors.Join(inter_error.ErrServerInternal,
			errors.New("failed to find captcha"),
			err)
	}
	slog.Debug("find latest captcha result", "result", ret)
	return ret.ToAggregate(), err
}

// Save implements captcha.CaptchaRepository.
func (m *MongoDBCaptchaRepository) Save(captcha *captcha.Captcha) (err error) {
	model := db.NewCaptchaFromAggregate(captcha)
	slog.Debug("save captcha aggreation", "captcha", captcha)
	slog.Debug("save captcha model", "model", model)
	var result any
	if model.ID.IsZero() {
		result, err = m.captchaCollection.InsertOne(context.TODO(), model)
		if err != nil {
			err = errors.Join(errors.New("failed to save captcha"), err)
		}
	} else {
		// update
		result, err = m.captchaCollection.UpdateOne(context.TODO(),
			bson.M{"_id": model.ID}, bson.M{"$set": model})
		if err != nil {
			err = errors.Join(errors.New("failed to update captcha"), err)
		}
		slog.Debug("update captcha result", "result", result)
	}
	slog.Debug("save captcha result", "result", result)

	if err != nil {
		err = errors.Join(inter_error.ErrServerInternal, ErrFailedToSaveCaptcha, err)
	}
	return err
}

var _ captcha.CaptchaRepository = (*MongoDBCaptchaRepository)(nil)
