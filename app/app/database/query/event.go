package query

import (
	"context"
	"raumzeitalpaka/app/database/model"
)

type GetEvent Handler[GetEventRequest, model.EventModel]

type GetEventRequest struct {
	EventId int `json:"eventId"`
}

type GetEventHandler struct {
	DB Database
}

func (q *GetEventHandler) Query(ctx context.Context, request GetEventRequest) (e model.EventModel, err error) {
	id := request.EventId
	row := q.DB.QueryRow(`SELECT id, name, event_start, event_end, slug, guid FROM raumzeitalpaka.events e WHERE id = $1`, id)
	if err = row.Err(); err != nil {
		return model.EventModel{}, nil
	}

	err = row.Scan(&e.ID, &e.Name, &e.Start, &e.End, &e.Slug, &e.GUID)
	if err != nil {
		return e, err
	}
	e.CalculateTotalDays()
	return e, nil
}
