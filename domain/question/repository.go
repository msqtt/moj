package question

import "context"

type QuestionRepository interface {
	FindQuestionByID(ctx context.Context, questionID string) (*Question, error)
	Save(context.Context, *Question) (questionID string, err error)
}
