package game

type GameRepository interface {
	findGameByID(gameID int) (*Game, error)
	save(game *Game) error
	insertSignUpAccount(GameID, accountID int, time int64) error
	deletSignUpAccount(GameID, accountID int) error
}
