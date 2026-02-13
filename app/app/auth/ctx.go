package auth

import (
	"context"
)

type identityCtxKey string

var identityKey = identityCtxKey("identity")

type Identity struct {
	User int
}

func WithIdentity(ctx context.Context, identity Identity) context.Context {
	return context.WithValue(ctx, identityKey, identity)
}

func IdentityFrom(ctx context.Context) (identity Identity, isAuthenticated bool) {
	fallbackIdentity := Identity{
		User: -1,
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
