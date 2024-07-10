package account

type LoginAccountEvent struct {
	AccountID   int
	LoginIPAddr string
	LoginDevice string
	LoginTime   int64
}
