package query

import . "raumzeitalpaka/app/database/model"

func (q *Queries) GetUserByLoginName(loginName string) (u UserModel, err error) {
	row := q.DB.QueryRow(`SELECT id, login_name, password FROM raumzeitalpaka.users WHERE login_name = $1`, loginName)
	if err = row.Err(); err != nil {
		return UserModel{}, err
	}

	err = row.Scan(&u.ID, &u.LoginName, &u.PasswordHash)
	return u, err
}
