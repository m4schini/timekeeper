package command

import "raumzeitalpaka/app/database/model"

type UpsertUser InsertHandler[UpsertUserRequest, int]

type UpsertUserRequest struct {
	ID           int
	LoginName    string
	DisplayName  string
	PasswordHash string
	Role         model.Role
}

type UpsertUserHandler struct {
	DB Database
}

func (c *UpsertUserHandler) Execute(m UpsertUserRequest) (id int, err error) {
	row := c.DB.QueryRow(`
INSERT INTO raumzeitalpaka.users (id, login_name, display_name, password, role) 
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (id) DO UPDATE SET role = $5, display_name = $3
RETURNING id`, m.ID, m.LoginName, m.DisplayName, m.PasswordHash, m.Role)
	if err = row.Err(); err != nil {
		return m.ID, err
	}

	err = row.Scan(&id)
	return id, err
}
