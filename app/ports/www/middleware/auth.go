package middleware

import (
	"net/http"
	"raumzeitalpaka/app/auth"
	"raumzeitalpaka/app/auth/authz"
	"raumzeitalpaka/app/database/model"
)

func IsOrganizer(r *http.Request, authorizer authz.Authorizer) bool {
	userId, _ := auth.UserFrom(r)
	return authorizer.HasRole(userId, model.RoleOrganizer)
}

func IsLoggedIn(r *http.Request) bool {
	_, loggedIn := auth.UserFrom(r)
	return loggedIn
}
