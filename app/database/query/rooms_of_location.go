package query

import . "timekeeper/app/database/model"

func (q *Queries) GetRoomsOfLocation(location int, offset, limit int) (rs []RoomModel, total int, err error) {
	row := q.DB.QueryRow(`SELECT COUNT(id) FROM timekeeper.rooms WHERE location = $1`, location)
	if err = row.Err(); err != nil {
		return nil, -1, err
	}
	err = row.Scan(&total)
	if err != nil {
		return nil, -1, err
	}
	if total == 0 || limit == 0 {
		return []RoomModel{}, total, nil
	}

	rows, err := q.DB.Query(`
SELECT r.id as id,
       location,
       r.name as name,
       location_x,
       location_y,
       location_w,
       location_h,
       l.id as location_id,
       l.name as location_name,
       l.file as file
FROM timekeeper.rooms r
JOIN timekeeper.locations l
ON r.location = l.id
WHERE l.id = $1
LIMIT $2 OFFSET $3`,
		location, limit, offset)
	if err != nil {
		return nil, total, err
	}

	rs = make([]RoomModel, 0, limit)
	for rows.Next() {
		var r RoomModel
		var l LocationModel
		err = row.Scan(&r.ID, &r.Name, &r.LocationX, &r.LocationY, &r.LocationW, &r.LocationH,
			l.ID, l.Name, l.File)
		if err != nil {
			return nil, 0, err
		}
		r.Location = l

		rs = append(rs, r)
	}
	return rs, total, nil
}
