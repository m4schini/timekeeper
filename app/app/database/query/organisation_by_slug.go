package query

import (
	"context"
	"raumzeitalpaka/app/database/model"
)

type GetOrganisationBySlug Handler[GetOrganisationBySlugRequest, model.OrganisationModel]

type GetOrganisationBySlugRequest struct {
	Slug string
}

type GetOrganisationBySlugHandler struct {
	DB Database
}

func (q *GetOrganisationBySlugHandler) Query(ctx context.Context, request GetOrganisationBySlugRequest) (u model.OrganisationModel, err error) {
	slug := request.Slug
	row := q.DB.QueryRow(`SELECT id, slug, name FROM raumzeitalpaka.organisations WHERE slug = $1`, slug)
	if err = row.Err(); err != nil {
		return model.OrganisationModel{}, err
	}

	err = row.Scan(&u.ID, &u.Slug, &u.Name)
	return u, err
}
