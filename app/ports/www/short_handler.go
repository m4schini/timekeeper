package www

import (
	"fmt"
	"net/http"
	"raumzeitalpaka/app/database/query"
	"raumzeitalpaka/ports/www/render"
	"time"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type ShortEventHandler struct {
	GetEventBySlug query.GetEventBySlug
}

func (s *ShortEventHandler) Method() string {
	return http.MethodGet
}

func (s *ShortEventHandler) Pattern() string {
	return "/e/{slug}"
}

func (s *ShortEventHandler) Handler() http.Handler {
	log := zap.L().Named("ports").Named("www").Named("short")
	return permanentRedirectHandler(log, s.GetEventBySlug, "/event/%v")
}

type ShortEventScheduleHandler struct {
	GetEventBySlug query.GetEventBySlug
}

func (s *ShortEventScheduleHandler) Method() string {
	return http.MethodGet
}

func (s *ShortEventScheduleHandler) Pattern() string {
	return "/s/{slug}"
}

func (s *ShortEventScheduleHandler) Handler() http.Handler {
	log := zap.L().Named("ports").Named("www").Named("short")
	return permanentRedirectHandler(log, s.GetEventBySlug, "/event/%v/schedule")

}

type ShortEventScheduleMHandler struct {
	GetEventBySlug query.GetEventBySlug
}

func (s *ShortEventScheduleMHandler) Method() string {
	return http.MethodGet
}

func (s *ShortEventScheduleMHandler) Pattern() string {
	return "/s/{slug}/m"
}

func (s *ShortEventScheduleMHandler) Handler() http.Handler {
	log := zap.L().Named("ports").Named("www").Named("short")
	return permanentRedirectHandler(log, s.GetEventBySlug, "/event/%v/schedule?role=Participant,Mentor")

}

func permanentRedirectHandler(log *zap.Logger, getEventBySlug query.GetEventBySlug, urltemplate string) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		slug := chi.URLParam(request, "slug")

		id, err := getEventBySlug.Query(query.GetEventBySlugRequest{Slug: slug})
		if err != nil {
			render.Error(log, writer, http.StatusNotFound, "failed to get event by slug", err)
			return
		}

		render.SetCache(writer, 24*time.Hour, nil)
		http.Redirect(writer, request, fmt.Sprintf(urltemplate, id), http.StatusPermanentRedirect)
	})
}
