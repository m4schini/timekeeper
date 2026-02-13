package query

import (
	"raumzeitalpaka/app/database/model"
)

type GetRoomsOfEventLocations Handler[GetRoomsOfEventLocationsRequest, []model.RoomModel]

type GetRoomsOfEventLocationsRequest struct {
	EventId int
}

type GetRoomsOfEventLocationsHandler struct {
	DB Database
}

func (q *GetRoomsOfEventLocationsHandler) Query(request GetRoomsOfEventLocationsRequest) (rs []model.RoomModel, err error) {
	id := request.EventId
	rows, err := q.DB.Query(`SELECT r.id as id,
       r.name as name,
       location_x,
       location_y,
       location_w,
       location_h,
       description,
       l.id as location_id,
       l.name as location_name,
       l.file as file
FROM raumzeitalpaka.rooms r
         JOIN raumzeitalpaka.locations l
              ON r.location = l.id
WHERE l.id 
          IN (SELECT location FROM raumzeitalpaka.event_has_location WHERE event = $1)
          ORDER BY l.name, r.name`,
		id)
	if err != nil {
		return nil, err
	}

	rs = make([]model.RoomModel, 0)
	for rows.Next() {
		var r model.RoomModel
		var l model.LocationModel
		err = rows.Scan(&r.ID, &r.Name, &r.LocationX, &r.LocationY, &r.LocationW, &r.LocationH, &r.Description,
			&l.ID, &l.Name, &l.File)
		if err != nil {
			return nil, err
		}
		r.Location = l

		rs = append(rs, r)
	}
	return rs, nil
}
