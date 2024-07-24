package account

type AccountRepository interface {
	FindAccountByID(accountID string) (*Account, error)
	FindAccountByEmail(email string) (*Account, error)
	Save(*Account) error
}
