package query

import (
	"database/sql"
	"raumzeitalpaka/app/database/model"
)

type GetEvent Handler[GetEventRequest, model.EventModel]

type GetEventRequest struct {
	EventId int
}

type GetEventHandler struct {
	DB Database
}

func NewGetEventHandler(db *sql.DB) *GetEventHandler {
	return &GetEventHandler{DB: db}
}

func (q *GetEventHandler) Query(request GetEventRequest) (e model.EventModel, err error) {
	id := request.EventId
	row := q.DB.QueryRow(`SELECT id, name, start, slug, guid FROM raumzeitalpaka.events WHERE id = $1`, id)
	if err = row.Err(); err != nil {
		return model.EventModel{}, nil
	}

	row2 := q.DB.QueryRow(`SELECT count(day) FROM (SELECT DISTINCT day FROM raumzeitalpaka.timeslots WHERE event = $1) AS day`, id)
	if err = row.Err(); err != nil {
		return model.EventModel{}, nil
	}

	var totalDays int
	err = row2.Scan(&totalDays)
	if err != nil {
		return model.EventModel{}, err
	}

	err = row.Scan(&e.ID, &e.Name, &e.Start, &e.Slug, &e.GUID)
	e.TotalDays = totalDays
	return e, nil
}
