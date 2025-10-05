package query

import . "timekeeper/app/database/model"

func (q *Queries) GetLocation(id int) (l LocationModel, err error) {
	row := q.DB.QueryRow(`SELECT id, name, file FROM locations WHERE id = $1`, id)
	if err = row.Err(); err != nil {
		return LocationModel{}, err
	}

	err = row.Scan(&l.ID, &l.Name, &l.File)
	return l, err
}
