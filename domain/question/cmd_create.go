package question

import "context"

type CreateQuestionCmd struct {
	AccountID        string
	Title            string
	Content          string
	Level            QuestionLevel
	AllowedLanguages []QuestionLanguage
	TimeLimit        int
	MemoryLimit      int
	Tags             []string
	Time             int64
	Cases            []Case
}

type CreateQuestionCmdHandler struct {
	repo QuestionRepository
}

func NewCreateQuestionCmdHandler(repo QuestionRepository) *CreateQuestionCmdHandler {
	return &CreateQuestionCmdHandler{
		repo: repo,
	}
}

func (h *CreateQuestionCmdHandler) Handle(ctx context.Context, cmd CreateQuestionCmd) (result any, err error) {
	ques, err := NewQuestion("", cmd.AccountID, cmd.Title, cmd.Content, cmd.Level, cmd.AllowedLanguages,
		cmd.TimeLimit, cmd.MemoryLimit, cmd.Tags, cmd.Time, 0, cmd.Cases)
	if err != nil {
		return nil, err
	}
	return h.repo.Save(ctx, ques)
}
