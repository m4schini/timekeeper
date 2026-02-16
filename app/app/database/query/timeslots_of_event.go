package query

import (
	"database/sql"
	"raumzeitalpaka/app/database/model"
	"raumzeitalpaka/config"
	"time"

	"github.com/lib/pq"
	"go.uber.org/zap"
)

type GetTimeslotsOfEvent Handler[GetTimeslotsOfEventRequest, GetTimeslotsOfEventResponse]

type GetTimeslotsOfEventRequest struct {
	EventId int
	Roles   []model.Role
	Offset  int
	Limit   int
}

type GetTimeslotsOfEventResponse struct {
	Timeslots []model.TimeslotModel
	Total     int
}

type GetTimeslotsOfEventHandler struct {
	DB Database
}

func (q *GetTimeslotsOfEventHandler) Query(request GetTimeslotsOfEventRequest) (GetTimeslotsOfEventResponse, error) {
	var (
		total  int
		event  = request.EventId
		roles  = request.Roles
		limit  = request.Limit
		offset = request.Offset
	)
	row := q.DB.QueryRow(`SELECT COUNT(id) FROM raumzeitalpaka.timeslots WHERE event = $1`, event)
	if err := row.Err(); err != nil {
		return GetTimeslotsOfEventResponse{}, err
	}
	err := row.Scan(&total)
	if err != nil {
		return GetTimeslotsOfEventResponse{}, err
	}
	if total == 0 || limit == 0 {
		return GetTimeslotsOfEventResponse{}, err
	}

	rows, err := q.DB.Query(`
SELECT ts.id as id,
       ts.parent_id as parent_id,
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
WHERE e.id = $1 AND ts.role = ANY($4) ORDER BY ts.start, ts.parent_id NULLS FIRST, ts.rank, ts.note LIMIT $2 OFFSET $3 `,
		event, limit, offset, pq.Array(roles))
	if err != nil {
		return GetTimeslotsOfEventResponse{}, err
	}

	_ts := make([]*model.TimeslotModel, 0, limit)
	_tsMap := make(map[int64]*model.TimeslotModel)
	for rows.Next() {
		var parentId sql.NullInt64
		var e model.EventModel
		var r model.RoomModel
		var l model.LocationModel
		var t model.TimeslotModel
		var durationInSeconds int
		err = rows.Scan(
			&t.ID, &parentId, &t.GUID, &t.Title, &t.Note, &t.Day, &t.Start, &t.Role, &durationInSeconds,
			&e.ID, &e.GUID, &e.Name, &e.Start,
			&r.ID, &r.GUID, &r.Name, &r.LocationX, &r.LocationY, &r.LocationW, &r.LocationH, &r.Description,
			&l.ID, &l.GUID, &l.Name, &l.File)
		if err != nil {
			return GetTimeslotsOfEventResponse{}, err
		}
		t.Duration = time.Duration(durationInSeconds) * time.Second
		e.Start = e.Start.In(config.Timezone())
		t.Event = e
		r.Location = l
		t.Room = r

		if !parentId.Valid {
			_ts = append(_ts, &t)
			_tsMap[int64(t.ID)] = &t
		} else {
			parent, ok := _tsMap[parentId.Int64]
			if ok {
				chldrn := parent.Children
				if chldrn == nil {
					chldrn = make([]model.TimeslotModel, 0)
				}
				chldrn = append(chldrn, t)
				parent.Children = chldrn
			} else {
				zap.L().Warn("parent timeslot doesn't exist yet", zap.Int64("parent_id", parentId.Int64), zap.Int("child_id", t.ID))
			}
		}
	}
	ts := make([]model.TimeslotModel, len(_ts))
	for i, t := range _ts {
		ts[i] = *t
	}

	return GetTimeslotsOfEventResponse{
		Timeslots: ts,
		Total:     total,
	}, nil
}
