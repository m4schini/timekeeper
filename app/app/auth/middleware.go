package auth

import (
	"net/http"
	"raumzeitalpaka/app/auth/user"
	"strings"
)

func IsOrganizer(r *http.Request) bool {
	_, isAuthenticated := UserFrom(r)
	return isAuthenticated //TODO && role == model.RoleOrganizer
}

func UseSessionCookie() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			request = parseSessionCookie(request)
			next.ServeHTTP(writer, request)
		})
	}
}

func parseSessionCookie(request *http.Request) *http.Request {
	ctx := request.Context()

	cookie, err := request.Cookie(CookieName)
	if err != nil {
		return request
	}

	userId, err := AuthenticateJWT(cookie.Value)
	if err != nil {
		return request
	}

	ctx = user.WithIdentity(ctx, user.Identity{
		User: userId,
	})
	return request.WithContext(ctx)
}

func UseBearerToken() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			request = parseAuthorizationHeader(request)
			next.ServeHTTP(writer, request)
		})
	}
}

func parseAuthorizationHeader(request *http.Request) *http.Request {
	ctx := request.Context()

	authHeader := request.Header.Get("Authorization")
	token, isBearer := strings.CutPrefix(authHeader, "Bearer ")
	if !isBearer {
		return request
	}

	userId, err := AuthenticateJWT(token)
	if err != nil {
		return request
	}

	ctx = user.WithIdentity(ctx, user.Identity{
		User: userId,
	})
	return request.WithContext(ctx)
}

func UserFrom(r *http.Request) (userId int, isAuthenticated bool) {
	ctx := r.Context()
	id, isAuthenticated := user.IdentityFrom(ctx)
	return id.User, isAuthenticated
}
