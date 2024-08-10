package listener

import (
	"context"
	"moj/apps/game/db"
	"moj/domain/game"
	"time"
)

type SignUpScoreLisener struct {
	signUpScoreDao db.SignUpScoreDao
}

// OnEvent implements Listener.
func (s *SignUpScoreLisener) OnEvent(event any) (err error) {
	ctx := context.Background()
	switch evt := event.(type) {
	case game.CalculateScoreEvent:
		score, err1 := s.signUpScoreDao.FindByID(ctx, evt.GameID, evt.AccountID)
		if err1 != nil {
			return err1
		}
		score.Score = evt.Score
		err = s.signUpScoreDao.Save(ctx, score)
	case game.SignUpGameEvent:
		model := &db.SignUpScoreViewModel{
			GameID:     evt.GameID,
			AccountID:  evt.AccountID,
			Score:      0,
			SignUpTime: time.Unix(evt.SignUpTime, 0),
		}
		err = s.signUpScoreDao.Save(ctx, model)
	case game.CancelSignUpGameEvent:
		err = s.signUpScoreDao.Delete(ctx, evt.GameID, evt.AccountID)
	default:
	}
	return
}

func NewSignUpScoreLisener(
	signUpScoreDao db.SignUpScoreDao,
) Listener {
	return &SignUpScoreLisener{
		signUpScoreDao: signUpScoreDao,
	}
}
