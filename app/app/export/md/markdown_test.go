package md

import (
	"fmt"
	"raumzeitalpaka/app/database/model"
	"testing"
	"time"
)

func TestExportMarkdownTimetable(t *testing.T) {
	event := model.EventModel{
		ID:        1,
		Name:      "Test Event",
		TotalDays: 3,
		Start:     time.Now(),
	}
	room := model.RoomModel{
		ID: 1,
		Location: model.LocationModel{
			ID:   1,
			Name: "Location Name",
			File: "test.png",
		},
		Name:      "Testroom",
		LocationX: 0,
		LocationY: 0,
		LocationW: 0,
		LocationH: 0,
	}
	timeslots := []model.TimeslotModel{
		{
			ID:    1,
			Event: event,
			Title: "First Event",
			Note:  "Note",
			Day:   0,
			Start: time.Now().Add(1 * time.Hour),
			Room:  room,
			Role:  model.RoleParticipant,
		},
		{
			ID:    2,
			Event: event,
			Title: "Second Event",
			Day:   0,
			Start: time.Now().Add(2 * time.Hour),
			Room:  room,
			Role:  model.RoleMentor,
		},
		{
			ID:    3,
			Event: event,
			Title: "Third Event",
			Note:  "Note",
			Day:   0,
			Start: time.Now().Add(3 * time.Hour),
			Role:  model.RoleOrganizer,
		},
	}

	md, err := ExportMarkdownTimetable(timeslots)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	fmt.Println(md)
}
