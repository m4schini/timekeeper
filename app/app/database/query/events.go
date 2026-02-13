package query

import (
	"raumzeitalpaka/app/database/model"
)

type GetEvents Handler[GetEventsRequest, []model.EventModel]

type GetEventsRequest struct {
	Offset int
	Limit  int
}

type GetEventsHandler struct {
	DB Database
}

func (q *GetEventsHandler) Query(request GetEventsRequest) (es []model.EventModel, err error) {
	offset := request.Offset
	limit := request.Limit
	rows, err := q.DB.Query(`
	SELECT 
		e.id, 
		e.name, 
		e.start, 
		e.slug,
		COUNT(DISTINCT t.day) AS total_days
	FROM 
		raumzeitalpaka.events e
	LEFT JOIN 
		raumzeitalpaka.timeslots t 
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
		err := rows.Scan(&r.ID, &r.Name, &r.Start, &r.Slug, &r.TotalDays)
		if err != nil {
			return nil, err
		}

		es = append(es, r)
	}
	return es, nil
}
