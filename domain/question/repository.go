package question

type QuestionRepository interface {
	FindQuestionByID(questionID string) (*Question, error)
	Save(*Question) error
}
