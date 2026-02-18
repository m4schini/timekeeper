package command

import (
	"time"
)

type UpdateLastLogin Handler[UpdateLastLoginRequest]

type UpdateLastLoginRequest struct {
	ID        int
	Timestamp time.Time
}

type UpdateLastLoginHandler struct {
	DB Database
}

func (c *UpdateLastLoginHandler) Execute(m UpdateLastLoginRequest) (err error) {
	row := c.DB.QueryRow(`
UPDATE raumzeitalpaka.users
SET
    last_login = $1
WHERE id = $2
`, m.Timestamp, m.ID)
	if err = row.Err(); err != nil {
		return err
	}

	return nil
}
