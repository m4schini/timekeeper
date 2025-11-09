package voc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"time"
	"timekeeper/app/database/model"
	"timekeeper/config"

	"go.uber.org/zap"
)

func ExportVocScheduleTo(event model.EventModel, timeslots []model.TimeslotModel, writer io.Writer) error {
	log := zap.L().Named("export").Named("vocschedule").With(zap.Int("event", event.ID), zap.Int("timeslots_count", len(timeslots)))
	log.Debug("exporting schedule as voc-schedule")

	conf := NewConference(fmt.Sprintf("timekeeper_event_%v", event.ID), event.Name, event.Start, event.TotalDays)
	tracksSet := make(map[model.Role]Track)

	timeslots = model.FlattenTimeslots(timeslots)

	for _, t := range timeslots {
		tracksSet[t.Role] = TrackFromRole(t.Role)

		day := t.Day
		room := t.Room.Name

		hours := int(t.Duration.Hours())
		minutes := int(t.Duration.Minutes()) % 60

		t.Note = config.PixelHackPlaceholderRx.ReplaceAllString(t.Note, "")
		t.Title = config.PixelHackPlaceholderRx.ReplaceAllString(t.Title, "")

		eventDate := t.Date()
		event := ConferenceEvent{
			Abstract:    t.Note,
			Description: t.Note,
			Date:        eventDate,
			Duration:    fmt.Sprintf("%02d:%02d", hours, minutes),
			Guid:        t.GUID,
			Id:          t.ID,
			Language:    "de",
			Track:       t.Role.Title(),
			Room:        room,
			Slug:        fmt.Sprintf("event-%d", t.ID),
			Start:       t.Start.Format("15:04"),
			Subtitle:    t.Note,
			Title:       t.Title,
			Type:        "other",
			Url:         fmt.Sprintf("%v/event/%v#%v", config.BaseUrl(), event.ID, t.ID),
			Links:       make([]interface{}, 0),
			Persons:     make([]Person, 0),
		}
		conf.Days[day].Rooms[room] = append(conf.Days[day].Rooms[room], event)
	}

	tracks := make([]Track, 0, len(tracksSet))
	for _, track := range tracksSet {
		tracks = append(tracks, track)
	}

	conf.Tracks = tracks

	log.Debug("exported schedule as voc-schedule")
	schedule := NewSchedule(event.EventURL(), fmt.Sprintf("0.0.%v", time.Now().Unix()), event.EventURL(), conf)
	return json.NewEncoder(writer).Encode(schedule)
}

func ExportVocSchedule(event model.EventModel, timeslots []model.TimeslotModel) ([]byte, error) {
	var buf bytes.Buffer
	err := ExportVocScheduleTo(event, timeslots, &buf)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
