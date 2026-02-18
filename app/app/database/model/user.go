package model

import "time"

type UserModel struct {
	ID           int
	LoginName    string
	DisplayName  string
	PasswordHash string
	LastLogin    time.Time
}

type CreateUserModel struct {
	LoginName    string
	PasswordHash string
}
