package query

import (
	"raumzeitalpaka/app/database/model"
)

type GetLocations Handler[GetLocationsRequest, []model.LocationModel]

type GetLocationsRequest struct {
	Offset int
	Limit  int
}

type GetLocationsHandler struct {
	DB Database
}

func (q *GetLocationsHandler) Query(request GetLocationsRequest) (ls []model.LocationModel, err error) {
	limit := request.Limit
	offset := request.Offset
	rows, err := q.DB.Query(`
SELECT id, name, file 
FROM raumzeitalpaka.locations
ORDER BY name
LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, err
	}

	ls = make([]model.LocationModel, 0, limit)
	for rows.Next() {
		var l model.LocationModel
		err = rows.Scan(&l.ID, &l.Name, &l.File)
		if err != nil {
			return nil, err
		}

		ls = append(ls, l)
	}
	return ls, nil
}
