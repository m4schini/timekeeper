package pages

import (
	_ "embed"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/components"
	. "maragu.dev/gomponents/html"
)

func Shell(children ...Node) Node {
	return HTML5(HTML5Props{
		Title:       "Timekeeper",
		Description: "Zeitplan",
		Language:    "de",
		Head: []Node{
			Script(Src("/static/htmx.min.js")),
			Link(Href("/static/style.css"), Rel("stylesheet")),
		},
		Body:      children,
		HTMLAttrs: nil,
	})

}
