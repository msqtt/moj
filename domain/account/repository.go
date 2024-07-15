package account

type AccountRepository interface {
	FindAccountByID(accountID int) (*Account, error)
	Save(*Account) error
}
