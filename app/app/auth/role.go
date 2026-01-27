package auth

import (
	"raumzeitalpaka/app/database/model"
	"strings"
)

type Roler interface {
	Role(groups []string) (model.Role, error)
}

type AlpakaRoleRules struct {
}

func (a AlpakaRoleRules) Role(groups []string) (model.Role, error) {
	var role = model.RoleParticipant

	for _, group := range groups {
		if strings.HasSuffix(strings.ToLower(group), ":orga") {
			return model.RoleOrganizer, nil
		}
		//if strings.HasSuffix(group, ":Menti") {
		//	role = model.RoleMentor
		//}
	}

	return role, nil
}
