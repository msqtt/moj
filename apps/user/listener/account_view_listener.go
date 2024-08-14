package listener

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"moj/apps/user/db"
	"moj/domain/account"
)

var ErrFailedToUpdateAccountView = errors.New("failed to update account view")

type AccountViewListener struct {
	accountViewDAO db.AccountViewDAO
}

func NewAccountViewListener(dao db.AccountViewDAO) Listener {
	return &AccountViewListener{
		accountViewDAO: dao,
	}
}

// OnEvent implements Listener.
func (a *AccountViewListener) OnEvent(event any) (err error) {
	ctx := context.Background()
	switch evt := event.(type) {
	case account.CreateAccountEvent:
		view := &db.AccountViewModel{
			AccountID:            evt.AccountID,
			Email:                evt.Email,
			AvatarLink:           evt.AvatarLink,
			NickName:             evt.NickName,
			Enabled:              true,
			IsAdmin:              false,
			LastLoginTime:        time.Time{},
			LastLoginIPAddr:      "",
			LastLoginDevice:      "",
			LastPasswdChangeTime: time.Time{},
			RegisterTime:         time.Unix(evt.RegisterTime, 0),
			DeleteTime:           time.Time{},
		}
		err = a.accountViewDAO.Insert(ctx, view)
	case account.DeleteAccountEvent:
		view, err1 := a.accountViewDAO.FindByAccountID(ctx, evt.AccountID)
		if err1 != nil {
			err = errors.Join(err1, err)
			return
		}
		view.Enabled = false
		view.DeleteTime = time.Unix(evt.DeleteTime, 0)

		err2 := a.accountViewDAO.Update(ctx, view)
		if err2 != nil {
			err = errors.Join(err2, err)
		}
	case account.LoginAccountEvent:
		view, err1 := a.accountViewDAO.FindByAccountID(ctx, evt.AccountID)
		if err1 != nil {
			err = errors.Join(err1, err)
			return
		}
		view.LastLoginTime = time.Unix(evt.LoginTime, 0)
		view.LastLoginIPAddr = evt.LoginIPAddr
		view.LastLoginDevice = evt.LoginDevice

		err2 := a.accountViewDAO.Update(ctx, view)
		if err2 != nil {
			err = errors.Join(err2, err)
		}
	case account.ModifyAccountInfoEvent:
		view, err1 := a.accountViewDAO.FindByAccountID(ctx, evt.AccountID)
		if err1 != nil {
			err = errors.Join(err1, err)
			return
		}
		view.NickName = evt.NickName
		view.AvatarLink = evt.AvatarLink

		err2 := a.accountViewDAO.Update(ctx, view)
		if err2 != nil {
			err = errors.Join(err2, err)
		}
	case account.ChangePasswdAccountEvent:
		view, err1 := a.accountViewDAO.FindByAccountID(ctx, evt.AccountID)
		if err1 != nil {
			err = errors.Join(err1, err)
			return
		}
		view.LastPasswdChangeTime = time.Unix(evt.ChangeTime, 0)

		err2 := a.accountViewDAO.Update(ctx, view)
		if err2 != nil {
			err = errors.Join(err2, err)
		}
	case account.SetAdminAccountEvent:
		view, err1 := a.accountViewDAO.FindByAccountID(ctx, evt.AccountID)
		if err1 != nil {
			err = errors.Join(err1, err)
			return
		}
		view.IsAdmin = evt.IsAdmin

		err2 := a.accountViewDAO.Update(ctx, view)
		if err2 != nil {
			err = errors.Join(err2, err)
		}
	case account.SetStatusAccountEvent:
		view, err1 := a.accountViewDAO.FindByAccountID(ctx, evt.AccountID)
		if err1 != nil {
			err = errors.Join(err1, err)
			return
		}
		view.Enabled = evt.Enabled

		err2 := a.accountViewDAO.Update(ctx, view)
		if err2 != nil {
			err = errors.Join(err2, err)
		}
	default:
	}
	if err != nil {
		slog.Error("AccountView failed", "err", err, "event", event)
		err = errors.Join(ErrFailedToUpdateAccountView, err)
	}
	return
}

var _ Listener = (*AccountViewListener)(nil)
