package www

import (
	"embed"
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
	"raumzeitalpaka/ports/www/render"
	"strings"
	"time"
)

//go:embed static/*
var embedFs embed.FS

//go:embed static/font/*
var fontFs embed.FS

//go:embed static/pixelhack/*
var pixelhackFs embed.FS

func PixelHackItems() map[string]string {
	entries, err := pixelhackFs.ReadDir("static/pixelhack")
	if err != nil {
		panic(err)
	}

	names := make(map[string]string)
	for _, entry := range entries {
		name := strings.TrimSuffix(entry.Name(), ".svg")
		names[name] = fmt.Sprintf("/static/pixelhack/%v", entry.Name())
	}

	return names
}

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

type FontFileRoute struct {
}

func (s FontFileRoute) Method() string {
	return http.MethodGet
}

func (s FontFileRoute) Pattern() string {
	return "/static/font/{file}"
}

func (s FontFileRoute) Handler() http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		fileName := chi.URLParam(request, "file")
		var fileContent []byte

		fileContent, err := fontFs.ReadFile("static/font/" + fileName)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusNotFound)
			return
		}

		revalidate := 365 * 24 * time.Hour
		render.SetCache(writer, 30*24*time.Hour, &revalidate)
		writer.Write(fileContent)
	})
}

type PixelhackFileRoute struct {
}

func (s PixelhackFileRoute) Method() string {
	return http.MethodGet
}

func (s PixelhackFileRoute) Pattern() string {
	return "/static/pixelhack/{file}"
}

func (s PixelhackFileRoute) Handler() http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		fileName := chi.URLParam(request, "file")
		var fileContent []byte

		fileContent, err := pixelhackFs.ReadFile("static/pixelhack/" + fileName)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusNotFound)
			return
		}

		writer.Header().Set("Content-Type", "image/svg+xml")
		revalidate := 365 * 24 * time.Hour
		render.SetCache(writer, 24*time.Hour, &revalidate)
		writer.Write(fileContent)
	})
}
