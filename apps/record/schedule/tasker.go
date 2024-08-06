package schedule

import (
	"log/slog"
	"time"
)

type Tasker interface {
	Launch()
	Close()
}

type TikerTasker struct {
	tiker  *time.Ticker
	tasker Tasker
}

func NewTickerTasker(tiker *time.Ticker, t Tasker) *TikerTasker {
	return &TikerTasker{
		tiker:  tiker,
		tasker: t,
	}
}

func (w *TikerTasker) Launch() (cancel func()) {
	cel := make(chan struct{})

	go func() {
		for {
			select {
			case <-w.tiker.C:
				slog.Info("ticker worker doing...")
				w.tasker.Launch()
				slog.Info("ticker worker finish doing")
			case <-cel:
				slog.Info("ticker worker cancel")
				return
			}
		}
	}()

	return func() {
		w.tasker.Close()
		cel <- struct{}{}
	}
}
