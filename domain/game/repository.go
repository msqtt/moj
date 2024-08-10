package game

import "context"

type GameRepository interface {
	FindGameByID(ctx context.Context, gameID string) (*Game, error)
	Save(ctx context.Context, game *Game) (gameID string, err error)
	InsertSignUpAccount(ctx context.Context, GameID, accountID string, time int64) error
	DeletSignUpAccount(ctx context.Context, GameID, accountID string) error
}
