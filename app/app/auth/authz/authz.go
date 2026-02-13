package authz

import (
	"raumzeitalpaka/app/database"
	"raumzeitalpaka/app/database/model"
	"raumzeitalpaka/app/database/query"
)

type Authorizer interface {
	HasRole(userId int, role model.Role) bool

	HasGroup(userId, groupId int) bool
	HasGroupRole(userId, groupId int, role model.Role) (hasGroup, hasRole bool)
	HasEventRole(userId, eventId int, role model.Role) (hasGroup, hasRole bool)
}

type DatabaseAuthz struct {
	UserHasGroupRole query.UserHasGroupRole
	UserHasEventRole query.UserHasEventRole
	UserHasRole      query.UserHasRole
}

func NewDatabaseAuthz(database *database.Database) *DatabaseAuthz {
	return &DatabaseAuthz{
		UserHasGroupRole: database.Queries.UserHasGroupRole,
		UserHasEventRole: database.Queries.UserHasEventRole,
		UserHasRole:      database.Queries.UserHasRole,
	}
}

func (d *DatabaseAuthz) HasRole(userId int, role model.Role) bool {
	resp, err := d.UserHasRole.Query(query.UserHasRoleRequest{UserId: userId, Role: role})
	if err != nil {
		return false
	}
	return resp.HasRole
}

func (d *DatabaseAuthz) HasGroup(userId, groupId int) bool {
	return true
	//resp, err := d.UserHasGroupRole.Query(query.UserHasGroupRoleRequest{
	//	UserId:  userId,
	//	GroupId: groupId,
	//	Role:    model.RoleParticipant,
	//})
	//if err != nil {
	//	return false
	//}
	//return resp.HasGroup
}

func (d *DatabaseAuthz) HasGroupRole(userId, groupId int, role model.Role) (hasGroup, hasRole bool) {
	resp, err := d.UserHasRole.Query(query.UserHasRoleRequest{
		UserId: userId,
		Role:   role,
	})
	if err != nil {
		return false, false
	}
	return true, resp.HasRole
}

func (d *DatabaseAuthz) HasEventRole(userId, eventId int, role model.Role) (hasGroup, hasRole bool) {
	resp, err := d.UserHasRole.Query(query.UserHasRoleRequest{
		UserId: userId,
		Role:   role,
	})
	if err != nil {
		return false, false
	}
	return true, resp.HasRole
}
