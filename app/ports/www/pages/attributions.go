package pages

import (
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
	"net/http"
	"raumzeitalpaka/app/database/model"
	"raumzeitalpaka/ports/www/components"
	"raumzeitalpaka/ports/www/render"
)

func AttributionsPage() Node {

	return Shell("",
		Main(
			components.PageHeader(model.EventModel{}),
			H2(Text("Attributions")),
			PixelhackSection(),
			FontSection(),
		),
	)
}

func PixelhackSection() Node {
	return Div(
		H3(Text("Pixelart: Pixelhack")),
		P(Text(components.PixelHackAttribution())),
		A(Text("License"), Href("http://creativecommons.org/licenses/by-sa/4.0/"), Target("_blank")),
	)
}

func FontSection() Node {
	return Div(
		H3(Text("Font: Atkinson Hyperlegible Next")),
		P(Style("white-space: pre-line"), Text(`Copyright 2020, Braille Institute of America, Inc. (https://www.brailleinstitute.org/), with Reserved Font Names: “ATKINSON” and “HYPERLEGIBLE”.
This Font Software is licensed under the SIL Open Font License, Version 1.1. This license is copied below, and is also available with a FAQ at: https://openfontlicense.org`)),
		A(Text("License"), Href("/static/font/Atkinson-Hyperlegible-SIL-OPEN-FONT-LICENSE-Version 1.1-v2 ACC.pdf"), Target("_blank")),
	)
}

type AttributionsPageRoute struct {
}

func (l *AttributionsPageRoute) Method() string {
	return http.MethodGet
}

func (l *AttributionsPageRoute) Pattern() string {
	return "/help/attributions"
}

func (l *AttributionsPageRoute) Handler() http.Handler {
	log := components.Logger(l)
	page := AttributionsPage
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		render.HTML(log, writer, request, page())
	})
}
