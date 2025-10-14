package model

import "timekeeper/adapters"

type LocationModel struct {
	ID    int
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

	Address *adapters.OsmAddress
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
	ID   int
	Name string
	Note string
}

type DeleteLocationFromEventModel struct {
	ID int
}
