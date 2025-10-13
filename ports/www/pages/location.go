package pages

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
	"net/http"
	"strconv"
	"timekeeper/app/database"
	"timekeeper/app/database/model"
	"timekeeper/ports/www/components"
	. "timekeeper/ports/www/render"
)

func LocationPage(event model.EventModel, rooms []model.RoomModel) Node {
	boxes := make([]Box, len(rooms))
	for i, room := range rooms {
		boxes[i] = Box{
			Title:  room.Name,
			Anchor: fmt.Sprintf("%v", room.ID),
			X:      float64(room.LocationX),
			Y:      float64(room.LocationY),
			W:      float64(room.LocationW),
			H:      float64(room.LocationH),
		}
	}

	return Shell(
		Main(
			components.PageHeader(event, false),
			//<iframe width="425" height="350" src="https://www.openstreetmap.org/export/embed.html?bbox=13.465826511383058%2C52.51108885548261%2C13.48895788192749%2C52.52363705879041&amp;layer=mapnik" style="border: 1px solid black"></iframe><br/><small><a href="https://www.openstreetmap.org/?#map=16/52.51736/13.47739">View Larger Map</a></small>
			IFrame(Width("425"), Height("350"), Src("https://www.openstreetmap.org/export/embed.html?bbox=13.474296927452087%2C52.51540800941198%2C13.480079770088198%2C52.51854508785383&amp;layer=mapnik&amp;marker=52.51697657663071%2C13.477188348770142")),
			ImageWithBoxes("/static/betahaus2.png", 3424, 2080, boxes),
		),
	)
}

type Box struct {
	Title  string
	Anchor string
	X      float64 // left position (%, px, etc.)
	Y      float64 // top position
	W      float64 // width
	H      float64 // height
}

// ImageWithBoxes renders an image with overlayed boxes
func ImageWithBoxes(imgSrc string, imgWidth, imgHeight float64, boxes []Box) Node {

	return Div(ID("karte"),
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

			return Div(ID(b.Anchor), Class("box"),
				Style(fmt.Sprintf(
					"position: absolute; left:%.2f%%; top:%.2f%%; "+
						"width:%.2f%%; height:%.2f%%; ",
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
	DB *database.Database
}

func (l *LocationPageRoute) Method() string {
	return http.MethodGet
}

func (l *LocationPageRoute) Pattern() string {
	return "/event/{event}/location/{location}"
}

func (l *LocationPageRoute) UseCache() bool {
	return false
}

func (l *LocationPageRoute) Handler() http.Handler {
	log := zap.L().Named(l.Pattern())
	queries := l.DB.Queries
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		eventId, err := strconv.ParseInt(chi.URLParam(request, "event"), 10, 64)
		if err != nil {
			RenderError(log, writer, http.StatusBadRequest, "invalid eventId", err)
			return
		}

		event, err := queries.GetEvent(int(eventId))

		//location, err := strconv.ParseInt(chi.URLParam(request, "location"), 10, 64)
		//if err != nil {
		//	RenderError(log, writer, http.StatusBadRequest, "invalid locationId", err)
		//	return
		//}

		rooms, _, err := queries.GetRooms(0, 100)
		if err != nil {
			RenderError(log, writer, http.StatusInternalServerError, "failed to get rooms", err)
			return
		}

		Render(log, writer, request, LocationPage(event, rooms))
	})
}
