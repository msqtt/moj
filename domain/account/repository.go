package account

import "context"

type AccountRepository interface {
	FindAccountByID(ctx context.Context, accountID string) (*Account, error)
	FindAccountByEmail(ctx context.Context, email string) (*Account, error)
	Save(ctx context.Context, account *Account) error
}
