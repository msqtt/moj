package game

type GameRepository interface {
	FindGameByID(gameID string) (*Game, error)
	Save(game *Game) error
	InsertSignUpAccount(GameID, accountID string, time int64) error
	DeletSignUpAccount(GameID, accountID string) error
}
