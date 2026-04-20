package routes

import "net/http"

type MultiMethodRoute struct {
	// Route pattern
	Route string
	Get   http.Handler
	Post  http.Handler
}

func (m *MultiMethodRoute) Method() string {
	return ""
}

func (m *MultiMethodRoute) Pattern() string {
	return m.Route
}

func (m *MultiMethodRoute) Handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case m.Get != nil && r.Method == http.MethodGet:
			m.Get.ServeHTTP(w, r)
			break
		case m.Post != nil && r.Method == http.MethodPost:
			m.Post.ServeHTTP(w, r)
			break
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
	})
}
