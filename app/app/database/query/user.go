package query

import (
	"database/sql"
	"raumzeitalpaka/app/database/model"
)

type GetUser Handler[GetUserRequest, model.UserModel]

type GetUserRequest struct {
	ID int
}

type GetUserHandler struct {
	DB Database
}

func (q *GetUserHandler) Query(request GetUserRequest) (u model.UserModel, err error) {
	id := request.ID
	row := q.DB.QueryRow(`SELECT id, login_name, display_name, password, last_login FROM raumzeitalpaka.users WHERE id = $1`, id)
	if err = row.Err(); err != nil {
		return model.UserModel{}, err
	}

	var ts sql.NullTime
	err = row.Scan(&u.ID, &u.LoginName, &u.DisplayName, &u.PasswordHash, &ts)
	if ts.Valid {
		u.LastLogin = ts.Time
	}
	return u, err
}
