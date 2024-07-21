package account

type LoginAccountEvent struct {
	AccountID   string
	LoginIPAddr string
	LoginDevice string
	LoginTime   int64
}
