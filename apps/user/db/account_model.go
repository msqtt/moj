package db

import (
	"moj/domain/account"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AccountModel struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	Email      string             `bson:"email"`
	Password   string             `bson:"password"`
	AvatarLink string             `bson:"avatar_link"`
	NickName   string             `bson:"nick_name"`
	Enabled    bool               `bson:"enabled"`
	IsAdmin    bool               `bson:"is_admin"`
}

func NewAccountModelFromAggregate(agg *account.Account) *AccountModel {
	id, _ := primitive.ObjectIDFromHex(agg.AccountID)
	return &AccountModel{
		ID:         id,
		Email:      agg.Email,
		Password:   agg.Password,
		AvatarLink: agg.AvatarLink,
		NickName:   agg.NickName,
		Enabled:    agg.Enabled,
		IsAdmin:    agg.IsAdmin,
	}
}

func (am *AccountModel) ToAggregate() *account.Account {
	return &account.Account{
		AccountID:  am.ID.Hex(),
		Email:      am.Email,
		Password:   am.Password,
		AvatarLink: am.AvatarLink,
		NickName:   am.NickName,
		Enabled:    am.Enabled,
		IsAdmin:    am.IsAdmin,
	}
}
