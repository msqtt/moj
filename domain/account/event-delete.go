package account

type DeleteAccountEvent struct {
	AccountID  string
	DeleteTime int64
	Enabled    bool
}
