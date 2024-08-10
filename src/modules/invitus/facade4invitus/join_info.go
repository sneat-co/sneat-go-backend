package facade4invitus

import (
	"context"
	"errors"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/briefs4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/dal4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/invitus/dbo4invitus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/strongo/validation"
	"strconv"
	"time"
)

// JoinInfoRequest request
type JoinInfoRequest struct {
	InviteID string `json:"inviteID"` // InviteDbo ID
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
	Created time.Time              `json:"created"`
	Status  string                 `json:"status"`
	From    dbo4invitus.InviteFrom `json:"from"`
	To      *dbo4invitus.InviteTo  `json:"to"`
	Message string                 `json:"message,omitempty"`
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
	Space  dbo4invitus.InviteSpace                             `json:"space"`
	Invite InviteInfo                                          `json:"invite"`
	Member *dbmodels.DtoWithID[*briefs4contactus.ContactBrief] `json:"member"`
}

func (v JoinInfoResponse) Validated() error {
	if err := v.Space.Validate(); err != nil {
		return validation.NewErrBadRecordFieldValue("space", err.Error())
	}
	if err := v.Invite.Validate(); err != nil {
		return validation.NewErrBadRecordFieldValue("space", err.Error())
	}
	if nil == v.Member {
		return validation.NewErrRecordIsMissingRequiredField("member")
	}
	if err := v.Member.Validate(); err != nil {
		return validation.NewErrBadRecordFieldValue("member", err.Error())
	}
	return nil
}

// GetSpaceJoinInfo return join info
func GetSpaceJoinInfo(ctx context.Context, request JoinInfoRequest) (response JoinInfoResponse, err error) {
	if err = request.Validate(); err != nil {
		return
	}
	var db dal.DB
	if db, err = facade.GetDatabase(ctx); err != nil {
		return
	}

	var inviteDto *dbo4invitus.InviteDbo
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
		member = dal4contactus.NewContactEntry(inviteDto.SpaceID, inviteDto.To.MemberID)
		if err = db.Get(ctx, member.Record); err != nil {
			err = fmt.Errorf("failed to get team member's contact record: %w", err)
			return
		}
	}
	response.Space = inviteDto.Space
	response.Space.ID = inviteDto.SpaceID
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
