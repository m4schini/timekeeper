package query

import (
	"raumzeitalpaka/app/database/model"
)

type GetUser Handler[GetUserRequest, model.UserModel]

type GetUserRequest struct {
	ID int
}

type GetUserRequestHandler struct {
	DB Database
}

func (q *GetUserRequestHandler) Query(request GetUserRequest) (u model.UserModel, err error) {
	id := request.ID
	row := q.DB.QueryRow(`SELECT id, login_name, password FROM raumzeitalpaka.users WHERE id = $1`, id)
	if err = row.Err(); err != nil {
		return model.UserModel{}, err
	}

	err = row.Scan(&u.ID, &u.LoginName, &u.PasswordHash)
	return u, err
}
