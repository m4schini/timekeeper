package render

import (
	"fmt"
	"go.uber.org/zap"
	"maragu.dev/gomponents"
	"net/http"
	"time"
	"timekeeper/ports/www/middleware"
)

func Render(w http.ResponseWriter, r *http.Request, node gomponents.Node) error {
	if !middleware.IsOrganizer(r) {
		revalidate := 30 * time.Second
		SetCache(w, 15*time.Second, &revalidate)
	} else {
		w.Header().Set("Cache-Control", "no-cache")
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	return node.Render(w)
}

func SetCache(w http.ResponseWriter, maxAge time.Duration, revalidate *time.Duration) {
	if revalidate == nil {
		revalidate = &maxAge
	}

	w.Header().Set("Cache-Control", fmt.Sprintf("public, max-age=%v, stale-while-revalidate=%v, immutable", int(maxAge.Seconds()), int(revalidate.Seconds())))
	w.Header().Set("Expires", time.Now().Add(maxAge).Format(http.TimeFormat))
}

func RenderError(log *zap.Logger, w http.ResponseWriter, code int, message string, err error) {
	log.Error(message, zap.Error(err), zap.Int("status", code))
	http.Error(w, message, code)
}
