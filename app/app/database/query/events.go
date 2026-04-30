package query

import (
	"context"
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

func (q *GetEventsHandler) Query(ctx context.Context, request GetEventsRequest) (es []model.EventModel, err error) {
	offset := request.Offset
	limit := request.Limit
	rows, err := q.DB.Query(`
	SELECT 
		e.id, 
		e.guid,
		e.name, 
		e.event_start, 
		e.event_end,
		e.slug,
		COUNT(DISTINCT t.day) AS total_days
	FROM 
		raumzeitalpaka.events e
	LEFT JOIN 
		raumzeitalpaka.timeslots t 
	ON 
		e.id = t.event
	GROUP BY 
		e.id, e.name, e.event_start
	OFFSET $1 LIMIT $2 
`, offset, limit)
	if err != nil {
		return nil, err
	}

	es = make([]model.EventModel, 0)
	for rows.Next() {
		var r model.EventModel
		err := rows.Scan(&r.ID, &r.GUID, &r.Name, &r.Start, &r.End, &r.Slug, &r.TotalDays)
		if err != nil {
			return nil, err
		}

		r.EventDays()
		es = append(es, r)
	}
	return es, nil
}
