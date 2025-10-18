package www

import (
	"embed"
	"github.com/go-chi/chi/v5"
	"net/http"
	"time"
	"timekeeper/ports/www/render"
)

//go:embed static/*
var embedFs embed.FS

type StaticFileRoute struct {
}

func (s StaticFileRoute) Method() string {
	return http.MethodGet
}

func (s StaticFileRoute) Pattern() string {
	return "/static/{file}"
}

func (s StaticFileRoute) Handler() http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		fileName := chi.URLParam(request, "file")
		var fileContent []byte

		fileContent, err := embedFs.ReadFile("static/" + fileName)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusNotFound)
			return
		}

		render.SetCache(writer, 24*time.Hour, nil)
		writer.Write(fileContent)
	})
}
