package components

import (
	"fmt"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
	"slices"
	"sync"
)

var AvailablePixelHackIcons []string

func SetAvailablePixelHackIcons(pixelHackItems map[string]string) {
	names := make([]string, 0, len(pixelHackItems))
	for name, _ := range pixelHackItems {
		names = append(names, name)
	}
	slices.Sort(names)

	AvailablePixelHackIcons = names
}

func PixelHackIcon(name string, height int) Node {
	return Img(Src(fmt.Sprintf("/static/pixelhack/%v.svg", name)), Title(PixelHackAttribution()), Alt(name), Height(fmt.Sprintf("%vpx", height)))
}

func PixelHackAttribution() string {
	return "CC-BY-SA 4.0 Jugend Hackt (Hanno Sternberg)"
}

var generatePixelHackOpions = sync.OnceValue[Group](func() Group {
	g := make(Group, len(AvailablePixelHackIcons))
	for _, icon := range AvailablePixelHackIcons {
		g = append(g, Option(Style(fmt.Sprintf("background-image: url(/static/pixelhack/%v.svg); background-size: 24px", icon)), Text(icon), Value(icon)))
	}
	return g
})

func PixelHackSelectOptions() Node {
	return generatePixelHackOpions()
}

func PixelHackSelect() Node {

	return Div(
		Label(For("pixelhack"), Text("PixelHack Icon "), A(Text("(Liste)"), Href("/help/pixelhack"), Target("_blank"))),
		Select(Name("pixelhack"), PixelHackSelectOptions()),
	)
}
