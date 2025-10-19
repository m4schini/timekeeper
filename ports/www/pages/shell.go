package pages

import (
	_ "embed"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/components"
	. "maragu.dev/gomponents/html"
	"timekeeper/ports/www/components"
)

func Shell(title string, children ...Node) Node {
	if title == "" {
		title = "Timekeeper"
	}
	return ShellWithHead(title, []Node{}, children...)
}

func ShellWithHead(title string, head []Node, children ...Node) Node {
	if title == "" {
		title = "Timekeeper"
	}
	return HTML5(HTML5Props{
		Title:       title,
		Description: "Zeitplan",
		Language:    "de",
		Head: append(head,
			Script(Src("/static/htmx.min.js")),
			Link(Href("/static/style.css"), Rel("stylesheet")),
			Link(Rel("icon"), Type("image/png"), Href("/static/timekeeper_icon.png")),
		),
		Body: append(children, Footer(Class("footer"),
			Div(
				components.PixelHackIcon("flag_pride", 24),
				components.PixelHackIcon("flag_trans", 24),
				components.PixelHackIcon("flag_nonbinary", 24),
			),
			Div(
				A(Text("Code of Conduct"), Href("https://jugendhackt.org/code-of-conduct/")),
				A(Text("Open Source"), Href("https://github.com/m4schini/timekeeper")),
				A(Text("Attributions"), Href("/help/attributions")),
				A(Text("Impressum"), Href("/help/legal")),
			),
			Div(
				components.PixelHackIcon("resitor_nonbinary", 24),
				components.PixelHackIcon("resistor_trans", 24),
				components.PixelHackIcon("resistor_pride", 24),
			),
		),
		),
		HTMLAttrs: nil,
	})
}
