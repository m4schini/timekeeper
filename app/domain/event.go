package domain

import (
	"context"
	"fmt"
	"raumzeitalpaka/app/auth/user"
	"raumzeitalpaka/app/database/command"
	"raumzeitalpaka/app/database/model"
	"raumzeitalpaka/app/database/query"
	"time"
)

type InvalidEventDateErr struct {
	Start time.Time
	End   time.Time
}

func (i InvalidEventDateErr) Error() string {
	return fmt.Sprintf("invalid event date: %v-%v", i.Start, i.End)
}

func CreateEvent(ctx context.Context, createEvent command.CreateEvent, request command.CreateEventRequest) (int, error) {
	_, authenticated := user.IdentityFrom(ctx)
	if !authenticated {
		return -1, UnauthenticatedErr{}
	}

	// valid event dates: start needs to be before end
	if request.Start.After(request.End) {
		return -1, InvalidEventDateErr{
			Start: request.Start,
			End:   request.End,
		}
	}

	return createEvent.Execute(ctx, request)
}

func GetEvent(ctx context.Context, getEvent query.GetEvent, request query.GetEventRequest) (model.EventModel, error) {
	return getEvent.Query(ctx, request)
}

func GetEvents(ctx context.Context, getEvents query.GetEvents) ([]model.EventModel, error) {
	return getEvents.Query(ctx, query.GetEventsRequest{
		Offset: 0,
		Limit:  100,
	})
}

func UpdateEvent(ctx context.Context, hasRole query.UserHasRole, updateEvent command.UpdateEvent, request command.UpdateEventRequest) error {
	// authenticate
	user, authenticated := user.IdentityFrom(ctx)
	if !authenticated {
		return UnauthenticatedErr{}
	}

	// authorize
	hasRoleResponse, err := hasRole.Query(ctx, query.UserHasRoleRequest{
		UserId: user.User,
		Role:   model.RoleOrganizer,
	})
	if err != nil {
		return fmt.Errorf("%w: %w", UnauthorizedErr{UserId: user.User, Action: "UPDATE", Subject: request.ID}, err)
	}
	if !hasRoleResponse.HasRole {
		return UnauthorizedErr{UserId: user.User, Action: "UPDATE", Subject: request.ID}
	}

	return updateEvent.Execute(ctx, request)
}
