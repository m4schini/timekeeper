package query

import (
	"context"
	"raumzeitalpaka/app/database/model"
)

type GetOrganisation Handler[GetOrganisationRequest, model.OrganisationModel]

type GetOrganisationRequest struct {
	ID int
}

type GetOrganisationHandler struct {
	DB Database
}

func (q *GetOrganisationHandler) Query(ctx context.Context, request GetOrganisationRequest) (u model.OrganisationModel, err error) {
	row := q.DB.QueryRow(`SELECT id, slug, name FROM raumzeitalpaka.organisations WHERE id = $1`, request.ID)
	if err = row.Err(); err != nil {
		return model.OrganisationModel{}, err
	}

	err = row.Scan(&u.ID, &u.Slug, &u.Name)
	return u, err
}
