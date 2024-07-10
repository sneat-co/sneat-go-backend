package facade4invitus

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/briefs4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/invitus/dbo4invitus"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/sneat-co/sneat-go-core/models/dbprofile"
	"github.com/strongo/random"
	"github.com/strongo/validation"
	"net/mail"
	"strings"
	"time"
)

func NewInviteKey(inviteID string) *dal.Key {
	return dal.NewKeyWithID(InvitesCollection, inviteID)
}

var randomInviteID = func() string {
	return random.ID(6)
}

var randomPinCode = func() string {
	return random.Digits(4)
}

// FailedToSendEmail error message
const FailedToSendEmail = "failed to send email"

func createInviteForMember(
	ctx context.Context,
	tx dal.ReadwriteTransaction,
	uid string,
	remoteClient dbmodels.RemoteClientInfo,
	space dbo4invitus.InviteSpace,
	from dbo4invitus.InviteFrom,
	to dbo4invitus.InviteToMember,
	composeOnly bool,
	inviterUserID,
	message string,
	toAvatar *dbprofile.Avatar,
) (id string, personalInvite *dbo4invitus.PersonalInviteDbo, err error) {
	if err = space.Validate(); err != nil {
		err = fmt.Errorf("parameter 'space' is not valid: %w", err)
		return
	}
	if err = to.Validate(); err != nil {
		err = fmt.Errorf("parameter 'to' is not valid: %w", err)
		return
	}
	if err = from.Validate(); err != nil {
		err = fmt.Errorf("parameter 'from' is not valid: %w", err)
		return
	}
	spaceID := space.ID
	if spaceID == "" {
		err = validation.NewErrRecordIsMissingRequiredField("space.InviteID")
		return
	}
	space.ID = ""
	if space.Type == "family" && space.Title != "" {
		space.Title = ""
	}
	var toAddress *mail.Address
	if to.Address != "" {
		toAddress, err = mail.ParseAddress(to.Address)
		if err != nil {
			err = fmt.Errorf("failed to parse to.Address: %w", err)
			return
		}
	}
	var toAddressLower string
	if toAddress != nil {
		toAddressLower = strings.ToLower(toAddress.Address)
	}
	from.UserID = uid
	personalInvite = &dbo4invitus.PersonalInviteDbo{
		InviteDbo: dbo4invitus.InviteDbo{
			Status:  "active",
			Pin:     randomPinCode(),
			SpaceID: spaceID,
			InviteBase: dbo4invitus.InviteBase{
				Type:    "personal",
				Channel: to.Channel,
				From:    from, // TODO: get user email
				To: &dbo4invitus.InviteTo{
					InviteContact: to.InviteContact,
				},
				ComposeOnly: composeOnly,
			},
			CreatedAt: time.Now(),
			Created: dbmodels.CreatedInfo{
				Client: remoteClient,
			},
			Space:   space,
			Message: message,
			Roles:   []string{"contributor"},
		},
		Address:         toAddressLower,
		ToSpaceMemberID: briefs4contactus.GetFullContactID(spaceID, to.MemberID),
		ToAvatar:        toAvatar,
	}
	id = randomInviteID()
	inviteKey := NewInviteKey(id)
	if err = personalInvite.Validate(); err != nil {
		err = fmt.Errorf("personal invite record data are not valid: %w", err)
		return
	}
	inviteRecord := dal.NewRecordWithData(inviteKey, personalInvite)
	if err = tx.Insert(ctx, inviteRecord); err != nil {
		err = fmt.Errorf("failed to insert a new invite record into database: %w", err)
		return
	}
	return
}
