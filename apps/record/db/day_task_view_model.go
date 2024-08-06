package db

import "time"

type DayTaskViewModel struct {
	SubmitNumber int
	FinishNumber int
	Time         time.Time
}
