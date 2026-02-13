package query

import (
	"raumzeitalpaka/app/database/model"
)

type GetRoom Handler[GetRoomRequest, model.RoomModel]

type GetRoomRequest struct {
	RoomId int
}

type GetRoomHandler struct {
	DB Database
}

func (q *GetRoomHandler) Query(request GetRoomRequest) (r model.RoomModel, err error) {
	id := request.RoomId
	row := q.DB.QueryRow(`
SELECT r.id as id,
       location,
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
WHERE r.id = $1`, id)
	if err = row.Err(); err != nil {
		return model.RoomModel{}, err
	}

	var l model.LocationModel
	err = row.Scan(&r.ID, &r.Name, &r.LocationX, &r.LocationY, &r.LocationW, &r.LocationH, &r.Description,
		l.ID, l.Name, l.File)
	r.Location = l
	return r, err
}
