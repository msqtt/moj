package question

type ModifyQuestionCmd struct {
	QuestionID       int
	Enabled          bool
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

type ModifyQuestionCmdHandler struct {
	repo QuestionRepository
}

func (h *ModifyQuestionCmdHandler) Handle(cmd ModifyQuestionCmd) error {
	ques, err := h.repo.findQuestionByID(cmd.QuestionID)
	if err != nil {
		return err
	}
	ques2, err := NewQuestion(ques.QuestionID, cmd.Title, cmd.Text, cmd.Level,
		cmd.AllowedLanguages, cmd.TimeLimit, cmd.MemoryLimit, cmd.Tags,
		ques.CreateTime, cmd.Time, cmd.Cases)
	if err != nil {
		return err
	}
	return h.repo.save(ques2)
}
