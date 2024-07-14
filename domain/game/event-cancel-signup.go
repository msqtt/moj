package game

type CancelSignUpGameEvent struct {
	GameID     int
	AccountID  int
	CancelTime int64
}
