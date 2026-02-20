package command

import "context"

type CreateUser InsertHandler[CreateUserRequest, int]

type CreateUserRequest struct {
	LoginName    string
	PasswordHash string
}

type CreateUserHandler struct {
	DB Database
}

func (c *CreateUserHandler) Execute(ctx context.Context, m CreateUserRequest) (id int, err error) {
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
