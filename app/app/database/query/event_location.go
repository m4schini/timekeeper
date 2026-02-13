package query

import (
	"raumzeitalpaka/app/database/model"
)

type GetEventLocation Handler[GetEventLocationRequest, model.EventLocationModel]

type GetEventLocationRequest struct {
	EventId    int
	LocationId int
}

type GetEventLocationHandler struct {
	DB Database
}

func (q *GetEventLocationHandler) Query(request GetEventLocationRequest) (l model.EventLocationModel, err error) {
	eventId := request.EventId
	locationId := request.LocationId
	row := q.DB.QueryRow(`
SELECT l.id, l.name, l.file, l.osm_id, ehl.name as relationship, ehl.id as relationship_id, ehl.note as relationship_note, ehl.visible as visible
FROM raumzeitalpaka.event_has_location ehl JOIN raumzeitalpaka.locations l ON l.id = ehl.location 
WHERE ehl.event = $1 AND ehl.location = $2
ORDER BY relationship, name`, eventId, locationId)
	if err := row.Err(); err != nil {
		return model.EventLocationModel{}, err
	}

	err = row.Scan(&l.ID, &l.Name, &l.File, &l.OsmId, &l.Relationship, &l.RelationshipId, &l.RelationshipNote, &l.Visible)
	if err != nil {
		return model.EventLocationModel{}, err
	}

	return l, nil
}
