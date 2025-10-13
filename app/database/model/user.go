package model

type UserModel struct {
	ID           int
	LoginName    string
	PasswordHash string
}

type CreateUserModel struct {
	LoginName    string
	PasswordHash string
}
