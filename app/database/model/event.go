package model

import "time"

type EventModel struct {
	ID        int
	Name      string
	TotalDays int
	Start     time.Time
}
