package api4contactus

import (
	"github.com/sneat-co/sneat-go-core/module"
	"net/http"
)

// RegisterHttpRoutes registers contact routes
func RegisterHttpRoutes(handle module.HTTPHandleFunc) {
	handle(http.MethodPost, "/v0/contactus/create_contact", httpPostCreateContact)
	handle(http.MethodDelete, "/v0/contactus/delete_contact", httpDeleteContact)
	handle(http.MethodPost, "/v0/contactus/set_contacts_status", httpSetContactStatus)
	handle(http.MethodPost, "/v0/contactus/update_contact", httpUpdateContact)
	handle(http.MethodPost, "/v0/contactus/archive_contact", httpPostArchiveContact)
	handle(http.MethodPost, "/v0/contactus/create_member", httpPostCreateMember)
	handle(http.MethodPost, "/v0/contactus/remove_team_member", httpPostRemoveSpaceMember)

	////
	//handle(http.MethodGet, "/v0/team/join_info", api4debtus.GetTeamJoinInfo)
	//handle(http.MethodPost, "/v0/team/join_team", api4debtus.JoinTeam)
	//handle(http.MethodPost, "/v0/team/refuse_to_join_team", api4debtus.RefuseToJoinTeam)
	//handle(http.MethodPost, "/v0/team/leave_team", api4debtus.LeaveTeam)
	//handle(http.MethodPost, "/v0/team/create_member", api4debtus.AddMember)
	//handle(http.MethodPost, "/v0/team/add_metric", api4debtus.AddMetric)
	//handle(http.MethodPost, "/v0/team/remove_member", api4debtus.RemoveMember)
	//handle(http.MethodPost, "/v0/team/change_member_role", api4debtus.ChangeMemberRole)
	//handle(http.MethodPost, "/v0/team/remove_metrics", api4debtus.RemoveMetrics)
	//handle(http.MethodGet, "/v0/team", api4debtus.GetTeam)
}
