package query

import (
	"raumzeitalpaka/app/database/model"
	"raumzeitalpaka/config"
	"time"
)

type GetTimeslot Handler[GetTimeslotRequest, model.TimeslotModel]

type GetTimeslotRequest struct {
	TimeslotId int
}

type GetTimeslotHandler struct {
	DB Database
}

func (q *GetTimeslotHandler) Query(request GetTimeslotRequest) (t model.TimeslotModel, err error) {
	id := request.TimeslotId
	row := q.DB.QueryRow(`
SELECT ts.id as id,
       ts.guid as ts_guid,
       title,
       note,
       day,
       ts.start as start,
       role,
       FLOOR(EXTRACT(EPOCH FROM duration)) AS total_seconds,

       e.id as event_id,
       e.guid as event_guid,
       e.name as event_name,
       e.start as event_start,

       r.id as room_id,
       r.guid as room_guid,
       r.name as room_name,
       r.location_x as room_x,
       r.location_y as room_y,
       r.location_w as room_w,
       r.location_h as room_h,
       r.description as room_description,

       l.id as location_id,
       l.guid as location_guid,
       l.name as location_name,
       l.file as location_file
FROM raumzeitalpaka.timeslots ts
JOIN raumzeitalpaka.rooms r ON r.id = ts.room
JOIN raumzeitalpaka.events e on e.id = ts.event
JOIN raumzeitalpaka.locations l on l.id = r.location
WHERE ts.id = $1 ORDER BY ts.start `, id)
	if err = row.Err(); err != nil {
		return model.TimeslotModel{}, err
	}

	var e model.EventModel
	var r model.RoomModel
	var l model.LocationModel
	var durationInSeconds int
	err = row.Scan(
		&t.ID, &t.GUID, &t.Title, &t.Note, &t.Day, &t.Start, &t.Role, &durationInSeconds,
		&e.ID, &e.GUID, &e.Name, &e.Start,
		&r.ID, &r.GUID, &r.Name, &r.LocationX, &r.LocationY, &r.LocationW, &r.LocationH, &r.Description,
		&l.ID, &l.GUID, &l.Name, &l.File)
	if err != nil {
		return t, err
	}
	t.Duration = time.Duration(durationInSeconds) * time.Second
	e.Start = e.Start.In(config.Timezone())
	t.Event = e
	r.Location = l
	t.Room = r

	return t, nil
}
