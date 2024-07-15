package question

type CreateQuestionCmd struct {
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

func (h *CreateQuestionCmdHandler) Handle(cmd CreateQuestionCmd) error {
	ques, err := NewQuestion(0, cmd.Title, cmd.Text, cmd.Level, cmd.AllowedLanguages,
		cmd.TimeLimit, cmd.MemoryLimit, cmd.Tags, cmd.Time, 0, cmd.Cases)
	if err != nil {
		return err
	}
	return h.repo.Save(ques)
}
