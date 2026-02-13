package command

import (
	"raumzeitalpaka/app/database/model"
)

type GroupAddUser Handler[GroupAddUserRequest]

type GroupAddUserRequest struct {
	UserId  int
	GroupId int
	Role    model.Role
}

type GroupAddUserHandler struct {
	DB Database
}

func (c *GroupAddUserHandler) Execute(m GroupAddUserRequest) error {
	_, err := c.DB.Exec(`
INSERT INTO raumzeitalpaka.group_has_user (user_id, group_id, role) 
VALUES ($1, $2, $3)
ON CONFLICT (user_id, group_id) 
DO UPDATE SET role = EXCLUDED.role;
`, m.UserId, m.GroupId, m.Role)
	return err
}
