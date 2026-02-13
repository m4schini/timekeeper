package command

import (
	"strings"
)

type CreateGroup InsertHandler[CreateGroupRequest, int]

type CreateGroupRequest struct {
	Name string
	Slug string
}

type CreateGroupHandler struct {
	DB Database
}

func (c *CreateGroupHandler) Execute(m CreateGroupRequest) (id int, err error) {
	row := c.DB.QueryRow(`
INSERT INTO raumzeitalpaka.groups (slug, name) 
VALUES ($1, $2)
RETURNING id`, m.slug(), m.Name)
	err = row.Err()
	if err != nil {
		return -1, err
	}

	err = row.Scan(&id)
	return id, err
}

func (g *CreateGroupRequest) slug() string {
	if g.Slug != "" {
		return g.Slug
	}

	slug := g.Name
	slug = strings.ToLower(slug)
	slug = strings.TrimSpace(slug)
	slug = strings.ReplaceAll(slug, "  ", " ")
	slug = strings.ReplaceAll(slug, " ", "-")

	return slug
}
