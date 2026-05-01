package api

import (
	"encoding/json"
	"net/http"
	"raumzeitalpaka/app/database"
	"raumzeitalpaka/app/database/command"
	"raumzeitalpaka/app/database/query"
)

func NewRouter(db *database.Database) http.Handler {
	c := db.Commands
	q := db.Queries

	return HandlerWithOptions(&apiServer{
		c: c,
		q: q,
	}, ChiServerOptions{})
}

type apiServer struct {
	c command.Commands
	q query.Queries
}

func (a *apiServer) CreateEvent(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}

func (a *apiServer) RemoveEventLocation(w http.ResponseWriter, r *http.Request, locationID int) {
	//TODO implement me
	panic("implement me")
}

func (a *apiServer) AddEventLocation(w http.ResponseWriter, r *http.Request, locationID int) {
	//TODO implement me
	panic("implement me")
}

func (a *apiServer) GetEvent(w http.ResponseWriter, r *http.Request, eventID int) {
	event, err := a.q.Event.Query(r.Context(), query.GetEventRequest{EventId: eventID})
	if err != nil {
		internalServerError(w, err)
		return
	}

	encode(w, event)
}

func (a *apiServer) UpdateEvent(w http.ResponseWriter, r *http.Request, eventID int) {
	//TODO implement me
	panic("implement me")
}

func (a *apiServer) GetEventRoomSchedule(w http.ResponseWriter, r *http.Request, eventID float32, roomID float32) {
	//TODO implement me
	panic("implement me")
}

func (a *apiServer) GetEventSchedule(w http.ResponseWriter, r *http.Request, eventID int) {
	//TODO implement me
	panic("implement me")
}

func (a *apiServer) CreateTimeslot(w http.ResponseWriter, r *http.Request, eventID int) {
	//TODO implement me
	panic("implement me")
}

func (a *apiServer) GetTimeslot(w http.ResponseWriter, r *http.Request, eventID int, timeslotID int) {
	//TODO implement me
	panic("implement me")
}

func (a *apiServer) UpdateTimeslot(w http.ResponseWriter, r *http.Request, eventID int, timeslotID int) {
	//TODO implement me
	panic("implement me")
}

func (a *apiServer) GetEventUserSchedule(w http.ResponseWriter, r *http.Request, eventID float32, userID float32) {
	//TODO implement me
	panic("implement me")
}

func (a *apiServer) GetEvents(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}

func (a *apiServer) CreateLocation(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}

func (a *apiServer) GetLocation(w http.ResponseWriter, r *http.Request, locationID int) {
	//TODO implement me
	panic("implement me")
}

func (a *apiServer) UpdateLocation(w http.ResponseWriter, r *http.Request, locationID int) {
	//TODO implement me
	panic("implement me")
}

func (a *apiServer) CreateRoom(w http.ResponseWriter, r *http.Request, locationID float32) {
	//TODO implement me
	panic("implement me")
}

func (a *apiServer) DeleteRoom(w http.ResponseWriter, r *http.Request, locationID float32, roomID float32) {
	//TODO implement me
	panic("implement me")
}

func (a *apiServer) GetRoom(w http.ResponseWriter, r *http.Request, locationID float32, roomID float32) {
	//TODO implement me
	panic("implement me")
}

func (a *apiServer) UpdateRoom(w http.ResponseWriter, r *http.Request, locationID float32, roomID float32) {
	//TODO implement me
	panic("implement me")
}

func (a *apiServer) GetLocations(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}

func (a *apiServer) CreateOrg(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}

func (a *apiServer) DeleteOrg(w http.ResponseWriter, r *http.Request, org int) {
	//TODO implement me
	panic("implement me")
}

func (a *apiServer) GetOrg(w http.ResponseWriter, r *http.Request, org int) {
	//TODO implement me
	panic("implement me")
}

func (a *apiServer) UpdateOrg(w http.ResponseWriter, r *http.Request, org int) {
	//TODO implement me
	panic("implement me")
}

func (a *apiServer) GetOrgEvents(w http.ResponseWriter, r *http.Request, org int) {
	//TODO implement me
	panic("implement me")
}

func (a *apiServer) GetOrgMembers(w http.ResponseWriter, r *http.Request, org int) {
	//TODO implement me
	panic("implement me")
}

func encode(w http.ResponseWriter, response any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(response)
}

func internalServerError(w http.ResponseWriter, err error) {
	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
}
