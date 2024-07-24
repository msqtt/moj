package db

import (
	"moj/domain/captcha"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CapthcaModel struct {
	ID         primitive.ObjectID  `bson:"_id,omitempty"`
	AccountID  string              `bson:"account_id"`
	Code       string              `bson:"code"`
	Email      string              `bson:"email"`
	Type       captcha.CaptchaType `bson:"type"`
	IpAddr     string              `bson:"ip_addr"`
	Enabled    bool                `bson:"enabled"`
	Duration   int64               `bson:"duration"`
	CreateTime time.Time           `bson:"create_time"`
	ExpireTime time.Time           `bson:"expire_time"`
}

func NewCaptchaFromAggregate(c *captcha.Captcha) *CapthcaModel {
	id, _ := primitive.ObjectIDFromHex(c.CaptchaID)
	return &CapthcaModel{
		ID:         id,
		AccountID:  c.AccountID,
		Code:       c.Code,
		Email:      c.Email,
		Type:       c.Type,
		IpAddr:     c.IpAddr,
		Enabled:    c.Enabled,
		Duration:   c.Duration,
		CreateTime: time.Unix(c.CreateTime, 0),
		ExpireTime: time.Unix(c.ExpireTime, 0),
	}
}

func (c *CapthcaModel) ToAggregate() *captcha.Captcha {
	return &captcha.Captcha{
		CaptchaID:  c.ID.Hex(),
		AccountID:  c.AccountID,
		Code:       c.Code,
		Email:      c.Email,
		Type:       c.Type,
		IpAddr:     c.IpAddr,
		Enabled:    c.Enabled,
		Duration:   c.Duration,
		CreateTime: c.CreateTime.Unix(),
		ExpireTime: c.ExpireTime.Unix(),
	}
}
