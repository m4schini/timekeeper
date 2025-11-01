package components

import (
	"fmt"
	"regexp"

	. "maragu.dev/gomponents"
)

var pixelhackPlaceholderRx = regexp.MustCompile(`:([a-z0-9_]+):`)

func Text2(text string, pixelhackHeight int) Node {
	text = pixelhackPlaceholderRx.ReplaceAllString(text, fmt.Sprintf(`<img src="/static/pixelhack/$1.svg" style="height: %vpx" title="$1" />`, pixelhackHeight))

	return Raw(text)
}
