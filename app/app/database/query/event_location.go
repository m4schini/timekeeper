package query

import . "raumzeitalpaka/app/database/model"

func (q *Queries) GetEventLocation(eventId, locationId int) (l EventLocationModel, err error) {
	row := q.DB.QueryRow(`
SELECT l.id, l.name, l.file, l.osm_id, ehl.name as relationship, ehl.id as relationship_id, ehl.note as relationship_note, ehl.visible as visible
FROM timekeeper.event_has_location ehl JOIN timekeeper.locations l ON l.id = ehl.location 
WHERE ehl.event = $1 AND ehl.location = $2
ORDER BY relationship, name`, eventId, locationId)
	if err := row.Err(); err != nil {
		return EventLocationModel{}, err
	}

	err = row.Scan(&l.ID, &l.Name, &l.File, &l.OsmId, &l.Relationship, &l.RelationshipId, &l.RelationshipNote, &l.Visible)
	if err != nil {
		return EventLocationModel{}, err
	}

	return l, nil
}
