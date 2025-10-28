package query

import . "timekeeper/app/database/model"

func (q *Queries) GetEvent(id int) (r EventModel, err error) {
	row := q.DB.QueryRow(`SELECT id, name, start, slug FROM timekeeper.events WHERE id = $1`, id)
	if err = row.Err(); err != nil {
		return EventModel{}, err
	}

	row2 := q.DB.QueryRow(`SELECT count(day) FROM (SELECT DISTINCT day FROM timekeeper.timeslots WHERE event = $1) AS day`, id)
	if err = row.Err(); err != nil {
		return EventModel{}, err
	}

	var totalDays int
	err = row2.Scan(&totalDays)
	if err != nil {
		return EventModel{}, err
	}

	err = row.Scan(&r.ID, &r.Name, &r.Start, &r.Slug)
	r.TotalDays = totalDays
	return r, nil
}
