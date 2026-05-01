package query

import (
	"context"
	"time"
)

type GetOrganisationMembers Handler[GetOrganisationMembersRequest, []OrganisationMember]

type GetOrganisationMembersRequest struct {
	OrganisationID int
}

type OrganisationMember struct {
	UserID      int       `json:"user_id"`
	LoginName   string    `json:"login_name"`
	DisplayName string    `json:"display_name"`
	LastLogin   time.Time `json:"last_login"`
	Role        string    `json:"role"`
}

type GetOrganisationMembersHandler struct {
	DB Database
}

func (q *GetOrganisationMembersHandler) Query(ctx context.Context, request GetOrganisationMembersRequest) (members []OrganisationMember, err error) {
	rows, err := q.DB.Query(`
SELECT id, login_name, display_name, last_login, ohu.role
FROM raumzeitalpaka.organisation_has_user as ohu
JOIN raumzeitalpaka.users u on u.id = ohu.user_id
WHERE organisation_id = $1`, request.OrganisationID)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var member OrganisationMember
		err = rows.Scan(&member.UserID, &member.LoginName, &member.DisplayName, &member.LastLogin, &member.Role)
		if err != nil {
			return nil, err
		}

		members = append(members, member)
	}

	return members, err
}
