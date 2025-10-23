package voc

import (
	"strings"
	"time"
	"timekeeper/app/database/model"
	"timekeeper/config"
)

const (
	ScheduleSchema   = "https://c3voc.de/schedule/schema.json"
	GeneratorName    = "timekeeper"
	GeneratorVersion = "dev"
)

func NewSchedule(url, version, baseUrl string, conf Conference) Schedule {
	return Schedule{
		Schema:    ScheduleSchema,
		Generator: Generator{Name: GeneratorName, Version: GeneratorVersion},
		Schedule: EventSchedule{
			Url:        url,
			Version:    version,
			BaseUrl:    baseUrl,
			Conference: conf,
		},
	}
}

func NewConference(acronym, title string, start time.Time, daysCount int) Conference {
	days := make([]Day, daysCount)
	for i := range daysCount {
		start := start.AddDate(0, 0, i)
		startTime := time.Date(start.Year(), start.Month(), start.Day(), 6, 0, 0, 0, config.Timezone())
		endTime := time.Date(start.Year(), start.Month(), start.Day(), 23, 0, 0, 0, config.Timezone())

		days[i] = Day{
			Index:     i + 1,
			Date:      start,
			DateStart: startTime,
			DateEnd:   endTime,
			Rooms:     make(map[string][]ConferenceEvent),
		}
	}

	return Conference{
		Acronym:          acronym,
		Title:            title,
		Start:            start,
		End:              start.AddDate(0, 0, daysCount),
		DaysCount:        daysCount,
		TimeslotDuration: 1 * time.Hour,
		TimeZoneName:     config.Timezone().String(),
		Colors:           map[string]string{"primary": "#000000"},
		Rooms:            nil,
		Days:             days,
	}
}

func TrackFromRole(role model.Role) Track {
	return Track{
		Name:  role.Title(),
		Slug:  strings.ToLower(string(role)),
		Color: role.Color(),
	}
}
