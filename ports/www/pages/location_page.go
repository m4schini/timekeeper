package pages

import (
	"fmt"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
	"net/http"
	"timekeeper/app/database/model"
	"timekeeper/ports/www/components"
	. "timekeeper/ports/www/render"
)

func LocationPage() Node {
	boxes := make([]Box, 0)

	return Shell(
		Main(
			components.PageHeader(model.EventModel{}),
			ImageWithBoxes("/static/betahaus.png", 3424, 2080, boxes),
		),
	)
}

type Box struct {
	Title string
	X     float64 // left position (%, px, etc.)
	Y     float64 // top position
	W     float64 // width
	H     float64 // height
}

// ImageWithBoxes renders an image with overlayed boxes
func ImageWithBoxes(imgSrc string, imgWidth, imgHeight float64, boxes []Box) Node {

	return Div(
		Style("position: relative; height: 100%; max-height: 95vh; margin: auto;"),
		// Overlayed boxes
		Map(boxes, func(b Box) Node {
			if b.Title == "betahaus" {
				return Group{}
			}

			leftPct := b.X / imgWidth * 100
			topPct := b.Y / imgHeight * 100
			widthPct := b.W / imgWidth * 100
			heightPct := b.H / imgHeight * 100

			return Div(
				Style(fmt.Sprintf(
					"position: absolute; left:%.2f%%; top:%.2f%%; "+
						"width:%.2f%%; height:%.2f%%; "+
						"display: flex; align-items: center; justify-content: center; "+
						"border: 2px solid red; background: rgba(0,0,0,0.4); "+
						"color: white; font-weight: bold;",
					leftPct, topPct, widthPct, heightPct,
				)),
				Text(b.Title),
			)

		}),
		// Image
		Img(
			Src(imgSrc),
			Style("width: 100%; height: auto; display: block;"),
		),
	)

}

type LocationPageRoute struct {
}

func (l *LocationPageRoute) Method() string {
	return http.MethodGet
}

func (l *LocationPageRoute) Pattern() string {
	return "/location"
}

func (l *LocationPageRoute) Handler() http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		Render(writer, request, LocationPage())
	})
}
