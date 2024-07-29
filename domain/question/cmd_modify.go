package question

type ModifyQuestionCmd struct {
	QuestionID       string
	Enabled          bool
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

type ModifyQuestionCmdHandler struct {
	repo QuestionRepository
}

func NewModifyQuestionCmdHandler(repo QuestionRepository) *ModifyQuestionCmdHandler {
	return &ModifyQuestionCmdHandler{
		repo: repo,
	}
}

func (h *ModifyQuestionCmdHandler) Handle(cmd ModifyQuestionCmd) (any, error) {
	ques, err := h.repo.FindQuestionByID(cmd.QuestionID)
	if err != nil {
		return nil, err
	}
	ques2, err := NewQuestion(ques.QuestionID, ques.AccountID, cmd.Title, cmd.Content, cmd.Level,
		cmd.AllowedLanguages, cmd.TimeLimit, cmd.MemoryLimit, cmd.Tags,
		ques.CreateTime, cmd.Time, cmd.Cases)
	if err != nil {
		return nil, err
	}
	return h.repo.Save(ques2)
}
