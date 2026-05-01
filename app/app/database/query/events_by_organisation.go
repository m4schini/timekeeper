package query

import (
	"context"
	"raumzeitalpaka/app/database/model"
)

type GetEventsByOrganisation Handler[GetEventsByOrganisationRequest, []model.EventModel]

type GetEventsByOrganisationRequest struct {
	OrganisationID int
}

type GetEventsByOrganisationHandler struct {
	DB Database
}

func (q *GetEventsByOrganisationHandler) Query(ctx context.Context, request GetEventsByOrganisationRequest) (events []model.EventModel, err error) {
	rows, err := q.DB.Query(`
SELECT id, guid, name, slug, event_start, event_end, setup, teardown
FROM raumzeitalpaka.events 
WHERE organisation_id = $1`, request.OrganisationID)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var event model.EventModel
		err = rows.Scan(&event.ID, &event.GUID, &event.Name, &event.Slug, &event.Start, &event.End, &event.Setup, &event.Teardown)
		if err != nil {
			return nil, err
		}
		event.EventDays()
		events = append(events, event)
	}

	return events, err
}
