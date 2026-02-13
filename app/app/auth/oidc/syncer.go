package oidc

import (
	"raumzeitalpaka/app/database"
	"raumzeitalpaka/app/database/command"
	"raumzeitalpaka/app/database/model"
	"raumzeitalpaka/app/database/query"
	"regexp"
	"strings"

	"go.uber.org/zap"
)

type Syncer interface {
	Sync(userId int, userName string, groups []string) error
}

type AlpakaSyncer struct {
	insertUser                     command.UpsertUser
	createGroup                    command.CreateGroup
	getGroupBySlug                 query.GetGroupBySlug
	updateManagedGroupsAssignments command.UpdateManagedGroupsAssignments
}

func NewAlpakaSyncer(db *database.Database) *AlpakaSyncer {
	return &AlpakaSyncer{
		insertUser:                     db.Commands.InsertUser,
		createGroup:                    db.Commands.CreateGroup,
		getGroupBySlug:                 db.Queries.GroupBySlug,
		updateManagedGroupsAssignments: db.Commands.UpdateManagedGroupsAssignments,
	}
}

func (s *AlpakaSyncer) Sync(userId int, userName string, groups []string) error {
	log := zap.L()
	log.Info("syncing user", zap.Int("user", userId), zap.Strings("groups", groups))
	//groupAssignments := make(map[int]model.Role)
	var isOrganizer = false
	var isMenti = false

	for _, oidcGroupStr := range groups {
		_, r, valid := parseGroup(oidcGroupStr, AlpakaRoleMapper)
		if !valid {
			log.Debug("user claims contained invalid group", zap.String("group", oidcGroupStr))
			continue
		}

		if r == model.RoleOrganizer {
			isOrganizer = true
			break
		}
		if !isMenti && r == model.RoleMentor {
			isMenti = true
		}

		//groupId, err := createGroupIfNew(s.createGroup, s.getGroupBySlug, command.CreateGroupRequest{
		//	Name: g,
		//	Slug: g,
		//})
		//if err != nil {
		//	log.Error("failed to create group", zap.Error(err))
		//	continue
		//}
		//
		//log.Info("user group parsed", zap.Int("id", groupId), zap.String("group", g), zap.Any("role", r))
		//
		//role, exists := groupAssignments[groupId]
		//if !exists {
		//	groupAssignments[groupId] = r
		//} else {
		//	groupAssignments[groupId] = model.HigherRole(role, r)
		//}
	}

	var role = model.RoleParticipant
	if isMenti {
		role = model.RoleMentor
	}
	if isOrganizer {
		role = model.RoleOrganizer
	}

	log.Info("upserting user", zap.Int("user", userId), zap.Any("role", role))
	_, err := s.insertUser.Execute(command.UpsertUserRequest{
		ID:           userId,
		LoginName:    userName,
		PasswordHash: "",
		Role:         role,
	})
	if err != nil {
		return err
	}

	//err := s.updateManagedGroupsAssignments.Execute(command.UpdateManagedGroupsAssignmentsRequest{
	//	UserId:      userId,
	//	Assignments: groupAssignmentsFrom(groupAssignments),
	//})
	//if err != nil {
	//	return err
	//}

	return nil
}

// alpaka uffd groups:
//["uffd_access", "uffd_admin", "app_bbb_access", "app_pad_access", "app_forms_access", "app_pretalx_access", "app_presign_access", "app_cloud_access", "app_git_access", "app_ticket_access", "app_matrix_access", "app_monitoring_access", "app_backup_access", "uffd_signup", "app_netbox_access", "org_event_hamburg", "it-crew", "app_dns_access", "app_status_access", "app_dashboard_access", "app_termin_access", "uffd_login", "app_wiki_access", "app_wiki_user", "app_wiki_mentoring", "app_wiki_lab", "app_wiki_event", "app_wiki_admin", "administrators", "org_event_hamburg_mentorinnen", "org_event_hamburg_teili", "event_hamburg", "event_hamburg2024", "app_dashboard_editor", "app_dashboard_admin", "app_uptime_access", "org_event_berlin_mentorinnen", "event_berlin", "event_berlin2024", "app_wiki_hamburg", "org_event_koeln_mentorinnen", "org_event_koeln_teili", "event_hamburg2025", "app_proxmox_admin", "app_proxmox_access", "app_pretix_user", "git_group_17_50", "app_crm", "org_lab_hamburg", "org_lab_hamburg_lead", "org_lab_hamburg_menti", "git_group_256_30", "event_koeln2025", "git_group_260_50", "git_group_280_30", "git_group_280_50", "event_berlin2025", "git_group_13_50", "git_group_12_30", "git_group_12_50", "git_group_320_30", "git_group_320_50", "app_meet_control", "app_meet_changeurl", "app_raumzeitalpaka"]
// ["org_event_hamburg", "org_event_hamburg_mentorinnen", "org_event_hamburg_teili", "event_hamburg", "event_hamburg2024", "event_hamburg2025",
//  "org_event_berlin_mentorinnen", "event_berlin", "event_berlin2024", "event_berlin2025",
//  "org_event_koeln_mentorinnen", "org_event_koeln_teili",  "event_koeln2025"
//  "org_lab_hamburg", "org_lab_hamburg_lead", "org_lab_hamburg_menti",
//

var groupParserRx = regexp.MustCompile(`Event:([^:]+):(?:[0-9]+:)?([^:]+)$`)
var isOrgaRx = regexp.MustCompile(`org_event_([a-z]+)$`)

func isOrga(group string) bool {
	return isOrgaRx.MatchString(group)
}

func parseGroup(oidcGroup string, roleMapper map[string]model.Role) (groupName string, role model.Role, valid bool) {
	//oidcGroup = strings.ToLower(oidcGroup)
	if isOrga(oidcGroup) {
		return "unknown", model.RoleOrganizer, true
	} else {
		return "unknown", model.RoleParticipant, false
	}

	//valid = groupParserRx.MatchString(oidcGroup)
	//if !valid {
	//	return "", "", false
	//}
	//mapRole := func(r string) model.Role {
	//	role, match := roleMapper[r]
	//	if !match {
	//		return model.Role(r)
	//	}
	//	return role
	//}
	//
	//parts := groupParserRx.FindStringSubmatch(oidcGroup)
	//groupName = parts[1]
	//rawRole := parts[2]
	//role = mapRole(rawRole)
	//
	//return groupName, role, model.ValidRole(role)
}

func createGroupIfNew(createGroup command.CreateGroup, getGroupBySlug query.GetGroupBySlug, request command.CreateGroupRequest) (groupId int, err error) {
	groupId, err = createGroup.Execute(request)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			group, err2 := getGroupBySlug.Query(query.GetGroupBySlugRequest{Slug: request.Slug})
			if err2 != nil {
				return -1, err
			}

			return group.ID, nil
		}
	}

	return groupId, nil
}

func groupAssignmentsFrom(r map[int]model.Role) []command.GroupAssignment {
	assignments := make([]command.GroupAssignment, 0, len(r))
	for i, role := range r {
		assignments = append(assignments, command.GroupAssignment{
			GroupId: i,
			Role:    role,
		})
	}
	return assignments
}
