package account

type CreateAccountEvent struct {
	AccountID    string
	Email        string
	NickName     string
	RegisterTime int64
	Enabled      bool
}
