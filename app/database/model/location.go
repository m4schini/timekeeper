package model

import "timekeeper/adapters"

type LocationModel struct {
	ID    int
	Name  string
	File  string
	OsmId string
}

type EventLocationModel struct {
	ID             int
	Name           string
	File           string
	OsmId          string
	Relationship   string
	RelationshipId int

	Address *adapters.OsmAddress
}

type CreateLocationModel struct {
	Name    string
	MapFile string
	OsmId   string
}

type AddLocationToEventModel struct {
	Name       string
	EventId    int
	LocationId int
}

type DeleteLocationFromEventModel struct {
	ID int
}
