package policy

import (
	"context"
	"errors"
	"log/slog"

	"moj/domain/game"
	domain_err "moj/domain/pkg/error"
	"moj/domain/pkg/queue"
	"moj/domain/record"
)

var ErrFailedToCalculateScore error = errors.New("failed to calculate score")

type CalculateScorePolicy struct {
	calculateScoreCmdHandler *game.CalculateScoreCmdHandler
	queue                    queue.EventQueue
}

func NewCalculateScorePolicy(
	calculateScoreCmdHandler *game.CalculateScoreCmdHandler,
	queue queue.EventQueue,
) *CalculateScorePolicy {
	return &CalculateScorePolicy{
		calculateScoreCmdHandler: calculateScoreCmdHandler,
		queue:                    queue,
	}
}

func (p *CalculateScorePolicy) OnEvent(event any) error {
	ctx := context.Background()
	evt, ok := event.(record.ModifyRecordEvent)
	if !ok {
		return domain_err.ErrEventTypeInvalid
	}

	if evt.GameID == "" {
		return nil
	}

	cmd := game.CalculateScoreCmd{
		GameID:           evt.GameID,
		RecordID:         evt.RecordID,
		AccountID:        evt.AccountID,
		QuestionID:       evt.QuestionID,
		NumberFinishedAt: evt.NumberFinishedAt,
		TotalQuestion:    evt.TotalQuestion,
	}

	slog.Info("start to calculate score", "cmd", cmd)
	err := p.calculateScoreCmdHandler.Handle(ctx, p.queue, cmd)
	if err != nil {
		return errors.Join(ErrFailedToCalculateScore, err)
	}
	return nil
}
