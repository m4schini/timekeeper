package query

import (
	"raumzeitalpaka/app/database/model"
)

type UserHasGroupRole Handler[UserHasGroupRoleRequest, UserHasGroupRoleResponse]

type UserHasGroupRoleRequest struct {
	UserId  int
	GroupId int
	Role    model.Role
}

type UserHasGroupRoleResponse struct {
	HasGroup bool
	HasRole  bool
	Role     model.Role
}

type UserHasGroupRoleHandler struct {
	DB Database
}

func (q *UserHasGroupRoleHandler) Query(request UserHasGroupRoleRequest) (UserHasGroupRoleResponse, error) {

	row := q.DB.QueryRow(`SELECT role FROM raumzeitalpaka.group_has_user WHERE user_id = $1 AND group_id = $2`,
		request.UserId, request.GroupId)
	if err := row.Err(); err != nil {
		return UserHasGroupRoleResponse{}, err
	}

	var role model.Role
	err := row.Scan(&role)
	if err != nil {
		return UserHasGroupRoleResponse{}, err
	}

	r := UserHasGroupRoleResponse{
		HasGroup: true,
		HasRole:  request.Role == role,
		Role:     role,
	}

	return r, err
}
