package oidc

import (
	"context"
	"raumzeitalpaka/app/database"
	"raumzeitalpaka/app/database/command"
	"raumzeitalpaka/app/database/model"
	"raumzeitalpaka/app/database/query"
	"regexp"
	"sync"

	"go.uber.org/zap"
)

type Syncer interface {
	Sync(userId int, userName, displayName string, groups []string) error
}

type AlpakaSyncer struct {
	insertUser                     command.UpsertUser
	createGroup                    command.CreateOrganisation
	getGroupBySlug                 query.GetOrganisationBySlug
	updateManagedGroupsAssignments command.UpdateManagedOrganisationAssignments
}

func NewAlpakaSyncer(db *database.Database) *AlpakaSyncer {
	return &AlpakaSyncer{
		insertUser:                     db.Commands.InsertUser,
		createGroup:                    db.Commands.CreateOrganisation,
		getGroupBySlug:                 db.Queries.OrganisationBySlug,
		updateManagedGroupsAssignments: db.Commands.UpdateManagedOrganisationAssignments,
	}
}

func (s *AlpakaSyncer) Sync(userId int, userName, displayName string, groups []string) error {
	log := zap.L()
	log.Info("syncing user", zap.Int("user", userId), zap.Strings("groups", groups))
	ctx := context.TODO()
	groupAssignments := make(map[int]model.Role)
	var isOrganizer = false
	var isMenti = false

	for _, oidcGroupStr := range groups {
		group, r, valid := parseGroup(oidcGroupStr, AlpakaRoleMapper)
		if !valid {
			log.Debug("user claims contained invalid group", zap.String("group", oidcGroupStr))
			continue
		}

		groupID, err := getOrCreateGroup(ctx, s.createGroup, s.getGroupBySlug, command.CreateOrganisationRequest{
			Name: group,
			Slug: group,
		})
		if err != nil {
			return err
		}
		log.Info("found group assignment", zap.String("group", group), zap.Int("groupID", groupID), zap.Any("role", r), zap.Bool("valid", valid))

		groupRole, exists := groupAssignments[groupID]
		if !exists || groupRole != model.RoleOrganizer {
			isOrganizer = true

			groupAssignments[groupID] = r
		}
	}

	orgs := make([]command.OrganisationAssignment, 0, len(groupAssignments))
	for groupID, role := range groupAssignments {
		orgs = append(orgs, command.OrganisationAssignment{
			OrganisationId: groupID,
			Role:           role,
		})
	}

	err := s.updateManagedGroupsAssignments.Execute(ctx, command.UpdateManagedOrganisationAssignmentsRequest{
		UserId:      userId,
		Assignments: orgs,
	})
	if err != nil {
		return err
	}

	var role = model.RoleParticipant
	if isMenti {
		role = model.RoleMentor
	}
	if isOrganizer {
		role = model.RoleOrganizer
	}

	log.Info("upserting user", zap.Int("user", userId), zap.Any("role", role))
	_, err = s.insertUser.Execute(context.TODO(), command.UpsertUserRequest{
		ID:           userId,
		LoginName:    userName,
		DisplayName:  displayName,
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
var isMentiRx = regexp.MustCompile(`org_event_([a-z]+)_mentorinnen$`)

func isOrga(group string) (string, bool) {
	match := isOrgaRx.MatchString(group)
	if !match {
		return "", false
	}

	matches := isOrgaRx.FindStringSubmatch(group)
	return matches[1], true
}

func isMenti(group string) (string, bool) {
	match := isMentiRx.MatchString(group)
	if !match {
		return "", false
	}

	matches := isMentiRx.FindStringSubmatch(group)
	return matches[1], true
}

func parseGroup(oidcGroup string, roleMapper map[string]model.Role) (groupName string, role model.Role, valid bool) {
	//oidcGroup = strings.ToLower(oidcGroup)
	//if isOrga(oidcGroup) {
	//	return "unknown", model.RoleOrganizer, true
	//} else {
	//	return "unknown", model.RoleParticipant, false
	//}

	groupName, isO := isOrga(oidcGroup)
	if isO {
		return groupName, model.RoleOrganizer, true
	}
	groupName, isM := isMenti(oidcGroup)
	if isM {
		return groupName, model.RoleMentor, true
	}
	return "unknown", model.RoleParticipant, false

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

var groupSyncMutex sync.Mutex

func getOrCreateGroup(ctx context.Context, createGroup command.CreateOrganisation, getGroupBySlug query.GetOrganisationBySlug, request command.CreateOrganisationRequest) (groupId int, err error) {
	groupSyncMutex.Lock()
	defer groupSyncMutex.Unlock()
	group, err2 := getGroupBySlug.Query(ctx, query.GetOrganisationBySlugRequest{Slug: request.Slug})
	if err2 == nil {
		return group.ID, err
	}

	groupId, err = createGroup.Execute(ctx, request)
	if err != nil {
		return -1, err
	}

	return groupId, nil
}

func groupAssignmentsFrom(r map[int]model.Role) []command.OrganisationAssignment {
	assignments := make([]command.OrganisationAssignment, 0, len(r))
	for i, role := range r {
		assignments = append(assignments, command.OrganisationAssignment{
			OrganisationId: i,
			Role:           role,
		})
	}
	return assignments
}
