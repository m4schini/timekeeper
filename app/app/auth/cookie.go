package auth

import (
	"net/http"
	"time"
)

var (
	CookieName = "SESSION"
)

func SetSessionCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     CookieName,
		Value:    token,
		Path:     "/",
		Domain:   "localhost:8080",
		Expires:  time.Now().Add(24 * 6 * time.Hour),
		MaxAge:   int((18 * time.Hour).Seconds()),
		Secure:   false,
		HttpOnly: true,

		SameSite: http.SameSiteStrictMode,
	})
}

func ClearSessionCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     CookieName,
		Value:    "",
		Path:     "/",
		Domain:   "",
		Expires:  time.Now(),
		MaxAge:   0,
		Secure:   false,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})
}
