package query

import . "timekeeper/app/database/model"

func (q *Queries) GetLocationsOfEvent(eventId int) (ls []EventLocationModel, err error) {
	rows, err := q.DB.Query(`
SELECT l.id, l.name, l.file, l.osm_id, ehl.name as relationship, ehl.id as relationship_id, ehl.note as relationship_note
FROM timekeeper.event_has_location ehl JOIN timekeeper.locations l ON l.id = ehl.location 
WHERE ehl.event = $1
ORDER BY relationship`, eventId)
	if err != nil {
		return nil, err
	}

	ls = make([]EventLocationModel, 0)
	for rows.Next() {
		var l EventLocationModel
		err = rows.Scan(&l.ID, &l.Name, &l.File, &l.OsmId, &l.Relationship, &l.RelationshipId, &l.RelationshipNote)
		if err != nil {
			return nil, err
		}

		ls = append(ls, l)
	}

	return ls, nil
}
