package game

type GameRepository interface {
	FindGameByID(gameID int) (*Game, error)
	Save(game *Game) error
	InsertSignUpAccount(GameID, accountID int, time int64) error
	DeletSignUpAccount(GameID, accountID int) error
}
