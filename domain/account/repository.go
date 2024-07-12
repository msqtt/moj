package account

type AccountRepository interface {
	findAccountByID(accountID int) (*Account, error)
	save(*Account) error
}
