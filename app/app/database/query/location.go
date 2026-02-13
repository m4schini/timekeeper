package query

import (
	"raumzeitalpaka/app/database/model"
)

type GetLocation Handler[GetLocationRequest, model.LocationModel]

type GetLocationRequest struct {
	LocationId int
}

type GetLocationHandler struct {
	DB Database
}

func (q *GetLocationHandler) Query(request GetLocationRequest) (l model.LocationModel, err error) {
	id := request.LocationId
	row := q.DB.QueryRow(`SELECT id, name, file, osm_id FROM raumzeitalpaka.locations WHERE id = $1`, id)
	if err = row.Err(); err != nil {
		return model.LocationModel{}, err
	}

	err = row.Scan(&l.ID, &l.Name, &l.File, &l.OsmId)
	return l, err
}
