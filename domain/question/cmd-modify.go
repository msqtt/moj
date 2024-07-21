package question

type ModifyQuestionCmd struct {
	QuestionID       string
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

func NewModifyQuestionCmdHandler(repo QuestionRepository) *ModifyQuestionCmdHandler {
	return &ModifyQuestionCmdHandler{
		repo: repo,
	}
}

func (h *ModifyQuestionCmdHandler) Handle(cmd ModifyQuestionCmd) error {
	ques, err := h.repo.FindQuestionByID(cmd.QuestionID)
	if err != nil {
		return err
	}
	if ques == nil {
		return ErrQuestionNotFound
	}
	ques2, err := NewQuestion(ques.QuestionID, cmd.Title, cmd.Text, cmd.Level,
		cmd.AllowedLanguages, cmd.TimeLimit, cmd.MemoryLimit, cmd.Tags,
		ques.CreateTime, cmd.Time, cmd.Cases)
	if err != nil {
		return err
	}
	return h.repo.Save(ques2)
}
