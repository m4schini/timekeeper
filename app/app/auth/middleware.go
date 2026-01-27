package auth

import (
	"net/http"
	"raumzeitalpaka/app/database/model"
)

func IsOrganizer(r *http.Request) bool {
	_, role, isAuthenticated := LoadUser(r)
	return isAuthenticated && role == model.RoleOrganizer
}

func UseJWT() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			request = parseJWT(request)
			next.ServeHTTP(writer, request)
		})
	}
}

func parseJWT(request *http.Request) *http.Request {
	ctx := request.Context()

	cookie, err := request.Cookie(CookieName)
	if err != nil {
		return request
	}

	userId, username, role, err := AuthenticateJWT(cookie.Value)
	if err != nil {
		return request
	}

	ctx = WithIdentity(ctx, userId, username, role)
	return request.WithContext(ctx)
}

func LoadUser(r *http.Request) (userId int, role model.Role, isAuthenticated bool) {
	ctx := r.Context()
	id, isAuthenticated := IdentityFrom(ctx)
	return id.User, id.Role, isAuthenticated
}
