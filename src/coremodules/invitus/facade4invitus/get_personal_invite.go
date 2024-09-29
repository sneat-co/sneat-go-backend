package facade4invitus

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/coremodules/contactus/briefs4contactus"
	"github.com/sneat-co/sneat-go-backend/src/coremodules/contactus/dal4contactus"
	"github.com/sneat-co/sneat-go-backend/src/coremodules/invitus/dbo4invitus"
	"github.com/sneat-co/sneat-go-backend/src/coremodules/spaceus/dto4spaceus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/validation"
)

// GetPersonalInviteRequest holds parameters for creating a personal invite
type GetPersonalInviteRequest struct {
	dto4spaceus.SpaceRequest
	InviteID string `json:"inviteID"`
}

// Validate validates request
func (v *GetPersonalInviteRequest) Validate() error {
	if err := v.SpaceRequest.Validate(); err != nil {
		return err
	}
	if v.InviteID == "" {
		return validation.NewErrRecordIsMissingRequiredField("invite")
	}
	//if len(v.InviteID) != 8 {
	//	return models2spotbuddies.NewErrBadRequestFieldValue("invite", "unexpected length of invite id")
	//}
	return nil
}

// PersonalInviteResponse holds response data for created personal invite
type PersonalInviteResponse struct {
	Invite  *dbo4invitus.PersonalInviteDbo            `json:"invite,omitempty"`
	Members map[string]*briefs4contactus.ContactBrief `json:"members,omitempty"`
}

func getPersonalInviteRecords(ctx context.Context, getter dal.ReadSession, params *dal4contactus.ContactusSpaceWorkerParams, inviteID, memberID string) (
	invite PersonalInviteEntry,
	member dal4contactus.ContactEntry,
	err error,
) {
	if inviteID == "" {
		err = validation.NewErrRequestIsMissingRequiredField("inviteID")
		return
	}
	invite = NewPersonalInviteEntry(inviteID)

	records := []dal.Record{invite.Record}
	if memberID != "" {
		member = dal4contactus.NewContactEntry(params.Space.ID, memberID)
		records = append(records, member.Record)
	}
	if err = params.GetRecords(ctx, getter, records...); err != nil {
		return
	}
	if !params.SpaceModuleEntry.Record.Exists() {
		err = validation.NewErrBadRequestFieldValue("spaceID",
			fmt.Sprintf("contactusSpace record not found by key=%v: record.Error=%v",
				params.SpaceModuleEntry.Key, params.SpaceModuleEntry.Record.Error()),
		)
		return
	}
	if !invite.Record.Exists() {
		err = validation.NewErrBadRequestFieldValue("inviteID",
			fmt.Sprintf("invite record not found in database by key=%v: record.Error=%v",
				invite.Record.Key(), invite.Record.Error()))
		return
	}
	if member.Record != nil && !member.Record.Exists() {
		err = validation.NewErrBadRequestFieldValue("memberID",
			fmt.Sprintf("member record not found in database by key=%v: record.Error=%v",
				member.Record.Key(), member.Record.Error()))
		return
	}
	return
}

// GetPersonal returns personal invite data
func GetPersonal(ctx context.Context, userCtx facade.UserContext, request GetPersonalInviteRequest) (response PersonalInviteResponse, err error) {
	if err = request.Validate(); err != nil {
		return
	}
	return response, dal4contactus.RunReadonlyContactusSpaceWorker(ctx, userCtx, request.SpaceRequest, func(ctx context.Context, tx dal.ReadTransaction, params *dal4contactus.ContactusSpaceWorkerParams) error {
		invite, _, err := getPersonalInviteRecords(ctx, tx, params, request.InviteID, "")
		if err != nil {
			return err
		}
		invite.Data.Pin = "" // Hide PIN code from visitor
		response = PersonalInviteResponse{
			Invite:  invite.Data,
			Members: make(map[string]*briefs4contactus.ContactBrief, len(params.SpaceModuleEntry.Data.Contacts)),
		}
		// TODO: Is this is a security breach in current implementation?
		//for id, contact := range contactusSpace.Data.Contacts {
		//	response.Members[id] = contact
		//}
		return nil
	})
}
