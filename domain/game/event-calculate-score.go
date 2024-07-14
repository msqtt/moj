package game

type CalculateScoreEvent struct {
	GameID     int
	AccountID  int
	QuestionID int
	Score      int
}
