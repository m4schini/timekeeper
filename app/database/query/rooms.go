package query

import . "timekeeper/app/database/model"

func (q *Queries) GetRooms(offset, limit int) (rs []RoomModel, total int, err error) {
	row := q.DB.QueryRow(`SELECT COUNT(id) FROM timekeeper.rooms`)
	if err = row.Err(); err != nil {
		return nil, -1, err
	}
	err = row.Scan(&total)
	if err != nil {
		return nil, -1, err
	}
	if total == 0 {
		return []RoomModel{}, total, nil
	}

	rows, err := q.DB.Query(`
SELECT r.id as id,
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
LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, -1, err
	}

	rs = make([]RoomModel, 0, limit)
	for rows.Next() {
		var r RoomModel
		var l LocationModel
		err = rows.Scan(&r.ID, &r.Name, &r.LocationX, &r.LocationY, &r.LocationW, &r.LocationH,
			&l.ID, &l.Name, &l.File)
		if err != nil {
			return nil, 0, err
		}
		r.Location = l

		rs = append(rs, r)
	}
	return rs, -1, nil
}
