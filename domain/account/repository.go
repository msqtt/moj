package account

type AccountRepository interface {
	FindAccountByID(accountID string) (*Account, error)
	Save(*Account) error
}
