package query

import (
	"context"
	"raumzeitalpaka/app/database/model"
)

type GetGroupBySlug Handler[GetGroupBySlugRequest, model.GroupModel]

type GetGroupBySlugRequest struct {
	Slug string
}

type GetGroupBySlugHandler struct {
	DB Database
}

func (q *GetGroupBySlugHandler) Query(ctx context.Context, request GetGroupBySlugRequest) (u model.GroupModel, err error) {
	slug := request.Slug
	row := q.DB.QueryRow(`SELECT id, slug, name FROM raumzeitalpaka.groups WHERE slug = $1`, slug)
	if err = row.Err(); err != nil {
		return model.GroupModel{}, err
	}

	err = row.Scan(&u.ID, &u.Slug, &u.Name)
	return u, err
}
