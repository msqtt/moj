package pkg

import "time"

func GetDateTime(t1 time.Time) (dateTime time.Time) {
	year, month, day := t1.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, time.Local)
}
