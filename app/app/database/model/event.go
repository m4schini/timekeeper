package model

import (
	"fmt"
	"raumzeitalpaka/config"
	"time"
)

type EventModel struct {
	ID        int
	GUID      string
	Slug      string
	Name      string
	TotalDays int
	Start     time.Time
}

func (e EventModel) EventURL() string {
	return fmt.Sprintf("%v/event/%v", config.BaseUrl(), e.ID)
}

func (e EventModel) ScheduleURL() string {
	return fmt.Sprintf("%v/event/%v/schedule", config.BaseUrl(), e.ID)
}

func (e EventModel) Day(day int) time.Time {
	return e.Start.AddDate(0, 0, day)
}

type CreateEventModel struct {
	Name  string
	Slug  string
	Start time.Time
}

type UpdateEventModel struct {
	ID    int
	Name  string
	Slug  string
	Start time.Time
}
