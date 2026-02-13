package query

import (
	"raumzeitalpaka/app/database/model"
)

type UserHasRole Handler[UserHasRoleRequest, UserHasRoleResponse]

type UserHasRoleRequest struct {
	UserId int
	Role   model.Role
}

type UserHasRoleResponse struct {
	HasRole bool
	Role    model.Role
}

type UserHasRoleHandler struct {
	DB Database
}

func (q *UserHasRoleHandler) Query(request UserHasRoleRequest) (UserHasRoleResponse, error) {

	row := q.DB.QueryRow(`SELECT role FROM raumzeitalpaka.users WHERE id = $1`,
		request.UserId)
	if err := row.Err(); err != nil {
		return UserHasRoleResponse{}, err
	}

	var role model.Role
	err := row.Scan(&role)
	if err != nil {
		return UserHasRoleResponse{}, err
	}

	r := UserHasRoleResponse{
		HasRole: request.Role == role,
		Role:    role,
	}

	return r, err
}
