package model

type Role string

const (
	RoleOrganizer   Role = "Organizer"
	RoleMentor      Role = "Mentor"
	RoleParticipant Role = "Participant"
)

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
