package account

type AccountRepo interface {
	findAccountByID(accountID int) (Account, error)
	save(*Account) error
}
