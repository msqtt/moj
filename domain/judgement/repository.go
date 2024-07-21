package judgement

type JudgementRepository interface {
	FindJudgementByID(id string) (*Judgement, error)
	FindJudgementByHash(questionID string, hash string, questionTime int64) (*Judgement, error)
	Save(judgement *Judgement) error
}
