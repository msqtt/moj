package judgement

type JudgementRepository interface {
	FindJudgementByID(id string) (*Judgement, error)
	FindJudgementByHash(questionID string, hash string, questionModifyTime int64) (*Judgement, error)
	Save(judgement *Judgement) error
}
