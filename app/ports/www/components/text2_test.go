package components

import "testing"

func TestPixelHackPlaceholderRx(t *testing.T) {
	in := `:apple: :birb:`

	text := pixelhackPlaceholderRx.ReplaceAllString(in, `<img src="/static/pixelhack/$1.svg" style="height: 0.5rem" />`)

	t.Log(in)
	t.Log(text)
}
