package middleware

import "net/http"

func AllowAllCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		request.Header.Set("Access-Control-Allow-Origin", "*")
		request.Header.Set("Access-Control-Allow-Credentials", "true")
		request.Header.Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		request.Header.Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if request.Method == "OPTIONS" {
			writer.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(writer, request)
	})
}
