package listener

import (
	"context"
	"errors"
	"log/slog"
	"moj/apps/game/pkg/app_err"
	"moj/apps/record/db"
	"moj/domain/judgement"
	"moj/domain/question"
	"moj/domain/record"
	"time"
)

type PassedQuestionViewListener struct {
	dao          db.PassQuestionViewDao
	questionRepo question.QuestionRepository
}

// OnEvent implements Listener.
func (p *PassedQuestionViewListener) OnEvent(event any) error {
	ctx := context.Background()
	switch evt := event.(type) {
	case record.ModifyRecordEvent:
		_, err := p.dao.FindByAccountIDAndQuestionID(ctx, evt.AccountID, evt.QuestionID)
		if errors.Is(err, app_err.ErrModelNotFound) {
			quest, err := p.questionRepo.FindQuestionByID(ctx, evt.QuestionID)
			if err != nil {
				slog.Error("failed to get question when listen to pass record", "error", err)
				return err
			}

			status := db.PassStatusWorking
			if evt.JudgeStatus == string(judgement.JudgeStatusAC) {
				status = db.PassStatusPass
			}
			model := db.PassedQuestionViewModel{
				AccountID:  evt.AccountID,
				QuestionID: evt.QuestionID,
				Status:     status,
				Level:      quest.Level.String(),
				RecordID:   evt.RecordID,
				GameID:     evt.GameID,
				FinishTime: time.Time{},
			}
			_, err = p.dao.Save(ctx, &model)
			if err != nil {
				slog.Error("failed to save passed question view when listen to pass record", "error", err)
				return err
			}
		}
	default:
	}
	return nil
}

func NewPassedQuestionViewListener(
	dao db.PassQuestionViewDao,
	questionRepo question.QuestionRepository,
) *PassedQuestionViewListener {
	return &PassedQuestionViewListener{
		dao:          dao,
		questionRepo: nil,
	}
}
