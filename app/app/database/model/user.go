package model

import "time"

type UserModel struct {
	ID           int
	LoginName    string
	DisplayName  string
	PasswordHash string
	LastLogin    time.Time
}

type UserOrganisationMembership struct {
	OrganisationID int
	Slug           string
	Name           string
	Role           string
}

type CreateUserModel struct {
	LoginName    string
	PasswordHash string
}
