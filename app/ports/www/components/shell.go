package components

import (
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/components"
	. "maragu.dev/gomponents/html"
)

func Shell(title string, children ...Node) Node {
	if title == "" {
		title = "RaumZeitAlpaka"
	}
	return ShellWithHead(title, PageFooter(), []Node{}, children...)
}

func ShellWithHead(title string, footer Node, head []Node, children ...Node) Node {
	if title == "" {
		title = "Timekeeper"
	}
	return HTML5(HTML5Props{
		Title:       title,
		Description: "Zeitplan",
		Language:    "de",
		Head: append(head,
			Script(Src("/static/htmx.min.js")),
			Link(Href("/static/style2.css"), Rel("stylesheet")),
			Link(Rel("icon"), Type("image/png"), Href("/static/timekeeper_icon.png")),
		),
		Body:      append(children, footer),
		HTMLAttrs: nil,
	})
}

func PageFooter() Node {
	return Footer(Class("footer"),
		Div(
			PixelHackIcon("flag_pride", 24),
			PixelHackIcon("flag_trans", 24),
			PixelHackIcon("flag_nonbinary", 24),
		),
		Div(
			A(Text("Code of Conduct"), Href("https://jugendhackt.org/code-of-conduct/")),
			A(Text("Open Source"), Href("https://codeberg.org/aur0ra/raumzeitalpaka")),
			A(Text("Report a bug"), Href("https://codeberg.org/aur0ra/raumzeitalpaka/issues/new?template=.github%2fISSUE_TEMPLATE%2fbug_report.md")),
			A(Text("Attributions"), Href("/help/attributions")),
			A(Text("Impressum"), Href("/help/legal")),
		),
		Div(
			PixelHackIcon("resitor_nonbinary", 24),
			PixelHackIcon("resistor_trans", 24),
			PixelHackIcon("resistor_pride", 24),
		),
	)
}
