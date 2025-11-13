package model

import (
	"raumzeitalpaka/adapters/nominatim"
	"strings"
)

type LocationModel struct {
	ID    int
	GUID  string
	Name  string
	File  string
	OsmId string
}

type EventLocationModel struct {
	ID               int
	Name             string
	File             string
	OsmId            string
	Relationship     string
	RelationshipId   int
	RelationshipNote string
	Visible          bool

	Address *nominatim.Address
}

func (e EventLocationModel) RelationshipLabel() (label string) {
	el := strings.TrimSpace(strings.ToLower(e.Relationship))
	switch el {
	case "sleep_location":
		label = "Übernachtungsort"
		break
	case "event_location":
		label = "Eventort"
		break
	default:
		label = el
		break
	}

	return label
}

type CreateLocationModel struct {
	Name    string
	MapFile string
	OsmId   string
}

type UpdateLocationModel struct {
	ID      int
	Name    string
	MapFile string
	OsmId   string
}

type AddLocationToEventModel struct {
	Name       string
	EventId    int
	LocationId int
	Note       string
}

type UpdateLocationToEventModel struct {
	ID      int
	Name    string
	Note    string
	Visible bool
}

type DeleteLocationFromEventModel struct {
	ID int
}
