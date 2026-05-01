package query

import (
	"context"
	"raumzeitalpaka/app/database/model"
)

type GetUserOrganisations Handler[GetUserOrganisationsRequest, []model.UserOrganisationMembership]

type GetUserOrganisationsRequest struct {
	UserID int
}

type GetUserOrganisationsHandler struct {
	DB Database
}

func (q *GetUserOrganisationsHandler) Query(ctx context.Context, request GetUserOrganisationsRequest) (ms []model.UserOrganisationMembership, err error) {
	id := request.UserID
	rows, err := q.DB.Query(`
SELECT organisation_id, role, slug, name
FROM raumzeitalpaka.organisation_has_user 
JOIN raumzeitalpaka.organisations o on o.id = organisation_has_user.organisation_id
WHERE user_id = $1`, id)
	if err != nil {
		return nil, err
	}

	ms = make([]model.UserOrganisationMembership, 0)
	for rows.Next() {
		var m model.UserOrganisationMembership
		err = rows.Scan(&m.OrganisationID, &m.Role, &m.Slug, &m.Name)
		if err != nil {
			return nil, err
		}
		ms = append(ms, m)
	}

	return ms, nil
}
