package db

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DailyTaskViewModel struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	SubmitNumber int                `bson:"submit_number"`
	FinishNumber int                `bson:"finish_number"`
	Time         time.Time          `bson:"time"`
}
