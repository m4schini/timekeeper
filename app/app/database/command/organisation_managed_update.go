package command

import (
	"context"
	"raumzeitalpaka/app/database/model"

	"go.uber.org/zap"
)

type UpdateManagedOrganisationAssignments Handler[UpdateManagedOrganisationAssignmentsRequest]

type UpdateManagedOrganisationAssignmentsRequest struct {
	UserId      int
	Assignments []OrganisationAssignment
}

type OrganisationAssignment struct {
	OrganisationId int
	Role           model.Role
}

type UpdateManagedOrganisationAssignmentsHandler struct {
	DB Database
}

func (c *UpdateManagedOrganisationAssignmentsHandler) Execute(ctx context.Context, m UpdateManagedOrganisationAssignmentsRequest) error {
	log := zap.L().Named("command")
	tx, err := c.DB.Begin()
	if err != nil {
		return err
	}
	rollback := func() {
		if err := tx.Rollback(); err != nil {
			log.WithOptions(zap.AddCallerSkip(1)).Error("failed to rollback")
		}
	}

	log.Debug("clearing existing managed group assignments", zap.Int("user", m.UserId))
	_, err = tx.Exec(`
DELETE FROM raumzeitalpaka.organisation_has_user WHERE user_id = $1 AND managed = true;
`, m.UserId)
	if err != nil {
		rollback()
		return err
	}

	for _, assignment := range m.Assignments {
		log.Debug("adding managed group assignment", zap.Int("user", m.UserId), zap.Any("role", assignment.Role), zap.Int("group", assignment.OrganisationId))
		_, err = tx.Exec(`
INSERT INTO raumzeitalpaka.organisation_has_user(user_id, organisation_id, role, managed) 
VALUES ($1, $2, $3, true)`, m.UserId, assignment.OrganisationId, assignment.Role)
		if err != nil {
			rollback()
			return err
		}
	}

	return tx.Commit()
}
