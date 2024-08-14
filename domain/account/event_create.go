package account

type CreateAccountEvent struct {
	AccountID    string
	AvatarLink   string
	Email        string
	NickName     string
	RegisterTime int64
	Enabled      bool
}
