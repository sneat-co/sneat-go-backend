package facade4invitus

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/dal4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/invitus/dbo4invitus"
	"github.com/sneat-co/sneat-go-backend/src/modules/spaceus/dto4spaceus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/strongo/validation"
)

// InviteMemberRequest is a request DTO
type InviteMemberRequest struct {
	dto4spaceus.SpaceRequest
	RemoteClient dbmodels.RemoteClientInfo `json:"remoteClient"`

	To    dbo4invitus.InviteToMember `json:"to"`
	Roles []string                   `json:"roles,omitempty"`
	//
	Send    bool   `json:"send,omitempty"`
	Message string `json:"message,omitempty"`
}

const maxMessageSize = 1000

// Validate returns error if not valid
func (v InviteMemberRequest) Validate() error {
	if err := v.SpaceRequest.Validate(); err != nil {
		return err
	}
	//if err := v.From.Validate(); err != nil {
	//	return validation.NewErrBadRequestFieldValue("from", err.Error())
	//}
	if err := v.To.Validate(); err != nil {
		return validation.NewErrBadRequestFieldValue("to", err.Error())
	}
	if len(v.Message) > maxMessageSize {
		return validation.NewErrBadRequestFieldValue("message", fmt.Sprintf("message length limit is %d characters max", maxMessageSize))
	}
	if v.To.Channel != "email" && v.Send {
		return fmt.Errorf("%w: at the moment invites can be sent only by email, channel='%s'", facade.ErrBadRequest, v.To.Channel)
	}
	return nil
}

// CreateInviteResponse is a response DTO
type CreateInviteResponse struct {
	Invite dbo4invitus.InviteBrief `json:"invite"`
}

// CreateOrReuseInviteForMember creates or reuses an invitation for a member
func CreateOrReuseInviteForMember(ctx context.Context, userCtx facade.UserContext, request InviteMemberRequest) (response CreateInviteResponse, err error) {
	if err = request.Validate(); err != nil {
		err = fmt.Errorf("invalid request: %w", err)
		return
	}
	err = dal4contactus.RunContactusSpaceWorker(ctx, userCtx, request.SpaceRequest,
		func(ctx context.Context, tx dal.ReadwriteTransaction, params *dal4contactus.ContactusSpaceWorkerParams) (err error) {
			if err = tx.Get(ctx, params.SpaceModuleEntry.Record); err != nil {
				return
			}
			userID := params.UserID()
			fromContactID, fromContactBrief := params.SpaceModuleEntry.Data.GetContactBriefByUserID(userID)

			if fromContactBrief == nil {
				// TODO: Should return specific error so we can return HTTP 401
				return fmt.Errorf("current user does not belong to the space: user={id=%s}, space={id=%s}", userID, params.Space.ID)
			}
			var (
				inviteID       string
				personalInvite *dbo4invitus.PersonalInviteDbo
			)

			fromContact := dal4contactus.NewContactEntry(request.SpaceID, fromContactID)
			if err = tx.Get(ctx, fromContact.Record); err != nil {
				return err
			}
			memberInviteBrief := fromContact.Data.GetInviteBriefByChannelAndToMemberID(request.To.Channel, request.To.MemberID)
			if memberInviteBrief != nil {
				personalInvite, _, err = GetPersonalInviteByID(ctx, tx, memberInviteBrief.ID)
				if err != nil {
					if dal.IsNotFound(err) {
						err = nil
					} else {
						return err
					}
				}
				if personalInvite.Status == "active" {
					inviteID = memberInviteBrief.ID
				} else {
					personalInvite = nil
					inviteBriefs := make([]*dbo4invitus.MemberInviteBrief, 0, len(fromContact.Data.Invites)-1)
					for _, mi := range fromContact.Data.Invites {
						if mi.ID != memberInviteBrief.ID {
							inviteBriefs = append(inviteBriefs, mi)
						}
					}
					fromContact.Data.Invites = inviteBriefs
				}
			}
			if personalInvite == nil {
				inviteID, personalInvite, err =
					createPersonalInvite(ctx, tx, params.UserID(), request, params, fromContact)
				if err != nil {
					return fmt.Errorf("failed to create personal invite record: %w", err)
				}
			}
			response.Invite = dbo4invitus.NewInviteBriefFromDto(inviteID, personalInvite.InviteDbo)
			if !request.Send {
				response.Invite.Pin = personalInvite.Pin
			}
			return err
		},
	)
	response.Invite.To = nil
	response.Invite.From = nil
	return response, err
}

func createPersonalInvite(
	ctx context.Context,
	tx dal.ReadwriteTransaction,
	uid string,
	request InviteMemberRequest,
	param *dal4contactus.ContactusSpaceWorkerParams,
	fromMember dal4contactus.ContactEntry,
) (
	inviteID string, personalInvite *dbo4invitus.PersonalInviteDbo, err error,
) {

	toMember := param.SpaceModuleEntry.Data.Contacts[request.To.MemberID]
	if toMember == nil {
		err = fmt.Errorf("space has no 'to' member with id=%s", request.To.MemberID)
		return
	}
	request.To.Title = toMember.GetTitle()
	from := dbo4invitus.InviteFrom{
		InviteContact: dbo4invitus.InviteContact{
			UserID:   uid,
			MemberID: fromMember.ID,
			Title:    fromMember.Data.GetTitle(),
		},
	}
	to := request.To
	to.Title = toMember.GetTitle()
	inviteSpace := dbo4invitus.InviteSpace{
		ID:    request.SpaceID,
		Type:  param.Space.Data.Type,
		Title: param.Space.Data.Title,
	}
	inviteID, personalInvite, err = createInviteForMember(
		ctx,
		tx,
		uid,
		request.RemoteClient,
		inviteSpace,
		from,
		to,
		!request.Send,
		uid,
		request.Message,
		toMember.Avatar,
	)
	if err != nil {
		err = fmt.Errorf("failed to create an invite record for a member: %w", err)
		return "", nil, err
	}
	if request.Send {
		if personalInvite.MessageID, err = sendInviteEmail(ctx, inviteID, personalInvite); err != nil {
			err = fmt.Errorf("%s: %w", FailedToSendEmail, err)
			return inviteID, personalInvite, err
		}
		inviteKey := NewInviteKey(inviteID)
		if err = tx.Update(ctx, inviteKey,
			[]dal.Update{
				{Field: "messageId", Value: personalInvite.MessageID},
			}); err != nil {
			err = fmt.Errorf("failed to update invite record with message ContactID: %w", err)
			return inviteID, personalInvite, err
		}
	}
	fromMember.Data.Invites = append(fromMember.Data.Invites, &dbo4invitus.MemberInviteBrief{
		ID:         inviteID,
		To:         *personalInvite.To,
		CreateTime: personalInvite.CreatedAt,
	})
	memberKey := dal4contactus.NewContactKey(request.SpaceID, fromMember.ID)
	if err = tx.Update(ctx, memberKey, []dal.Update{
		{Field: "invites", Value: fromMember.Data.Invites},
	}); err != nil {
		err = fmt.Errorf("failed to add invite brief into member record: %w", err)
		return inviteID, personalInvite, err
	}
	return inviteID, personalInvite, err
}
