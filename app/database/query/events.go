package query

import "timekeeper/app/database/model"

func (q *Queries) GetEvents(offset, limit int) (es []model.EventModel, err error) {
	rows, err := q.DB.Query(`
	SELECT 
		e.id, 
		e.name, 
		e.start, 
		COUNT(DISTINCT t.day) AS total_days
	FROM 
		timekeeper.events e
	LEFT JOIN 
		timekeeper.timeslots t 
	ON 
		e.id = t.event
	GROUP BY 
		e.id, e.name, e.start
	OFFSET $1 LIMIT $2 
`, offset, limit)
	if err != nil {
		return nil, err
	}

	es = make([]model.EventModel, 0)
	for rows.Next() {
		var r model.EventModel
		err := rows.Scan(&r.ID, &r.Name, &r.Start, &r.TotalDays)
		if err != nil {
			return nil, err
		}

		es = append(es, r)
	}

	return es, nil
}
