package query

import (
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
	row := q.DB.QueryRow(`SELECT id, login_name, password FROM raumzeitalpaka.users WHERE login_name = $1`, loginName)
	if err = row.Err(); err != nil {
		return model.UserModel{}, err
	}

	err = row.Scan(&u.ID, &u.LoginName, &u.PasswordHash)
	return u, err
}
