package model

type Role string

const (
	RoleOrganizer   Role = "Organizer"
	RoleMentor      Role = "Mentor"
	RoleParticipant Role = "Participant"
)

func RoleFrom(role string) Role {
	switch Role(role) {
	case RoleMentor:
		return RoleMentor
	case RoleParticipant:
		return RoleParticipant
	default:
		return RoleOrganizer
	}
}
