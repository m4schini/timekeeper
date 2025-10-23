package voc

import (
	"fmt"
	"testing"
	"time"
	"timekeeper/app/database/model"
)

func TestExportVocSchedule(t *testing.T) {
	room := model.RoomModel{
		ID: 1,
		Location: model.LocationModel{
			ID:   2,
			Name: "LocName",
		},
		Name: "Testraum",
	}
	event := model.EventModel{
		ID:        3,
		Name:      "JHHH25",
		TotalDays: 3,
		Start:     time.Now(),
	}
	timeSlots := []model.TimeslotModel{
		{
			ID:    4,
			Event: event,
			Title: "Timeslot 1",
			Note:  "Test",
			Day:   0,
			Start: time.Now().Add(5 * time.Minute),
			Room:  room,
			Role:  model.RoleParticipant,
		},
	}
	j, err := ExportVocSchedule(event, timeSlots)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	fmt.Println(string(j))
}
