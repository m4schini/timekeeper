package query

import . "raumzeitalpaka/app/database/model"

func (q *Queries) GetLocation(id int) (l LocationModel, err error) {
	row := q.DB.QueryRow(`SELECT id, name, file, osm_id FROM timekeeper.locations WHERE id = $1`, id)
	if err = row.Err(); err != nil {
		return LocationModel{}, err
	}

	err = row.Scan(&l.ID, &l.Name, &l.File, &l.OsmId)
	return l, err
}
