package middleware

import (
	"context"
	"go.uber.org/zap"
	"net/http"
	"timekeeper/app/auth"
	"timekeeper/app/database/model"
)

func UseAuth(authenticator auth.Authenticator) func(next http.Handler) http.Handler {
	log := zap.L().Named("middleware").Named("auth")
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			cookie, err := request.Cookie("SESSION")
			if err != nil {
				log.Debug("failed to get session cookie", zap.Error(err))
				next.ServeHTTP(writer, request)
				return
			}

			err = authenticator.AuthenticateToken(cookie.Value)
			if err != nil {
				log.Debug("failed to authenticate token", zap.Error(err))
				next.ServeHTTP(writer, request)
				return
			}

			next.ServeHTTP(writer, authenticateCookie(authenticator, log, request))
		})
	}
}

type authCtxKey string

var roleKey = authCtxKey("role")

func authenticateCookie(authenticator auth.Authenticator, log *zap.Logger, request *http.Request) *http.Request {
	cookie, err := request.Cookie("SESSION")
	if err != nil {
		log.Debug("failed to get session cookie", zap.Error(err))
		return request
	}

	err = authenticator.AuthenticateToken(cookie.Value)
	if err != nil {
		log.Debug("failed to authenticate token", zap.Error(err))
		return request
	}

	ctx := request.Context()
	ctx = context.WithValue(ctx, roleKey, model.RoleOrganizer)
	return request.WithContext(ctx)
}

func IsOrganizer(r *http.Request) bool {
	ctx := r.Context()
	if v := ctx.Value(roleKey); v != nil {
		role, ok := v.(model.Role)
		return ok && role == model.RoleOrganizer
	}
	return false
}
