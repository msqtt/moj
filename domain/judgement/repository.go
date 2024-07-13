package judgement

type JudgementRepository interface {
	findJudgementByID(id int) (*Judgement, error)
	findJudgementByHash(questionID int, hash string, questionTime int64) (*Judgement, error)
	save(judgement *Judgement) error
}
