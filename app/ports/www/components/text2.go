package components

import (
	"fmt"
	"timekeeper/config"

	. "maragu.dev/gomponents"
)

func Text2(text string, pixelhackHeight int) Node {
	text = config.PixelHackPlaceholderRx.ReplaceAllString(text, fmt.Sprintf(`<img src="/static/pixelhack/$1.svg" style="height: %vpx" title="$1" />`, pixelhackHeight))

	return Raw(text)
}
