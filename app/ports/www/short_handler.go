package www

import (
	"fmt"
	"net/http"
	"raumzeitalpaka/app/database"
	"raumzeitalpaka/ports/www/render"
	"time"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type ShortEventHandler struct {
	DB *database.Database
}

func (s *ShortEventHandler) Method() string {
	return http.MethodGet
}

func (s *ShortEventHandler) Pattern() string {
	return "/e/{slug}"
}

func (s *ShortEventHandler) Handler() http.Handler {
	queries := s.DB.Queries
	log := zap.L().Named("ports").Named("www").Named("short")
	return permanentRedirectHandler(log, queries, "/event/%v")
}

type ShortEventScheduleHandler struct {
	DB *database.Database
}

func (s *ShortEventScheduleHandler) Method() string {
	return http.MethodGet
}

func (s *ShortEventScheduleHandler) Pattern() string {
	return "/s/{slug}"
}

func (s *ShortEventScheduleHandler) Handler() http.Handler {
	queries := s.DB.Queries
	log := zap.L().Named("ports").Named("www").Named("short")
	return permanentRedirectHandler(log, queries, "/event/%v/schedule")

}

type ShortEventScheduleMHandler struct {
	DB *database.Database
}

func (s *ShortEventScheduleMHandler) Method() string {
	return http.MethodGet
}

func (s *ShortEventScheduleMHandler) Pattern() string {
	return "/s/{slug}/m"
}

func (s *ShortEventScheduleMHandler) Handler() http.Handler {
	queries := s.DB.Queries
	log := zap.L().Named("ports").Named("www").Named("short")
	return permanentRedirectHandler(log, queries, "/event/%v/schedule?role=Participant,Mentor")

}

func permanentRedirectHandler(log *zap.Logger, queries database.Queries, urltemplate string) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		slug := chi.URLParam(request, "slug")

		id, err := queries.GetEventIdBySlug(slug)
		if err != nil {
			render.Error(log, writer, http.StatusNotFound, "failed to get event by slug", err)
			return
		}

		render.SetCache(writer, 24*time.Hour, nil)
		http.Redirect(writer, request, fmt.Sprintf(urltemplate, id), http.StatusPermanentRedirect)
	})
}
