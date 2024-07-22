package game

type CancelSignUpGameEvent struct {
	GameID     string
	AccountID  string
	CancelTime int64
}
