package account

type DeleteAccountEvent struct {
	AccountID  int
	DeleteTime int64
	Enabled    bool
}
