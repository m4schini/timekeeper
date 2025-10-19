package pages

import (
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
	"net/http"
	"timekeeper/app/database/model"
	"timekeeper/ports/www/components"
	"timekeeper/ports/www/render"
)

func PixelHackPage() Node {
	g := Group{}
	for _, name := range components.AvailablePixelHackIcons {
		g = append(g, Li(Style("display: flex; align-items: center; gap: 1rem; margin: 0"),
			components.PixelHackIcon(name, 24),
			Textf(":%v:", name),
		))
	}
	return Shell("",
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
