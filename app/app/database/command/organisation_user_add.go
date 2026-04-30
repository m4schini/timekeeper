package command

import (
	"context"
	"raumzeitalpaka/app/database/model"
)

type OrganisationAddUser Handler[OrganisationAddUserRequest]

type OrganisationAddUserRequest struct {
	UserId         int
	OrganisationId int
	Role           model.Role
}

type OrganisationAddUserHandler struct {
	DB Database
}

func (c *OrganisationAddUserHandler) Execute(ctx context.Context, m OrganisationAddUserRequest) error {
	_, err := c.DB.Exec(`
INSERT INTO raumzeitalpaka.organisation_has_user (user_id, organisation_id, role) 
VALUES ($1, $2, $3)
ON CONFLICT (user_id, organisation_id) 
DO UPDATE SET role = EXCLUDED.role;
`, m.UserId, m.OrganisationId, m.Role)
	return err
}
