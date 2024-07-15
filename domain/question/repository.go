package question

type QuestionRepository interface {
	FindQuestionByID(questionID int) (*Question, error)
	Save(*Question) error
}
