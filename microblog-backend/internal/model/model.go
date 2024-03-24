package model

import "time"

type User struct {
	Id        int64     `db:"id"`
	Username  string    `db:"username"`
	CreatedAt time.Time `db:"created_at"`
	InviteId  *int64    `db:"invite_id"`
}

type Credentials struct {
	Id           int64      `db:"id"`
	UserId       int64      `db:"user_id"`
	Login        string     `db:"login"`
	PasswordHash string     `db:"password_hash"`
	CreatedAt    time.Time  `db:"created_at"`
	ObsoletedAt  *time.Time `db:"obsoleted_at"` // this field is ignored rn
}

type Invite struct {
	Id        int64      `db:"id"`
	Code      string     `db:"code"`
	UserId    int64      `db:"user_id"`
	CreatedAt time.Time  `db:"created_at"`
	UsedAt    *time.Time `db:"used_at"`
}

type Post struct {
	Id        int64     `db:"id"`
	UserId    int64     `db:"user_id"`
	Body      string    `db:"body"`
	CreatedAt time.Time `db:"created_at"`
	ExpiresAt time.Time `db:"expires_at"`
}

type Session struct {
	UserId    int64
	Key       string
	ExpiresAt time.Time
}
