package signup

import "context"

type UserRepo interface {
	CreateUser(ctx context.Context, username string) (int64, error)
}

type CredentialsRepo interface {
	StoreCredentials(ctx context.Context, userId int64, login string, password string) error
}
