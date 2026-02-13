package query

import (
	"raumzeitalpaka/app/database/model"
)

type GetRooms Handler[GetRoomsRequest, GetRoomsResponse]

type GetRoomsRequest struct {
	Offset int
	Limit  int
}

type GetRoomsResponse struct {
	Rooms []model.RoomModel
	Total int
}

type GetRoomsHandler struct {
	DB Database
}

func (q *GetRoomsHandler) Query(request GetRoomsRequest) (res GetRoomsResponse, err error) {
	offset := request.Offset
	limit := request.Limit
	var total int
	row := q.DB.QueryRow(`SELECT COUNT(id) FROM raumzeitalpaka.rooms`)
	if err = row.Err(); err != nil {
		return GetRoomsResponse{}, err
	}
	err = row.Scan(&total)
	if err != nil {
		return GetRoomsResponse{}, err
	}
	if total == 0 {
		return GetRoomsResponse{}, nil
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
ORDER BY location_name, name
LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return GetRoomsResponse{}, err
	}

	rs := make([]model.RoomModel, 0, limit)
	for rows.Next() {
		var r model.RoomModel
		var l model.LocationModel
		err = rows.Scan(&r.ID, &r.Name, &r.LocationX, &r.LocationY, &r.LocationW, &r.LocationH, &r.Description,
			&l.ID, &l.Name, &l.File)
		if err != nil {
			return GetRoomsResponse{}, err
		}
		r.Location = l

		rs = append(rs, r)
	}

	return GetRoomsResponse{
		Rooms: rs,
		Total: total,
	}, nil
}
