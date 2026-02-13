package oidc

import "testing"

func TestIsOrga_Positive(t *testing.T) {
	orgaGroups := []string{
		"org_event_hamburg",
		"org_event_berlin",
		"org_event_koeln",
	}

	for _, group := range orgaGroups {
		if !isOrga(group) {
			t.Fail()
		}
	}
}

func TestIsOrga_Negative(t *testing.T) {
	orgaGroups := []string{"uffd_access", "uffd_admin", "app_bbb_access", "app_pad_access", "app_forms_access", "app_pretalx_access", "app_presign_access", "app_cloud_access", "app_git_access", "app_ticket_access", "app_matrix_access", "app_monitoring_access", "app_backup_access", "uffd_signup", "app_netbox_access", "it-crew", "app_dns_access", "app_status_access", "app_dashboard_access", "app_termin_access", "uffd_login", "app_wiki_access", "app_wiki_user", "app_wiki_mentoring", "app_wiki_lab", "app_wiki_event", "app_wiki_admin", "administrators", "org_event_hamburg_mentorinnen", "org_event_hamburg_teili", "event_hamburg", "event_hamburg2024", "app_dashboard_editor", "app_dashboard_admin", "app_uptime_access", "org_event_berlin_mentorinnen", "event_berlin", "event_berlin2024", "app_wiki_hamburg", "org_event_koeln_mentorinnen", "org_event_koeln_teili", "event_hamburg2025", "app_proxmox_admin", "app_proxmox_access", "app_pretix_user", "git_group_17_50", "app_crm", "org_lab_hamburg", "org_lab_hamburg_lead", "org_lab_hamburg_menti", "git_group_256_30", "event_koeln2025", "git_group_260_50", "git_group_280_30", "git_group_280_50", "event_berlin2025", "git_group_13_50", "git_group_12_30", "git_group_12_50", "git_group_320_30", "git_group_320_50", "app_meet_control", "app_meet_changeurl", "app_raumzeitalpaka"}

	for _, group := range orgaGroups {
		if isOrga(group) {
			t.Fail()
		}
	}
}
