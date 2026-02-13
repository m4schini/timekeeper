package query

import (
	"raumzeitalpaka/app/database/model"
)

type GetRoomsOfLocation Handler[GetRoomsOfLocationRequest, GetRoomsOfLocationResponse]

type GetRoomsOfLocationRequest struct {
	LocationId int
	Offset     int
	Limit      int
}

type GetRoomsOfLocationResponse struct {
	Rooms []model.RoomModel
	Total int
}

type GetRoomsOfLocationHandler struct {
	DB Database
}

func (q *GetRoomsOfLocationHandler) Query(request GetRoomsOfLocationRequest) (res GetRoomsOfLocationResponse, err error) {
	var total int
	row := q.DB.QueryRow(`SELECT COUNT(id) FROM raumzeitalpaka.rooms WHERE location = $1`, request.LocationId)
	if err = row.Err(); err != nil {
		return GetRoomsOfLocationResponse{}, err
	}
	err = row.Scan(&total)
	if err != nil {
		return GetRoomsOfLocationResponse{}, err
	}
	if total == 0 || request.Limit == 0 {
		return GetRoomsOfLocationResponse{}, nil
	}

	rows, err := q.DB.Query(`
SELECT r.id as id,
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
WHERE r.location = $1
ORDER BY location_name, name
LIMIT $2 OFFSET $3`,
		request.LocationId, request.Limit, request.Offset)
	if err != nil {
		return GetRoomsOfLocationResponse{}, err
	}

	rs := make([]model.RoomModel, 0, request.Limit)
	for rows.Next() {
		var r model.RoomModel
		var l model.LocationModel
		err = rows.Scan(&r.ID, &r.Name, &r.LocationX, &r.LocationY, &r.LocationW, &r.LocationH, &r.Description,
			&l.ID, &l.Name, &l.File)
		if err != nil {
			return GetRoomsOfLocationResponse{}, err
		}
		r.Location = l

		rs = append(rs, r)
	}
	return GetRoomsOfLocationResponse{
		Rooms: rs,
		Total: total,
	}, nil
}
