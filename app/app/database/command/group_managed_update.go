package command

import (
	"raumzeitalpaka/app/database/model"

	"go.uber.org/zap"
)

type UpdateManagedGroupsAssignments Handler[UpdateManagedGroupsAssignmentsRequest]

type UpdateManagedGroupsAssignmentsRequest struct {
	UserId      int
	Assignments []GroupAssignment
}

type GroupAssignment struct {
	GroupId int
	Role    model.Role
}

type UpdateManagedGroupsAssignmentsHandler struct {
	DB Database
}

func (c *UpdateManagedGroupsAssignmentsHandler) Execute(m UpdateManagedGroupsAssignmentsRequest) error {
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
DELETE FROM raumzeitalpaka.group_has_user WHERE user_id = $1 AND managed = true;
`, m.UserId)
	if err != nil {
		rollback()
		return err
	}

	for _, assignment := range m.Assignments {
		log.Debug("adding managed group assignment", zap.Int("user", m.UserId), zap.Any("role", assignment.Role), zap.Int("group", assignment.GroupId))
		_, err = tx.Exec(`
INSERT INTO raumzeitalpaka.group_has_user(user_id, group_id, role, managed) 
VALUES ($1, $2, $3, true)`, m.UserId, assignment.GroupId, assignment.Role)
		if err != nil {
			rollback()
			return err
		}
	}

	return tx.Commit()
}
