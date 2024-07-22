package game

type CalculateScoreEvent struct {
	GameID     string
	AccountID  string
	QuestionID string
	Score      int
}
