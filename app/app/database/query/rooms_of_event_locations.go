package query

import . "raumzeitalpaka/app/database/model"

func (q *Queries) GetRoomsOfEventLocations(event int) (rs []RoomModel, err error) {
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
		event)
	if err != nil {
		return nil, err
	}

	rs = make([]RoomModel, 0)
	for rows.Next() {
		var r RoomModel
		var l LocationModel
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
