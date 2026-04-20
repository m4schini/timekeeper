package domain

import "fmt"

type UnauthenticatedErr struct{}

func (u UnauthenticatedErr) Error() string {
	return "unauthenticated"
}

type UnauthorizedErr struct {
	UserId  int
	Action  string
	Subject any
}

func (u UnauthorizedErr) Error() string {
	return fmt.Sprintf("unauthorized: user(%v) attempted %v on %v", u.UserId, u.Action, u.Subject)
}
