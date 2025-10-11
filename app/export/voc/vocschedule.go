package voc

import (
	"encoding/json"
	"fmt"
	"time"
	"timekeeper/app/database/model"
	"timekeeper/config"
)

func ExportVocSchedule(event model.EventModel, timeslots []model.TimeslotModel) ([]byte, error) {
	conf := NewConference(fmt.Sprintf("timekeeper_event_%v", event.ID), event.Name, event.Start, event.TotalDays)

	tracksSet := make(map[model.Role]Track)

	for _, t := range timeslots {
		tracksSet[t.Role] = TrackFromRole(t.Role)

		day := t.Day
		room := t.Room.Name

		eventDate := t.Date()
		event := ConferenceEvent{
			Abstract:    t.Note,
			Description: t.Note,
			Date:        eventDate,
			Duration:    "01:00",
			Guid:        fmt.Sprintf("00000000-0000-0000-0000-%012d", t.ID),
			Id:          t.ID,
			Language:    "de",
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

	schedule := NewSchedule(event.URL(), fmt.Sprintf("0.0.%v", time.Now().Unix()), event.URL(), conf)
	return json.Marshal(schedule)
}
