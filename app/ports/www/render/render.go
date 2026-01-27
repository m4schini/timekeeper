package render

import (
	"fmt"
	"go.uber.org/zap"
	"maragu.dev/gomponents"
	"net/http"
	"raumzeitalpaka/ports/www/middleware"
	"time"
)

func HTML(log *zap.Logger, w http.ResponseWriter, r *http.Request, node gomponents.Node) {
	if !middleware.IsOrganizer(r) {
		revalidate := 30 * time.Second
		SetCache(w, 15*time.Second, &revalidate)
	} else {
	}
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err := node.Render(w)
	if err != nil {
		log.Error("failed to render node", zap.Error(err), zap.String("method", r.Method), zap.String("route", r.URL.Path))
	}
}

func Error(log *zap.Logger, w http.ResponseWriter, code int, message string, err error) {
	log.Error(message, zap.Error(err), zap.Int("status", code))
	http.Error(w, message, code)
}

func SetCache(w http.ResponseWriter, maxAge time.Duration, revalidate *time.Duration) {
	if revalidate == nil {
		revalidate = &maxAge
	}

	w.Header().Set("Cache-Control", fmt.Sprintf("public, max-age=%v, stale-while-revalidate=%v, immutable", int(maxAge.Seconds()), int(revalidate.Seconds())))
	w.Header().Set("Expires", time.Now().Add(maxAge).Format(http.TimeFormat))
}
