package query

import (
	"time"
	. "timekeeper/app/database/model"
	"timekeeper/config"
)

func (q *Queries) GetTimeslotsOfEvent(event int, offset, limit int) (ts []TimeslotModel, total int, err error) {
	row := q.DB.QueryRow(`SELECT COUNT(id) FROM timekeeper.timeslots WHERE event = $1`, event)
	if err = row.Err(); err != nil {
		return nil, -1, err
	}
	err = row.Scan(&total)
	if err != nil {
		return nil, -1, err
	}
	if total == 0 || limit == 0 {
		return []TimeslotModel{}, total, err
	}

	rows, err := q.DB.Query(`
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


FROM timekeeper.timeslots ts
JOIN timekeeper.rooms r ON r.id = ts.room
JOIN timekeeper.events e on e.id = ts.event
JOIN timekeeper.locations l on l.id = r.location
WHERE e.id = $1 ORDER BY ts.start, ts.note LIMIT $2 OFFSET $3 `,
		event, limit, offset)
	if err != nil {
		return nil, total, err
	}

	ts = make([]TimeslotModel, 0, limit)
	for rows.Next() {
		var e EventModel
		var r RoomModel
		var l LocationModel
		var t TimeslotModel
		var durationInSeconds int
		err = rows.Scan(
			&t.ID, &t.GUID, &t.Title, &t.Note, &t.Day, &t.Start, &t.Role, &durationInSeconds,
			&e.ID, &e.GUID, &e.Name, &e.Start,
			&r.ID, &r.GUID, &r.Name, &r.LocationX, &r.LocationY, &r.LocationW, &r.LocationH, &r.Description,
			&l.ID, &l.GUID, &l.Name, &l.File)
		if err != nil {
			return nil, 0, err
		}
		t.Duration = time.Duration(durationInSeconds) * time.Second
		e.Start = e.Start.In(config.Timezone())
		t.Event = e
		r.Location = l
		t.Room = r

		ts = append(ts, t)
	}
	return ts, total, nil
}
