package model

type Role string

const (
	RoleOrganizer   Role = "Organizer"
	RoleMentor      Role = "Mentor"
	RoleParticipant Role = "Participant"
)

var RoleHierarchy = []Role{RoleOrganizer, RoleMentor, RoleParticipant}

func (r Role) Title() string {
	switch r {
	case RoleOrganizer:
		return "Orga"
	case RoleMentor:
		return "Mentor*innen"
	case RoleParticipant:
		return "Teilnehmer*innen"
	default:
		panic("unknown role")
	}
}

func (r Role) Color() string {
	switch r {
	case RoleOrganizer:
		return "#ffd003"
	case RoleMentor:
		return "#ea680c"
	case RoleParticipant:
		return "#4cad37"
	default:
		panic("unknown role")
	}
}

func RoleFrom(role string) Role {
	switch Role(role) {
	case RoleOrganizer:
		return RoleOrganizer
	case RoleMentor:
		return RoleMentor
	case RoleParticipant:
		return RoleParticipant
	default:
		return RoleParticipant
	}
}

func ValidRole(role Role) bool {
	return role == RoleOrganizer || role == RoleMentor || role == RoleParticipant
}

func HigherRole(r1, r2 Role) Role {
	if r1 == "" {
		return r2
	}
	if r2 == "" {
		return r1
	}
	for _, role := range RoleHierarchy {
		if r1 == role {
			return r1
		}
		if r2 == role {
			return r2
		}
	}

	return r1
}
