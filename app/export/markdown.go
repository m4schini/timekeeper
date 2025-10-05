package export

import (
	"bytes"
	"github.com/JohannesKaufmann/html-to-markdown/v2/plugin/base"
	"github.com/JohannesKaufmann/html-to-markdown/v2/plugin/commonmark"
	"github.com/JohannesKaufmann/html-to-markdown/v2/plugin/table"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
	"timekeeper/app/database/model"

	"github.com/JohannesKaufmann/html-to-markdown/v2/converter"
)

var conv = converter.NewConverter(
	converter.WithPlugins(
		base.NewBasePlugin(),
		table.NewTablePlugin(),
		commonmark.NewCommonmarkPlugin(
			commonmark.WithStrongDelimiter("__"),
			// ...additional configurations for the plugin
		),
	),
)

func ExportMarkdownTimetable(timeslots []model.TimeslotModel) (string, error) {
	rows := Group{}
	for _, timeslot := range timeslots {
		rows = append(rows, Tr(
			Td(Text(timeslot.Start.Format("15:04"))),
			Td(Text(timeslot.Room.Name)),
			Td(Text(timeslot.Title)),
			Td(Text(timeslot.Note)),
		))
	}

	node := Table(
		THead(Tr(
			Th(Text("Zeit")),
			Th(Text("Raum")),
			Th(Text("Name")),
			Th(Text("Notes")),
		)),
		TBody(rows),
	)

	var buf bytes.Buffer
	err := node.Render(&buf)
	if err != nil {
		return "", err
	}
	html := string(buf.Bytes())

	markdown, err := conv.ConvertString(html)
	if err != nil {
		return "", err
	}

	return markdown, nil
}
