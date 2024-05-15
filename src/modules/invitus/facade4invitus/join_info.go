package facade4invitus

import (
	"context"
	"errors"
	"fmt"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/briefs4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/dal4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/invitus/models4invitus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/strongo/validation"
	"strconv"
	"time"
)

// JoinInfoRequest request
type JoinInfoRequest struct {
	InviteID string `json:"inviteID"` // InviteDto ID
	Pin      string `json:"pin"`
}

// Validate validates request
func (v *JoinInfoRequest) Validate() error {
	if v.InviteID == "" {
		return validation.NewErrRecordIsMissingRequiredField("id")
	}
	if v.Pin == "" {
		return validation.NewErrRequestIsMissingRequiredField("pin")
	}
	if _, err := strconv.Atoi(v.Pin); err != nil {
		return validation.NewErrBadRequestFieldValue("pin", "%pin is expected to be an integer")
	}
	return nil
}

type InviteInfo struct {
	Created time.Time                 `json:"created"`
	Status  string                    `json:"status"`
	From    models4invitus.InviteFrom `json:"from"`
	To      *models4invitus.InviteTo  `json:"to"`
	Message string                    `json:"message,omitempty"`
}

func (v InviteInfo) Validate() error {
	if v.Status == "" {
		return validation.NewErrRecordIsMissingRequiredField("status")
	}
	if v.Created.IsZero() {
		return validation.NewErrRecordIsMissingRequiredField("created")
	}
	if err := v.From.Validate(); err != nil {
		return validation.NewErrBadRecordFieldValue("from", err.Error())
	}
	if err := v.To.Validate(); err != nil {
		return validation.NewErrBadRecordFieldValue("to", err.Error())
	}
	return nil
}

// JoinInfoResponse response
type JoinInfoResponse struct {
	Team   models4invitus.InviteTeam                           `json:"team"`
	Invite InviteInfo                                          `json:"invite"`
	Member *dbmodels.DtoWithID[*briefs4contactus.ContactBrief] `json:"member"`
}

func (v JoinInfoResponse) Validated() error {
	if err := v.Team.Validate(); err != nil {
		return validation.NewErrBadRecordFieldValue("team", err.Error())
	}
	if err := v.Invite.Validate(); err != nil {
		return validation.NewErrBadRecordFieldValue("team", err.Error())
	}
	if nil == v.Member {
		return validation.NewErrRecordIsMissingRequiredField("member")
	}
	if err := v.Member.Validate(); err != nil {
		return validation.NewErrBadRecordFieldValue("member", err.Error())
	}
	return nil
}

// GetTeamJoinInfo return join info
func GetTeamJoinInfo(ctx context.Context, request JoinInfoRequest) (response JoinInfoResponse, err error) {
	if err = request.Validate(); err != nil {
		return
	}
	db := facade.GetDatabase(ctx)

	var inviteDto *models4invitus.InviteDto
	inviteDto, _, err = GetInviteByID(ctx, db, request.InviteID)
	if err != nil {
		err = fmt.Errorf("failed to get invite record by ID=%s: %w", request.InviteID, err)
		return
	}
	if inviteDto == nil {
		err = errors.New("invite record not found by ID: " + request.InviteID)
		return
	}

	if inviteDto.Pin != request.Pin {
		err = fmt.Errorf("%v: %w",
			validation.NewErrBadRequestFieldValue("pin", "invalid pin"),
			facade.ErrForbidden,
		)
		return
	}
	var member dal4contactus.ContactEntry
	if inviteDto.To.MemberID != "" {
		member = dal4contactus.NewContactEntry(inviteDto.TeamID, inviteDto.To.MemberID)
		db := facade.GetDatabase(ctx)
		if err = db.Get(ctx, member.Record); err != nil {
			err = fmt.Errorf("failed to get team member's contact record: %w", err)
			return
		}
	}
	response.Team = inviteDto.Team
	response.Team.ID = inviteDto.TeamID
	response.Invite.Status = inviteDto.Status
	response.Invite.Created = inviteDto.CreatedAt
	response.Invite.From = inviteDto.From
	response.Invite.To = inviteDto.To
	response.Invite.Message = inviteDto.Message
	if inviteDto.To.MemberID != "" {
		response.Member = &dbmodels.DtoWithID[*briefs4contactus.ContactBrief]{
			ID:   inviteDto.To.MemberID,
			Data: &member.Data.ContactBrief,
		}
	}
	return response, nil
}
