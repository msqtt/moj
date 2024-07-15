package judgement

type JudgementRepository interface {
	FindJudgementByID(id int) (*Judgement, error)
	FindJudgementByHash(questionID int, hash string, questionTime int64) (*Judgement, error)
	Save(judgement *Judgement) error
}
