package handler

import (
	"context"
	"raumzeitalpaka/app/auth/user"
)

type UnauthenticatedErr struct{}

func (u UnauthenticatedErr) Error() string {
	return "unauthenticated"
}

func requireAuthentication(ctx context.Context) error {
	_, isAuthenticated := user.IdentityFrom(ctx)
	if !isAuthenticated {
		return UnauthenticatedErr{}
	}
	return nil
}
