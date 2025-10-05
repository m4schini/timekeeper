package components

import (
	"fmt"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

const locationMapSrc = "/static/betahaus.png"

func LocationMap() Node {
	return Div(Class("location-map"),
		Img(Src(locationMapSrc)),
	)
}

func LocationCrop(x, y, width, height, targetH int) Node {
	scale := float64(targetH) / float64(height)
	targetW := int(scale * float64(width))

	return Div(
		Style(fmt.Sprintf("width:%dpx; height:%dpx; overflow:hidden;", targetW, targetH)),
		Div(
			Style(fmt.Sprintf(`
				width:%dpx;
				height:%dpx;
				background-image:url('%s');
				background-position:-%dpx -%dpx;
				background-repeat:no-repeat;
				transform:scale(%f);
				transform-origin:top left;
			`,
				width, height, locationMapSrc, x, y, scale,
			)),
		),
	)

}
