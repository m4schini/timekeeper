package middleware

import (
	"net/http"
	"raumzeitalpaka/app/auth"
)

func IsOrganizer(r *http.Request) bool {
	return auth.IsOrganizer(r)
}
