package auth

import (
	"context"
	"raumzeitalpaka/app/database/model"
)

type identityCtxKey string

var identityKey = identityCtxKey("identity")

type Identity struct {
	User int
	Name string
	Role model.Role
}

func WithIdentity(ctx context.Context, userId int, userName string, role model.Role) context.Context {
	return context.WithValue(ctx, identityKey, Identity{
		User: userId,
		Name: userName,
		Role: role,
	})
}

func IdentityFrom(ctx context.Context) (identity Identity, isAuthenticated bool) {
	fallbackIdentity := Identity{
		User: -1,
		Name: "unknown",
		Role: model.RoleParticipant,
	}

	// get user
	if v := ctx.Value(identityKey); v != nil {
		var ok bool
		identity, ok = v.(Identity)
		if !ok {
			return fallbackIdentity, false
		}
	} else {
		return fallbackIdentity, false
	}

	return identity, true
}
