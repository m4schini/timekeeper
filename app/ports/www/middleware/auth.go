package middleware

import (
	"context"
	"go.uber.org/zap"
	"net/http"
	"raumzeitalpaka/app/auth"
	"raumzeitalpaka/app/database/model"
)

func UseAuth(authenticator auth.Authenticator) func(next http.Handler) http.Handler {
	log := zap.L().Named("middleware").Named("auth")
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			next.ServeHTTP(writer, authenticateCookie(authenticator, log, request))
		})
	}
}

type roleCtxKey string

var roleKey = roleCtxKey("role")

type userCtxKey string

var userKey = userCtxKey("user")

func authenticateCookie(authenticator auth.Authenticator, log *zap.Logger, request *http.Request) *http.Request {
	cookie, err := request.Cookie("SESSION")
	if err != nil {
		//log.Debug("failed to get session cookie", zap.Error(err))
		return request
	}

	userId, role, err := authenticator.AuthenticateToken(cookie.Value)
	if err != nil {
		log.Debug("failed to authenticate token", zap.Error(err))
		return request
	}

	ctx := request.Context()
	ctx = context.WithValue(ctx, roleKey, role)
	ctx = context.WithValue(ctx, userKey, userId)
	return request.WithContext(ctx)
}

func LoadUser(r *http.Request) (userId int, role model.Role, isAuthenticated bool) {
	ctx := r.Context()
	var ok bool

	// get user
	if v := ctx.Value(userKey); v != nil {
		userId, ok = v.(int)
		if !ok {
			return -1, model.RoleParticipant, false
		}
	} else {
		return -1, model.RoleParticipant, false
	}

	// get role
	if v := ctx.Value(roleKey); v != nil {
		role, ok = v.(model.Role)
		if !ok {
			return -1, model.RoleParticipant, false
		}
	} else {
		return -1, model.RoleParticipant, false
	}

	return userId, role, true
}

func IsOrganizer(r *http.Request) bool {
	_, role, isAuthenticated := LoadUser(r)
	return isAuthenticated && role == model.RoleOrganizer
}
