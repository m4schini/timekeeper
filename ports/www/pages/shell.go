package pages

import (
	_ "embed"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/components"
	. "maragu.dev/gomponents/html"
)

func Shell(title string, children ...Node) Node {
	if title == "" {
		title = "Timekeeper"
	}
	return HTML5(HTML5Props{
		Title:       title,
		Description: "Zeitplan",
		Language:    "de",
		Head: []Node{
			Script(Src("/static/htmx.min.js")),
			Link(Href("/static/style.css"), Rel("stylesheet")),
			Link(Rel("icon"), Type("image/png"), Href("/static/jh_logo_icon.png")),
		},
		Body:      children,
		HTMLAttrs: nil,
	})

}
