package account

type CreateAccountEvent struct {
	AccountID    int
	Email        string
	NickName     string
	RegisterTime int64
	Enabled      bool
}
