package pages

import (
	"fmt"
	"net/http"
	"raumzeitalpaka/app/database/model"
	"raumzeitalpaka/ports/www/components"
	"raumzeitalpaka/ports/www/render"

	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func PixelHackPage() Node {
	g := Group{}
	for _, name := range components.AvailablePixelHackIcons {
		g = append(g, Li(Style("display: flex; align-items: center; gap: 1rem; margin: 0"),
			components.PixelHackIcon(name, 24),
			components.CopyTextBox("", "", fmt.Sprintf(":%v:", name)),
		))
	}
	return components.Shell("",
		Main(
			components.PageHeader(model.EventModel{}),
			Ul(g),
		),
	)
}

type PixelHackPageRoute struct {
}

func (l *PixelHackPageRoute) Method() string {
	return http.MethodGet
}

func (l *PixelHackPageRoute) Pattern() string {
	return "/help/pixelhack"
}

func (l *PixelHackPageRoute) Handler() http.Handler {
	log := components.Logger(l)
	page := PixelHackPage()
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		render.HTML(log, writer, request, page)
	})
}
