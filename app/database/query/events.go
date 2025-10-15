package query

import "timekeeper/app/database/model"

func (q *Queries) GetEvents(offset, limit int) (es []model.EventModel, err error) {
	rows, err := q.DB.Query(`
SELECT id, name, start 
FROM timekeeper.events
ORDER BY start, name
LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, err
	}

	es = make([]model.EventModel, 0, limit)
	for rows.Next() {
		var e model.EventModel
		err = rows.Scan(&e.ID, &e.Name, &e.Start)
		if err != nil {
			return nil, err
		}

		es = append(es, e)
	}

	return es, nil
}
