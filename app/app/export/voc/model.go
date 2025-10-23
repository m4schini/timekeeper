package voc

import (
	"encoding/json"
	"fmt"
	"time"
)

type Schedule struct {
	Schema    string        `json:"$schema"`
	Generator Generator     `json:"generator"`
	Schedule  EventSchedule `json:"schedule"`
}

type Generator struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type EventSchedule struct {
	Url        string     `json:"url"`
	Version    string     `json:"version"`
	BaseUrl    string     `json:"baseUrl"`
	Conference Conference `json:"conference"`
}

type Conference struct {
	Acronym          string
	Title            string
	Start            time.Time
	End              time.Time
	DaysCount        int
	TimeslotDuration time.Duration
	TimeZoneName     string
	Colors           ConferenceColors
	Rooms            []Room
	Tracks           []Track
	Days             []Day
}

func (c Conference) MarshalJSON() ([]byte, error) {
	rooms := c.Rooms
	if rooms == nil {
		rooms = make([]Room, 0)
	}

	encodeModel := conference{
		Acronym:          c.Acronym,
		Title:            c.Title,
		Start:            c.Start.Format("2006-01-02"),
		End:              c.End.Format("2006-01-02"),
		DaysCount:        c.DaysCount,
		TimeslotDuration: fmt.Sprintf("%02.f:%02.f", c.TimeslotDuration.Hours(), float64(0)),
		TimeZoneName:     c.TimeZoneName,
		Colors:           c.Colors,
		Rooms:            rooms,
		Tracks:           c.Tracks,
		Days:             c.Days,
	}

	return json.Marshal(encodeModel)
}

type conference struct {
	Acronym          string           `json:"acronym"`
	Title            string           `json:"title"`
	Start            string           `json:"start"`
	End              string           `json:"end"`
	DaysCount        int              `json:"daysCount"`
	TimeslotDuration string           `json:"timeslot_duration"` //HH:mm
	TimeZoneName     string           `json:"time_zone_name"`
	Colors           ConferenceColors `json:"colors"`
	Rooms            []Room           `json:"rooms"`
	Tracks           []Track          `json:"tracks"`
	Days             []Day            `json:"days"`
}

type ConferenceColors map[string]string

type Room struct {
	Name        string  `json:"name"`
	Slug        string  `json:"slug"`
	Guid        string  `json:"guid"`
	Description *string `json:"description"`
	Capacity    *int    `json:"capacity"`
}

type Track struct {
	Name  string `json:"name"`
	Slug  string `json:"slug"`
	Color string `json:"color"`
}

type Day struct {
	Index     int
	Date      time.Time
	DateStart time.Time
	DateEnd   time.Time
	Rooms     map[string][]ConferenceEvent
}

func (d Day) MarshalJSON() ([]byte, error) {
	encodeModel := day{
		Index:     d.Index,
		Date:      d.Date.Format("2006-01-02"),
		DateStart: d.DateStart.Format("2006-01-02T15:04:05-07:00"),
		DateEnd:   d.DateStart.Format("2006-01-02T15:04:05-07:00"),
		Rooms:     d.Rooms,
	}
	return json.Marshal(encodeModel)
}

type day struct {
	Index     int                          `json:"index"`
	Date      string                       `json:"date"`
	DateStart string                       `json:"day_start"`
	DateEnd   string                       `json:"day_end"`
	Rooms     map[string][]ConferenceEvent `json:"rooms"`
}

type DayRoom struct {
	Guid        string        `json:"guid"`
	Code        string        `json:"code"`
	Id          int           `json:"id"`
	Logo        interface{}   `json:"logo"`
	Date        time.Time     `json:"date"`
	Start       string        `json:"start"`
	Duration    string        `json:"duration"`
	Room        string        `json:"room"`
	Slug        string        `json:"slug"`
	Url         string        `json:"url"`
	Title       string        `json:"title"`
	Subtitle    string        `json:"subtitle"`
	Track       string        `json:"track"`
	Type        string        `json:"type"`
	Language    string        `json:"language"`
	Description string        `json:"description"`
	Links       []interface{} `json:"links"`
	FeedbackUrl string        `json:"feedback_url"`
	OriginUrl   string        `json:"origin_url"`
	Attachments []interface{} `json:"attachments"`
}

type Person struct {
	Code       string      `json:"code"`
	Name       string      `json:"name"`
	Avatar     interface{} `json:"avatar"`
	Biography  interface{} `json:"biography"`
	PublicName string      `json:"public_name"`
	Guid       string      `json:"guid"`
	Url        string      `json:"url"`
}

type ConferenceEvent struct {
	Guid             string        `json:"guid"`
	Code             string        `json:"code,omitempty"`
	Id               int           `json:"id"`
	Logo             interface{}   `json:"logo"`
	Date             time.Time     `json:"date"`
	Start            string        `json:"start"`
	Duration         string        `json:"duration"`
	Room             string        `json:"room"`
	Slug             string        `json:"slug"`
	Url              string        `json:"url"`
	Title            string        `json:"title"`
	Subtitle         string        `json:"subtitle"`
	Track            string        `json:"track"`
	Type             string        `json:"type"`
	Language         string        `json:"language"`
	Abstract         string        `json:"abstract"`
	Description      string        `json:"description"`
	RecordingLicense string        `json:"recording_license"`
	DoNotRecord      bool          `json:"do_not_record"`
	Persons          []Person      `json:"persons"`
	Links            []interface{} `json:"links"`
	FeedbackUrl      string        `json:"feedback_url,omitempty"`
	OriginUrl        string        `json:"origin_url,omitempty"`
	Attachments      []interface{} `json:"attachments,omitempty"`
}
