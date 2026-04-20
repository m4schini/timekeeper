package model

import (
	"fmt"
	"raumzeitalpaka/config"
	"time"
)

type EventModel struct {
	ID        int       `json:"id"`
	GUID      string    `json:"guid"`
	Slug      string    `json:"slug"`
	Name      string    `json:"name"`
	TotalDays int       `json:"totalDays"`
	Start     time.Time `json:"start"`
	End       time.Time `json:"end"`
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

func (e *EventModel) CalculateTotalDays() int {
	e.TotalDays = (int(e.End.Sub(e.Start).Hours()) / 24) + 1
	return e.TotalDays
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
