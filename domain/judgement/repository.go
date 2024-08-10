package judgement

import "context"

type JudgementRepository interface {
	FindJudgementByID(ctx context.Context, id string) (*Judgement, error)
	FindJudgementByHash(ctx context.Context, questionID string, hash string, questionModifyTime int64) (*Judgement, error)
	Save(context.Context, *Judgement) error
}
