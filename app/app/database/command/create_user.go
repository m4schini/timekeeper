package command

import (
	. "raumzeitalpaka/app/database/model"
)

func (c *Commands) CreateUser(m CreateUserModel) (id int, err error) {
	row := c.DB.QueryRow(`
INSERT INTO raumzeitalpaka.users (login_name, password) 
VALUES ($1, $2)
RETURNING id`, m.LoginName, m.PasswordHash)
	if err = row.Err(); err != nil {
		return -1, err
	}

	err = row.Scan(&id)
	return id, err
}
