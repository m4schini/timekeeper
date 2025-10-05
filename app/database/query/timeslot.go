package query

import (
	"time"
	. "timekeeper/app/database/model"
)

func (q *Queries) GetTimeslot(id int) (t TimeslotModel, err error) {
	row := q.DB.QueryRow(`
SELECT ts.id as id,
       title,
       note,
       day,
       ts.start as start,
       role,

       e.id as event_id,
       e.name as event_name,
       e.start as event_start,

       r.id as room_id,
       r.name as room_name,
       r.location_x as room_x,
       r.location_y as room_y,
       r.location_w as room_w,
       r.location_h as room_h,

       l.id as location_id,
       l.name as location_name,
       l.file as location_file
FROM timekeeper.timeslots ts
JOIN timekeeper.rooms r ON r.id = ts.room
JOIN timekeeper.events e on e.id = ts.event
JOIN timekeeper.locations l on l.id = r.location
WHERE ts.id = $1 ORDER BY ts.start `, id)
	if err = row.Err(); err != nil {
		return TimeslotModel{}, err
	}

	var e EventModel
	var r RoomModel
	var l LocationModel
	err = row.Scan(&t.ID, &t.Title, &t.Note, &t.Day, &t.Start, &t.Role, &e.ID, &e.Name, &e.Start, &r.ID, &r.Name, &r.LocationX, &r.LocationY, &r.LocationW, &r.LocationH, &l.ID, &l.Name, &l.File)
	if err != nil {
		return t, err
	}
	now := time.Now()
	e.Start = e.Start.In(now.Location())
	t.Event = e
	r.Location = l
	t.Room = r

	return t, nil
}
