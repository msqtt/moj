package db

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GameView struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	AccountID   string             `bson:"account_id,omitempty"`
	Title       string             `bson:"title"`
	Description string             `bson:"description"`
	StartTime   time.Time          `bson:"start_time"`
	EndTime     time.Time          `bson:"end_time"`
	CreateTime  time.Time          `bson:"create_time"`
}
