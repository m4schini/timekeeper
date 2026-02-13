package query

import (
	"raumzeitalpaka/app/database/model"
)

type UserHasEventRole Handler[UserHasEventRoleRequest, UserHasEventRoleResponse]

type UserHasEventRoleRequest struct {
	UserId  int
	EventId int
	Role    model.Role
}

type UserHasEventRoleResponse struct {
	HasEvent bool
	HasRole  bool
	Role     model.Role
	Group    int
}

type UserHasEventRoleHandler struct {
	DB Database
}

func (q *UserHasEventRoleHandler) Query(request UserHasEventRoleRequest) (UserHasEventRoleResponse, error) {
	//TODO
	row := q.DB.QueryRow(`SELECT role FROM raumzeitalpaka.group_has_user WHERE user_id = $1 AND group_id = $2`,
		request.UserId, request.EventId)
	if err := row.Err(); err != nil {
		return UserHasEventRoleResponse{}, err
	}

	var role model.Role
	err := row.Scan(&role)
	if err != nil {
		return UserHasEventRoleResponse{}, err
	}

	r := UserHasEventRoleResponse{
		HasEvent: true,
		HasRole:  request.Role == role,
		Role:     role,
	}

	return r, err
}
