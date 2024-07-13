package question

type QuestionRepository interface {
	findQuestionByID(questionID int) (*Question, error)
	save(*Question) error
}
