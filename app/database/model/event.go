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
