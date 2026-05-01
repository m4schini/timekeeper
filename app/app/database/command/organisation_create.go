package command

import (
	"context"
	"strings"
)

type CreateOrganisation InsertHandler[CreateOrganisationRequest, int]

type CreateOrganisationRequest struct {
	Name string
	Slug string
}

type CreateOrganisationHandler struct {
	DB Database
}

func (c *CreateOrganisationHandler) Execute(ctx context.Context, m CreateOrganisationRequest) (id int, err error) {
	row := c.DB.QueryRow(`
INSERT INTO raumzeitalpaka.organisations (slug, name) 
VALUES ($1, $2)
RETURNING id`, m.slug(), m.Name)
	err = row.Err()
	if err != nil {
		return -1, err
	}

	err = row.Scan(&id)
	if err != nil {
		return -1, err
	}

	_, err = c.DB.Exec(`
INSERT INTO raumzeitalpaka.organisation_roles (organisation, id, name, required) 
VALUES ($1, $2, $3, true)
RETURNING id`, id, "participant", "Teilnehmer*in")
	if err != nil {
		return -1, err
	}

	_, err = c.DB.Exec(`
INSERT INTO raumzeitalpaka.organisation_roles (organisation, id, name, required) 
VALUES ($1, $2, $3, true)
RETURNING id`, id, "mentor", "Mentor*in")
	if err != nil {
		return -1, err
	}

	_, err = c.DB.Exec(`
INSERT INTO raumzeitalpaka.organisation_roles (organisation, id, name, required) 
VALUES ($1, $2, $3, true)
RETURNING id`, id, "organizer", "Orga")
	if err != nil {
		return -1, err
	}

	return id, err
}

func (g *CreateOrganisationRequest) slug() string {
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
