package query

import (
	"context"
	"raumzeitalpaka/app/database/model"
)

type UserHasOrganisationRole Handler[UserHasOrganisationRoleRequest, UserHasOrganisationRoleResponse]

type UserHasOrganisationRoleRequest struct {
	UserId         int
	OrganisationId int
	Role           model.Role
}

type UserHasOrganisationRoleResponse struct {
	HasOrganisation bool
	HasRole         bool
	Role            model.Role
}

type UserHasOrganisationRoleHandler struct {
	DB Database
}

func (q *UserHasOrganisationRoleHandler) Query(ctx context.Context, request UserHasOrganisationRoleRequest) (UserHasOrganisationRoleResponse, error) {

	row := q.DB.QueryRow(`SELECT role FROM raumzeitalpaka.organisation_has_user WHERE user_id = $1 AND organisation_id = $2`,
		request.UserId, request.OrganisationId)
	if err := row.Err(); err != nil {
		return UserHasOrganisationRoleResponse{}, err
	}

	var role model.Role
	err := row.Scan(&role)
	if err != nil {
		return UserHasOrganisationRoleResponse{}, err
	}

	r := UserHasOrganisationRoleResponse{
		HasOrganisation: true,
		HasRole:         request.Role == role,
		Role:            role,
	}

	return r, err
}
