package query

import (
	"raumzeitalpaka/app/database/model"
)

type GetEventLocations Handler[GetEventLocationsRequest, []model.EventLocationModel]

type GetEventLocationsRequest struct {
	EventId int
}

type GetEventLocationsHandler struct {
	DB Database
}

func (q *GetEventLocationsHandler) Query(request GetEventLocationsRequest) (ls []model.EventLocationModel, err error) {
	eventId := request.EventId
	rows, err := q.DB.Query(`
SELECT l.id, l.name, l.file, l.osm_id, ehl.name as relationship, ehl.id as relationship_id, ehl.note as relationship_note, ehl.visible as visible
FROM raumzeitalpaka.event_has_location ehl JOIN raumzeitalpaka.locations l ON l.id = ehl.location 
WHERE ehl.event = $1
ORDER BY relationship, name`, eventId)
	if err != nil {
		return nil, err
	}

	ls = make([]model.EventLocationModel, 0)
	for rows.Next() {
		var l model.EventLocationModel
		err = rows.Scan(&l.ID, &l.Name, &l.File, &l.OsmId, &l.Relationship, &l.RelationshipId, &l.RelationshipNote, &l.Visible)
		if err != nil {
			return nil, err
		}

		ls = append(ls, l)
	}

	return ls, nil
}
