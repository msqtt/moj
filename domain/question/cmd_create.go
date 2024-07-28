package question

type CreateQuestionCmd struct {
	AccountID        string
	Title            string
	Text             string
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

func (h *CreateQuestionCmdHandler) Handle(cmd CreateQuestionCmd) error {
	ques, err := NewQuestion("", cmd.AccountID, cmd.Title, cmd.Text, cmd.Level, cmd.AllowedLanguages,
		cmd.TimeLimit, cmd.MemoryLimit, cmd.Tags, cmd.Time, 0, cmd.Cases)
	if err != nil {
		return err
	}
	return h.repo.Save(ques)
}
