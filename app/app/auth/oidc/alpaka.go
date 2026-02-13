package oidc

import "raumzeitalpaka/app/database/model"

var AlpakaRoleMapper = map[string]model.Role{
	"Menti":        model.RoleMentor,
	"Mentor*innen": model.RoleMentor,
	"Orga":         model.RoleOrganizer,
	"":             model.RoleParticipant,
}
