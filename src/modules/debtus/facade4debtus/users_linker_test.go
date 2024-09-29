package facade4debtus

import (
	"github.com/dal-go/dalgo/dal"
	dal4contactus2 "github.com/sneat-co/sneat-core-modules/contactus/dal4contactus"
	"github.com/sneat-co/sneat-core-modules/userus/dal4userus"
	"github.com/sneat-co/sneat-core-modules/userus/dbo4userus"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/gae_app/debtstracker/dtmocks"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/models4debtus"
	"github.com/sneat-co/sneat-go-core/facade"
	"testing"

	"context"
)

func TestUsersLinker_LinkUsersWithinTransaction(t *testing.T) {
	t.Skip("TODO: fix")
	c := context.Background()
	dtmocks.SetupMocks(c)

	usersLinker := &usersLinker{}

	const spaceID = "s3"

	contactusSpace := dal4contactus2.NewContactusSpaceEntry(spaceID)

	var (
		err                                        error
		inviterUser, invitedUser                   dbo4userus.UserEntry
		inviterContact, invitedContact             dal4contactus2.ContactEntry
		inviterDebtusContact, invitedDebtusContact models4debtus.DebtusSpaceContactEntry
	)

	if inviterUser, err = dal4userus.GetUserByID(c, nil, "1"); err != nil {
		t.Error("Failed to get inviter user", err)
		return
	}

	if invitedUser, err = dal4userus.GetUserByID(c, nil, "3"); err != nil {
		t.Error("Failed to get invited user", err)
		return
	}

	if inviterDebtusContact, err = GetDebtusSpaceContactByID(c, nil, spaceID, "6"); err != nil {
		t.Error("Failed to get inviter user", err)
		return
	}

	if inviterDebtusContact.Data.CounterpartyUserID != "" {
		t.Error("inviterDebtusContact.CounterpartyUserID != 0")
	}

	if inviterDebtusContact.Data.CounterpartyContactID != "" {
		t.Error("inviterDebtusContact.CounterpartyContactID != 0")
	}

	err = facade.RunReadwriteTransaction(c, func(tctx context.Context, tx dal.ReadwriteTransaction) (err error) {
		usersLinker = newUsersLinker(&usersLinkingDbChanges{
			inviter: &userLinkingParty{
				user:          inviterUser,
				contact:       inviterContact,
				debtusContact: inviterDebtusContact,
			},
			invited: &userLinkingParty{
				user:          invitedUser,
				contact:       invitedContact,
				debtusContact: invitedDebtusContact,
			},
		})
		if err = usersLinker.linkUsersWithinTransaction(tctx, tx, "unit-test:1"); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		t.Error("Unexpected error:", err)
		return
	}

	if len(usersLinker.changes.Records()) == 0 {
		t.Error("len(usersLinker.changes.EntityHolders()) == 0")
		return
	}

	invitedContact = usersLinker.changes.invited.contact
	inviterContact = usersLinker.changes.inviter.contact
	invitedDebtusContact = usersLinker.changes.invited.debtusContact
	inviterDebtusContact = usersLinker.changes.inviter.debtusContact
	invitedUser = usersLinker.changes.invited.user
	inviterUser = usersLinker.changes.inviter.user

	if invitedDebtusContact.ID == "" {
		t.Error("invitedDebtusContact.ContactID == 0")
		return
	}

	if invitedDebtusContact.ID == inviterDebtusContact.ID {
		t.Errorf("invitedDebtusContact.ContactID == inviterDebtusContact.ContactID: %s", invitedDebtusContact.ID)
	}

	if invitedDebtusContact.Data == nil {
		t.Error("invitedDebtusContact.DebtusSpaceContactDbo == nil")
		return
	}

	if invitedContact.Data.UserID == "" {
		t.Error("invitedDebtusContact.UserID == 0")
		return
	}

	if invitedContact.Data.UserID != invitedUser.ID {
		t.Errorf("invitedDebtusContact.UserID == invitedUser.ContactID : %s != %s", invitedContact.Data.UserID, invitedUser.ID)
		return
	}

	if invitedDebtusContact.Data.CounterpartyUserID == "" {
		t.Error("invitedDebtusContact.CounterpartyUserID == 0")
		return
	}

	if invitedDebtusContact.Data.CounterpartyContactID == "" {
		t.Error("invitedDebtusContact.CounterpartyContactID == 0")
		return
	}

	if invitedDebtusContact.Data.CounterpartyUserID != inviterUser.ID {
		t.Errorf("invitedDebtusContact.CounterpartyUserID != inviterUser.ContactID : %s != %s", invitedDebtusContact.Data.CounterpartyUserID, inviterUser.ID)
		return
	}

	if invitedDebtusContact.Data.CounterpartyContactID != inviterDebtusContact.ID {
		t.Errorf("invitedDebtusContact.CounterpartyContactID != inviterDebtusContact.ContactID : %s != %s", invitedDebtusContact.Data.CounterpartyContactID, inviterDebtusContact.ID)
		return
	}

	if inviterDebtusContact.Data.CounterpartyUserID == "" {
		t.Error("inviterDebtusContact.CounterpartyUserID == 0")
		return
	}

	if inviterDebtusContact.Data.CounterpartyContactID == "" {
		t.Error("inviterDebtusContact.CounterpartyContactID == 0")
		return
	}

	if inviterDebtusContact.Data.CounterpartyUserID != invitedUser.ID {
		t.Errorf("inviterDebtusContact.CounterpartyUserID != invitedUser.ContactID : %s != %s", inviterDebtusContact.Data.CounterpartyUserID, invitedUser.ID)
		return
	}

	if inviterDebtusContact.Data.CounterpartyContactID != invitedDebtusContact.ID {
		t.Errorf("inviterDebtusContact.CounterpartyContactID != invitedDebtusContact.ContactID : %s != %s", inviterDebtusContact.Data.CounterpartyContactID, invitedDebtusContact.ID)
		return
	}

	if invitedContact.Data.Names.UserName != "" && invitedContact.Data.Names.UserName == inviterDebtusContact.Data.NameFields.UserName {
		t.Errorf("invitedDebtusContact.Username == inviterDebtusContact.Username: %v", invitedDebtusContact.Data.UserName)
		return
	}

	if invitedDebtusContact.Data.FirstName != "" && invitedDebtusContact.Data.FirstName == inviterDebtusContact.Data.FirstName {
		t.Errorf("invitedDebtusContact.FirstName == inviterDebtusContact.FirstName: %v", invitedDebtusContact.Data.FirstName)
		return
	}

	if invitedDebtusContact.Data.LastName != "" && invitedDebtusContact.Data.LastName == inviterDebtusContact.Data.LastName {
		t.Errorf("invitedDebtusContact.LastName == inviterDebtusContact.LastName: %v", invitedDebtusContact.Data.LastName)
		return
	}

	if invitedDebtusContact.Data.NickName != "" && invitedDebtusContact.Data.NickName == inviterDebtusContact.Data.NickName {
		t.Errorf("invitedDebtusContact.NickName == inviterDebtusContact.NickName: %v", invitedDebtusContact.Data.NickName)
		return
	}

	if invitedDebtusContact.Data.ScreenName != "" && invitedDebtusContact.Data.ScreenName == inviterDebtusContact.Data.ScreenName {
		t.Errorf("invitedDebtusContact.ScreenName == inviterDebtusContact.ScreenName: %v", invitedDebtusContact.Data.ScreenName)
		return
	}

	var isInvitedUserHasInvitedContact bool

	for invitedUserContactID := range contactusSpace.Data.Contacts {
		if invitedUserContactID == invitedDebtusContact.ID {
			isInvitedUserHasInvitedContact = true
			break
		}
	}

	if !isInvitedUserHasInvitedContact {
		t.Error("Invited user missing invited Contact in the CounterpartiesJson")
		return
	}
}
