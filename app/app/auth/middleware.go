package auth

import (
	"net/http"
)

func IsOrganizer(r *http.Request) bool {
	_, isAuthenticated := UserFrom(r)
	return isAuthenticated //TODO && role == model.RoleOrganizer
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

	userId, err := AuthenticateJWT(cookie.Value)
	if err != nil {
		return request
	}

	ctx = WithIdentity(ctx, Identity{
		User: userId,
	})
	return request.WithContext(ctx)
}

func UserFrom(r *http.Request) (userId int, isAuthenticated bool) {
	ctx := r.Context()
	id, isAuthenticated := IdentityFrom(ctx)
	return id.User, isAuthenticated
}
