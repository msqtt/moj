package policy

import (
	"errors"

	"moj/domain/game"
	"moj/domain/pkg/queue"
	"moj/domain/record"
)

var ErrFailedToCalculateScore error = errors.New("failed to calculate score")

type CalculateScorePolicy struct {
	recordRepository         record.RecordRepository
	calculateScoreCmdHandler game.CalculateScoreCmdHandler
	queue                    queue.EventQueue
}

func NewCalculateScorePolicy(
	recordRepository record.RecordRepository,
	calculateScoreCmdHandler game.CalculateScoreCmdHandler,
	queue queue.EventQueue,
) *CalculateScorePolicy {
	return &CalculateScorePolicy{
		recordRepository:         recordRepository,
		calculateScoreCmdHandler: calculateScoreCmdHandler,
		queue:                    queue,
	}
}

func (p *CalculateScorePolicy) OnEvent(event any) error {
	evt, ok := event.(record.ModifyRecordEvent)
	if !ok {
		return nil
	}

	if evt.GameID == "" {
		return nil
	}

	rec, err := p.recordRepository.FindBestGameRecord(evt.GameID, evt.AccountID)
	if err != nil {
		return err
	}

	if rec.NumberFinishedAt >= evt.NumberFinishedAt {
		return nil
	}

	cmd := game.CalculateScoreCmd{
		GameID:           evt.GameID,
		RecordID:         evt.RecordID,
		AccountID:        evt.AccountID,
		QuestionID:       evt.QuestionID,
		NumberFinishedAt: evt.NumberFinishedAt,
		LastFinishedAt:   rec.NumberFinishedAt,
		TotalQuestion:    evt.TotalQuestion,
	}

	err = p.calculateScoreCmdHandler.Handle(p.queue, cmd)
	if err != nil {
		return errors.Join(ErrFailedToCalculateScore, err)
	}
	return nil
}
