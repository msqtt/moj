package db

import (
	"time"
)

type AccountViewModel struct {
	AccountID            string    `bson:"account_id"`
	Email                string    `bson:"email"`
	AvatarLink           string    `bson:"avatar_link"`
	NickName             string    `bson:"nick_name"`
	Enabled              bool      `bson:"enabled"`
	IsAdmin              bool      `bson:"is_admin"`
	LastLoginTime        time.Time `bson:"last_login_time"`
	LastLoginIPAddr      string    `bson:"last_login_ip_addr"`
	LastLoginDevice      string    `bson:"last_login_device"`
	LastPasswdChangeTime time.Time `bson:"last_passwd_change_time"`
	RegisterTime         time.Time `bson:"register_time"`
	DeleteTime           time.Time `bson:"delete_time"`
}
