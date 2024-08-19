package listener

import (
	"context"
	"errors"
	"log/slog"
	"moj/game/pkg/app_err"
	"moj/record/db"
	"moj/record/pkg"
	"moj/domain/record"
	"time"
)

type DailyTaskViewListener struct {
	dao db.DailyTaskViewDao
}

// OnEvent implements Listener.
func (d *DailyTaskViewListener) OnEvent(event any) error {
	ctx := context.Background()
	switch event.(type) {
	case record.SubmitRecordEvent:
		dateTime := pkg.GetDateTime(time.Now())
		_, err := d.dao.FindByDate(ctx, dateTime)
		if errors.Is(err, app_err.ErrModelNotFound) {
			model := db.DailyTaskViewModel{
				SubmitNumber: 1,
				FinishNumber: 0,
				Time:         dateTime,
			}
			_, err := d.dao.Save(ctx, &model)
			if err != nil {
				slog.Error("failed to save daily task view", "err", err)
				return err
			}
		} else {
			err := d.dao.SumOneSubmitByDate(ctx, dateTime)
			if err != nil {
				slog.Error("failed to sum one for submit daily task view", "err", err)
				return err
			}
		}
	case record.ModifyRecordEvent:
		dateTime := pkg.GetDateTime(time.Now())
		_, err := d.dao.FindByDate(ctx, dateTime)
		if err != nil {
			slog.Error("failed to find daily task view", "err", err)
			return err
		}
		err = d.dao.SumOneFinishByDate(ctx, dateTime)
		if err != nil {
			slog.Error("failed to sum one for modify daily task view", "err", err)
			return err
		}

	default:
	}
	return nil
}

func NewDailyTaskViewListener(dao db.DailyTaskViewDao) Listener {
	return &DailyTaskViewListener{
		dao: dao,
	}
}
