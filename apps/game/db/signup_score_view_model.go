package db

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SignUpScoreViewModel struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	GameID     string             `bson:"game_id"`
	AccountID  string             `bson:"account_id"`
	Score      int
	SignUpTime time.Time `bson:"sign_up_time"`
}
