package api4invitus

import (
	"github.com/sneat-co/sneat-go-core/module"
	"net/http"
)

// RegisterHttpRoutes registers invites routes
func RegisterHttpRoutes(handle module.HTTPHandleFunc) {
	handle(http.MethodPost, "/v0/invites/create_invite_for_member", httpPostCreateOrReuseInviteForMember)
	handle(http.MethodGet, "/v0/invites/invite_link_for_member", httpGetOrCreateInviteLink)
	handle(http.MethodGet, "/v0/invites/personal_invite", httpGetPersonal)
	handle(http.MethodPost, "/v0/invites/create_mass_invite", httpPostCreateMassInvite)
	handle(http.MethodPost, "/v0/invites/accept_personal_invite", httpPostAcceptPersonalInvite)
	handle(http.MethodPost, "/v0/invites/reject_personal_invite", httpPostRejectPersonalInvite)
}
