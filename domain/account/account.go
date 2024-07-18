package account

import (
	"errors"
	"github.com/msqtt/moj/domain/pkg/common"
	"github.com/msqtt/moj/domain/pkg/crypt"
	"github.com/msqtt/moj/domain/pkg/queue"
	"regexp"
	"strings"
	"unicode"
)

var (
	ErrAccountNotFound   = errors.New("account not found")
	ErrInValidAvatarLink = errors.New("invalid avatar link")
	ErrInValidNickName   = errors.New("invalid nickname")
	ErrInValidEmail      = errors.New("invalid email")
	ErrInValidPasswd     = errors.New("invalid password")
	ErrDuplicateRemoval  = errors.New("duplicate removal")
)

type Account struct {
	AccountID  int
	Email      string
	Password   string
	AvatarLink string
	NickName   string
	Enabled    bool
	IsAdmin    bool
}

func NewAccount(cry crypt.Cryptor, email, passwd, nickName string) (acc *Account, err error) {
	if !common.IsEmail(email) {
		err = errors.Join(err, ErrInValidEmail)
	}

	if !isPasswd(passwd) {
		err = errors.Join(err, ErrInValidPasswd)
	}

	if !isNickName(nickName) {
		err = errors.Join(err, ErrInValidNickName)
	}

	// avoid unnecessary hash operations
	if err != nil {
		return
	}

	acc = &Account{
		Email:      email,
		Password:   cry.Encrypt(passwd),
		AvatarLink: "",
		NickName:   nickName,
		Enabled:    true,
		IsAdmin:    false,
	}
	return
}

func (a *Account) create(queue queue.EventQueue, cmd CreateAccountCmd) error {
	event := CreateAccountEvent{
		AccountID:    a.AccountID,
		Email:        cmd.Email,
		NickName:     cmd.NickName,
		RegisterTime: cmd.Time,
		Enabled:      true,
	}
	return queue.EnQueue(event)
}

func (a *Account) ValidPasswd(cry crypt.Cryptor, passwd string) (bool, error) {
	return cry.Valid(passwd, a.Password)
}

func (a *Account) login(queue queue.EventQueue, cmd LoginAccountCmd) error {
	event := LoginAccountEvent{
		AccountID:   cmd.AccountID,
		LoginIPAddr: cmd.IPAddr,
		LoginDevice: cmd.Device,
		LoginTime:   cmd.Time,
	}
	return queue.EnQueue(event)
}

func (a *Account) modifyInfo(queue queue.EventQueue, cmd ModifyInfoAccountCmd) error {
	if !common.IsURL(cmd.AvatarLink) {
		return ErrInValidAvatarLink
	}
	if !isNickName(cmd.NickName) {
		return ErrInValidNickName
	}

	a.AvatarLink = cmd.AvatarLink
	a.NickName = cmd.NickName

	event := ModifyAccountInfoEvent{
		AccountID:  a.AccountID,
		NickName:   cmd.NickName,
		AvatarLink: cmd.AvatarLink,
	}
	return queue.EnQueue(event)
}

func (a *Account) changePasswd(cry crypt.Cryptor, queue queue.EventQueue,
	cmd ChangePasswdAccountCmd) error {
	if !isPasswd(cmd.Password) {
		return ErrInValidPasswd
	}

	a.Password = cry.Encrypt(cmd.Password)

	event := ChangePasswdAccountEvent{
		AccountID:  a.AccountID,
		ChangeTime: cmd.Time,
	}
	return queue.EnQueue(event)

}

func (a *Account) delete(queue queue.EventQueue, cmd DeleteAccountCmd) error {
	if !a.Enabled {
		return ErrDuplicateRemoval
	}

	a.Enabled = false

	event := DeleteAccountEvent{
		AccountID:  a.AccountID,
		Enabled:    false,
		DeleteTime: cmd.Time,
	}
	return queue.EnQueue(event)
}

func (a *Account) SetAdmin(queue queue.EventQueue, cmd SetAdminAccountCmd) error {
	a.IsAdmin = cmd.IsAdmin
	event := SetAdminAccountEvent{
		AccountID: a.AccountID,
		IsAdmin:   cmd.IsAdmin,
	}
	return queue.EnQueue(event)
}

func (a *Account) SetStatus(queue queue.EventQueue, cmd SetStatusAccountCmd) error {
	a.Enabled = cmd.Enabled

	event := SetStatusAccountEvent{
		AccountID: a.AccountID,
		Enabled:   cmd.Enabled,
	}
	return queue.EnQueue(event)
}

// isNickName checks whether given name is a valid nickname.
func isNickName(name string) bool {
	regex := regexp.MustCompile(`^[a-zA-Z0-9_-]{3,15}$|^[\p{Han}]{1,5}$`)
	return regex.MatchString(name)
}

// func isPasswd(passwd string) bool {
// 	regex := regexp.MustCompile(
// 		`^(?=.*?[A-Z])(?=.*?[a-z])(?=.*?[0-9])(?=.*?[#?!@$ %^&*-]).{8,20}$`)
// 	return regex.MatchString(passwd)
// }

// isPasswd checks whether given passwd is a valid password.
func isPasswd(passwd string) bool {
	// Ensure password length is between 8 and 20 characters
	if len(passwd) < 8 || len(passwd) > 20 {
		return false
	}

	// Check for at least one uppercase letter, one lowercase letter, one digit,
	// and one special character among the allowed special characters.
	hasUpper := false
	hasLower := false
	hasDigit := false
	hasSpecial := false
	specialChars := "#?!@$ %^&*-"

	for _, ch := range passwd {
		switch {
		case unicode.IsUpper(ch):
			hasUpper = true
		case unicode.IsLower(ch):
			hasLower = true
		case unicode.IsDigit(ch):
			hasDigit = true
		case strings.ContainsRune(specialChars, ch):
			hasSpecial = true
		}
	}

	return hasUpper && hasLower && hasDigit && hasSpecial
}
