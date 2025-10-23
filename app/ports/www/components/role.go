package components

import (
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
	"timekeeper/app/database/model"
)

func RoleTag(role model.Role) Node {
	switch role {
	case model.RoleOrganizer:
		return Span(Class("role role-o"), Text("Orga"))
	case model.RoleMentor:
		return Span(Class("role role-m"), Text("Mentor*innen"))
	case model.RoleParticipant:
		return Span(Class("role role-t"), Text("Teilnehmer*innen"))
	default:
		return Span(Class("role role-o"), Text("Orga"))
	}
}
