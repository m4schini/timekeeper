package query

import . "timekeeper/app/database/model"

func (q *Queries) GetLocations(offset, limit int) (ls []LocationModel, err error) {
	rows, err := q.DB.Query(`SELECT id, name, file FROM timekeeper.locations LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, err
	}

	ls = make([]LocationModel, 0, limit)
	for rows.Next() {
		var l LocationModel
		err = rows.Scan(&l.ID, &l.Name, &l.File)
		if err != nil {
			return nil, err
		}

		ls = append(ls, l)
	}

	return ls, nil
}
