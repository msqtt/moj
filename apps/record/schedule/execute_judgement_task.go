package schedule

import (
	"context"
	"log/slog"
	"moj/apps/record/db"
	mq "moj/apps/record/mq/producer"
	"moj/domain/record"
)

type ExecuteJudgementTask struct {
	dao      db.RecordViewDao
	producer mq.Producer
}

// Close implements Tasker.
func (e *ExecuteJudgementTask) Close() {
	e.producer.Close()
}

// Launch implements Worker.
func (e *ExecuteJudgementTask) Launch() {
	records, err := e.dao.FindAllUnFinished(context.Background())
	if err != nil {
		slog.Error("failed to do execute judgement work", "error", err)
		return
	}

	for _, red := range records {
		msg := record.SubmitRecordEvent{
			RecordID:   red.ID.Hex(),
			AccountID:  red.AccountID,
			QuestionID: red.QuestionID,
			GameID:     red.GameID,
			Language:   red.Language,
			Code:       red.Code,
			CodeHash:   red.CodeHash,
			CreateTime: red.CreateTime.Unix(),
		}
		if err := e.producer.Send(msg); err != nil {
			slog.Error("failed to send message to execute judge queue", "error", err)
		}
	}
}

func NewExecuteJudgementTask(
	dao db.RecordViewDao,
	producer mq.Producer,
) Tasker {
	return &ExecuteJudgementTask{
		dao:      dao,
		producer: producer,
	}
}
