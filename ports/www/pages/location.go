package pages

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
	"net/http"
	"strconv"
	"timekeeper/adapters"
	"timekeeper/app/database"
	"timekeeper/app/database/model"
	"timekeeper/ports/www/components"
	. "timekeeper/ports/www/render"
)

func OsmContainer(el model.EventLocationModel, osm adapters.LookupResponse) Node {
	return Div(Class("osm-container"),
		Div(ID("map"), Style("width: 600px; height: 400px")),
		Div(
			P(Class("location-label"), Strong(Text(el.RelationshipLabel())), Br(), Text(el.RelationshipNote)),
			P(Class("address"), Textf(`%v
%v %v
%v %v`, osm.Name, osm.Address.Road, osm.Address.HouseNumber, osm.Address.Postcode, osm.Address.City)),
		),
	)
}

func LocationPage(event model.EventModel, location model.EventLocationModel, locationOsmData *adapters.LookupResponse, rooms []model.RoomModel) Node {
	roomItems := Group{}
	for _, room := range rooms {
		roomItems = append(roomItems, Div(Class("room"), ID(fmt.Sprintf("room-%v", room.ID)),
			H4(Text(room.Name)),
			P(Text(room.Description)),
		))
	}

	//osmUrl := fmt.Sprintf("https://www.openstreetmap.org/export/embed.html?bbox=%v%2C52.51540800941198%2C13.480079770088198%2C52.51854508785383&amp;layer=mapnik&amp;marker=52.51697657663071%2C13.477188348770142")
	return ShellWithHead(fmt.Sprintf("%v - %v", location.Name, event.Name),
		[]Node{
			Link(Rel("stylesheet"), Href("https://unpkg.com/leaflet@1.9.4/dist/leaflet.css"), Integrity("sha256-p4NxAoJBhIIN+hmNHrzRCf9tD/miZyoHS5obTRR9BMY="), CrossOrigin("")),
		},
		Group{
			components.PageHeader(event),
			Main(
				H2(Text(location.Name)),
				Iff(locationOsmData != nil, func() Node { return OsmContainer(location, *locationOsmData) }),
				H3(Text("Räume")),
				Div(roomItems),
			),
			Script(Src("https://unpkg.com/leaflet@1.9.4/dist/leaflet.js"), Integrity("sha256-20nQCchB9co0qIjJZRGuk2/Z9VM+kNiyxNV1lvTlZBo="), CrossOrigin("")),
			Iff(locationOsmData != nil, func() Node {
				return Script(Rawf(`
var map = L.map('map').setView([%v, %v], 16);

L.tileLayer('https://tile.openstreetmap.org/{z}/{x}/{y}.png', {
    maxZoom: 19,
    attribution: '&copy; <a href="http://www.openstreetmap.org/copyright">OpenStreetMap</a>'
}).addTo(map);

var marker = L.marker([%v, %v]).addTo(map);
`, locationOsmData.Lat, locationOsmData.Lon, locationOsmData.Lat, locationOsmData.Lon))
			}),
		},
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
	DB        *database.Database
	Nominatim *adapters.NominatimClient
}

func (l *LocationPageRoute) Method() string {
	return http.MethodGet
}

func (l *LocationPageRoute) Pattern() string {
	return "/event/{event}/location/{location}"
}

func (l *LocationPageRoute) Handler() http.Handler {
	log := components.Logger(l)
	queries := l.DB.Queries
	nominatim := l.Nominatim
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		eventId, err := strconv.ParseInt(chi.URLParam(request, "event"), 10, 64)
		if err != nil {
			RenderError(log, writer, http.StatusBadRequest, "invalid eventId", err)
			return
		}
		locationId, err := strconv.ParseInt(chi.URLParam(request, "location"), 10, 64)
		if err != nil {
			RenderError(log, writer, http.StatusBadRequest, "invalid locationId", err)
			return
		}

		event, err := queries.GetEvent(int(eventId))
		if err != nil {
			RenderError(log, writer, http.StatusInternalServerError, "failed to get event", err)
			return
		}

		location, err := queries.GetEventLocation(int(eventId), int(locationId))
		if err != nil {
			RenderError(log, writer, http.StatusInternalServerError, "failed to get location", err)
			return
		}

		var locationOsmData *adapters.LookupResponse
		resp, err := nominatim.Lookup(request.Context(), location.OsmId)
		if err == nil {
			locationOsmData = &resp
		}

		rooms, _, err := queries.GetRoomsOfLocation(int(locationId), 0, 100)
		if err != nil {
			RenderError(log, writer, http.StatusInternalServerError, "failed to get rooms", err)
			return
		}

		Render(log, writer, request, LocationPage(event, location, locationOsmData, rooms))
	})
}
