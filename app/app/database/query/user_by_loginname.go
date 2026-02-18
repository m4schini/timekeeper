package query

import (
	"database/sql"
	"raumzeitalpaka/app/database/model"
)

type GetUserByLoginName Handler[GetUserByLoginNameRequest, model.UserModel]

type GetUserByLoginNameRequest struct {
	LoginName string
}

type GetUserByLoginNameHandler struct {
	DB Database
}

func (q *GetUserByLoginNameHandler) Query(request GetUserByLoginNameRequest) (u model.UserModel, err error) {
	loginName := request.LoginName
	row := q.DB.QueryRow(`SELECT id, login_name, display_name, password, last_login FROM raumzeitalpaka.users WHERE login_name = $1`, loginName)
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
