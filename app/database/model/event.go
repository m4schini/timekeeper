package model

import (
	"fmt"
	"time"
	"timekeeper/config"
)

type EventModel struct {
	ID        int
	Name      string
	TotalDays int
	Start     time.Time
}

func (e EventModel) URL() string {
	return fmt.Sprintf("%v/event/%v", config.BaseUrl(), e.ID)
}

func (e EventModel) Day(day int) time.Time {
	return e.Start.AddDate(0, 0, day)
}

type CreateEventModel struct {
	Name  string
	Start time.Time
}

type UpdateEventModel struct {
	ID    int
	Name  string
	Start time.Time
}
