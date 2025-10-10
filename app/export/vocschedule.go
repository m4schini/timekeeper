package export

import (
	"encoding/json"
	"fmt"
	"time"
	"timekeeper/app/database/model"
	"timekeeper/config"
)

func ExportVocSchedule(event model.EventModel, timeslots []model.TimeslotModel) ([]byte, error) {
	conferenceStart := event.Start.Format("2006-01-02")
	conferenceEnd := event.Start.AddDate(0, 0, event.TotalDays-1).Format("2006-01-02")

	// Build a map of days -> rooms -> []events
	daysMap := map[int]map[string][]map[string]interface{}{}

	for _, t := range timeslots {
		day := t.Day
		room := t.Room.Name
		if _, ok := daysMap[day]; !ok {
			daysMap[day] = map[string][]map[string]interface{}{}
		}
		eventJSON := map[string]interface{}{
			"abstract": "",
			"date":     t.Start.Format(time.RFC3339),
			"duration": "01:00",
			"guid":     fmt.Sprintf("00000000-0000-0000-0000-%012d", t.ID),
			"id":       t.ID,
			"language": "de",
			"links":    []interface{}{},
			"persons":  []interface{}{},
			"room":     room,
			"slug":     fmt.Sprintf("event-%d", t.ID),
			"start":    t.Start.Format("15:04"),
			"subtitle": "",
			"title":    t.Title,
			"track":    nil,
			"type":     "other",
			"url":      fmt.Sprintf("%v/event/%v#%v", config.BaseUrl(), event.ID, t.ID),
		}
		daysMap[day][room] = append(daysMap[day][room], eventJSON)
	}

	// Construct minimal days array
	var daysArr []map[string]interface{}
	for day, rooms := range daysMap {
		dayDate := event.Start.AddDate(0, 0, day)
		dayEntry := map[string]interface{}{
			"index":     day,
			"date":      dayDate.Format("2006-01-02"),
			"day_start": dayDate.Format(time.RFC3339),
			"day_end":   dayDate.Add(23 * time.Hour).Format(time.RFC3339),
			"rooms":     rooms,
		}
		daysArr = append(daysArr, dayEntry)
	}

	// Build final schedule
	schedule := map[string]interface{}{
		"schedule": map[string]interface{}{
			"version": "1.0",
			"conference": map[string]interface{}{
				"title":             event.Name,
				"acronym":           fmt.Sprintf("event_%d", event.ID),
				"daysCount":         event.TotalDays,
				"start":             conferenceStart,
				"end":               conferenceEnd,
				"timeslot_duration": "10:00",
				"days":              daysArr,
			},
		},
	}

	return json.MarshalIndent(schedule, "", "  ")
}
