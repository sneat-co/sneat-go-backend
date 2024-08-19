package gaedal

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/briefs4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/const4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/dal4contactus"
	models4contactus "github.com/sneat-co/sneat-go-backend/src/modules/contactus/dbo4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/gae_app/debtstracker/dtdal"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/gae_app/general"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/models4debtus"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/dal4userus"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/dbo4userus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/logus"
	"github.com/strongo/strongoapp"
	"strconv"
	"strings"
	"time"
)

type InviteDalGae struct {
}

func NewInviteDalGae() InviteDalGae {
	return InviteDalGae{}
}

func (InviteDalGae) GetInvite(c context.Context, tx dal.ReadSession, inviteCode string) (invite models4debtus.Invite, err error) {
	if tx == nil {
		if tx, err = facade.GetDatabase(c); err != nil {
			return
		}
	}
	invite = models4debtus.NewInvite(inviteCode, nil)
	return invite, tx.Get(c, invite.Record)
}

// ClaimInvite claims invite by user - TODO compare with ClaimInvite2 and get rid of one of them
func (InviteDalGae) ClaimInvite(c context.Context, userID string, inviteCode, claimedOn, claimedVia string) (err error) {
	err = facade.RunReadwriteTransaction(c, func(tc context.Context, tx dal.ReadwriteTransaction) error {
		invite := models4debtus.NewInvite(inviteCode, nil)

		if err = tx.Get(tc, invite.Record); err != nil {
			return err
		}
		logus.Debugf(c, "Invite found")
		// TODO: Check invite.For
		//invite.ClaimedCount += 1
		inviteClaim := models4debtus.NewInviteClaimWithoutID(models4debtus.NewInviteClaimData(inviteCode, userID, claimedOn, claimedVia))
		user := dbo4userus.NewUserEntry(userID)
		if err = tx.Get(tc, user.Record); err != nil {
			return err
		}
		user.Data.InvitedByUserID = invite.Data.CreatedByUserID

		if err = tx.Insert(c, inviteClaim.Record); err != nil {
			return fmt.Errorf("failed to insert invite claim: %w", err)
		}
		if err := tx.Set(tc, user.Record); err != nil {
			return fmt.Errorf("failed to save user: %w", err)
		}
		inviteClaimID := inviteClaim.Key.ID.(int64)
		logus.Debugf(c, "inviteClaimKey.IntegerID(): %v", inviteClaimID)
		return DelayUpdateInviteClaimedCount(tc, inviteClaimID)
	})
	return
}

const (
	AUTO_GENERATE_INVITE_CODE = ""
	INVITE_CODE_LENGTH        = 5
	PERSONAL_INVITE           = 1
)

func (InviteDalGae) CreatePersonalInvite(ec strongoapp.ExecutionContext, userID string, inviteBy models4debtus.InviteBy, inviteToAddress, createdOnPlatform, createdOnID, related string) (models4debtus.Invite, error) {
	return createInvite(ec, models4debtus.InviteTypePersonal, userID, inviteBy, inviteToAddress, createdOnPlatform, createdOnID, INVITE_CODE_LENGTH, AUTO_GENERATE_INVITE_CODE, related, PERSONAL_INVITE)
}

func (InviteDalGae) CreateMassInvite(ec strongoapp.ExecutionContext, userID string, inviteCode string, maxClaimsCount int32, createdOnPlatform string) (invite models4debtus.Invite, err error) {
	invite, err = createInvite(ec, models4debtus.InviteTypePublic, userID, "", "", createdOnPlatform, "", uint8(len(inviteCode)), inviteCode, "", maxClaimsCount)
	return
}

func createInvite(ec strongoapp.ExecutionContext, inviteType models4debtus.InviteType, userID string, inviteBy models4debtus.InviteBy, inviteToAddress, createdOnPlatform, createdOnID string, inviteCodeLen uint8, inviteCode, related string, maxClaimsCount int32) (invite models4debtus.Invite, err error) {
	if inviteCode != AUTO_GENERATE_INVITE_CODE && !dtdal.InviteCodeRegex.Match([]byte(inviteCode)) {
		err = fmt.Errorf("Invalid invite code: %v", inviteCode)
		return
	}
	if related != "" && len(strings.Split(related, "=")) != 2 {
		panic(fmt.Sprintf("Invalid format for related: %v", related))
	}
	c := ec.Context()

	dtCreated := time.Now()
	invite = models4debtus.NewInvite(inviteCode, &models4debtus.InviteData{
		Type:    string(inviteType),
		Channel: string(inviteBy),
		CreatedOn: general.CreatedOn{
			CreatedOnPlatform: createdOnPlatform,
			CreatedOnID:       createdOnID,
		},
		DtCreated:       dtCreated,
		CreatedByUserID: userID,
		Related:         related,
		MaxClaimsCount:  maxClaimsCount,
		DtActiveFrom:    dtCreated,
		DtActiveTill:    dtCreated.AddDate(100, 0, 0), // By default is active for 100 years
	})
	switch inviteBy {
	case models4debtus.InviteByEmail:
		if inviteToAddress == "" {
			panic("Emmail address is not supplied")
		}
		if strings.Index(inviteToAddress, "@") <= 0 || strings.Index(inviteToAddress, ".") <= 0 {
			panic("Invalid email address")
		}
		invite.Data.ToEmail = strings.ToLower(inviteToAddress)
		if inviteToAddress != strings.ToLower(inviteToAddress) {
			invite.Data.ToEmailOriginal = inviteToAddress
		}
	case models4debtus.InviteBySms:
		var phoneNumber int64
		phoneNumber, err = strconv.ParseInt(inviteToAddress, 10, 64)
		if err != nil {
			return
		}
		invite.Data.ToPhoneNumber = phoneNumber
	}
	err = facade.RunReadwriteTransaction(c, func(tc context.Context, tx dal.ReadwriteTransaction) error {
		if inviteCode != AUTO_GENERATE_INVITE_CODE {
			//inviteKey = models.NewInviteKey(inviteCode)
		} else {
			for {
				if inviteCodeLen == 0 {
					inviteCodeLen = INVITE_CODE_LENGTH
				}
				inviteCode = dtdal.RandomCode(inviteCodeLen)
				existingInvite := models4debtus.NewInvite(inviteCode, nil)

				if err := tx.Get(c, existingInvite.Record); dal.IsNotFound(err) {
					//logus.Debugf(c, "New invite code: %v", inviteCode)
					break
				} else {
					logus.Warningf(c, "Already existing invite code: %v", inviteCode)
				}
			}
		}
		return tx.Set(c, invite.Record)
	}, nil)
	if err == nil {
		logus.Infof(c, "Invite created with code: %v", inviteCode)
	} else {
		logus.Errorf(c, "Failed to create invite with code: %v", err)
	}
	return
}

// ClaimInvite2 claims invite by user - TODO compare with ClaimInvite and get rid of one of them
func (InviteDalGae) ClaimInvite2(c context.Context, inviteCode string, invite models4debtus.Invite, claimedByUserID string, claimedOn, claimedVia string) (err error) {
	var db dal.DB // Needed for query records outside of transaction
	if db, err = facade.GetDatabase(c); err != nil {
		return
	}
	err = facade.RunReadwriteTransaction(c, func(tc context.Context, tx dal.ReadwriteTransaction) error {
		//userKey := models.NewAppUserKeyOBSOLETE(claimedByUserID)
		user := dbo4userus.NewUserEntry(claimedByUserID)
		if err = tx.GetMulti(tc, []dal.Record{invite.Record, user.Record}); err != nil {
			return err
		}

		invite.Data.ClaimedCount += 1
		if invite.Data.MaxClaimsCount > 0 && invite.Data.ClaimedCount > invite.Data.MaxClaimsCount {
			return fmt.Errorf("invite.ClaimedCount > invite.MaxClaimsCount: %v > %v", invite.Data.ClaimedCount, invite.Data.MaxClaimsCount)
		}
		inviteClaimData := models4debtus.NewInviteClaimData(inviteCode, claimedByUserID, claimedOn, claimedVia)
		inviteClaim := models4debtus.NewInviteClaim(0, inviteClaimData)
		if err = tx.Insert(c, inviteClaim.Record); err != nil {
			return err
		}
		recordsToUpdate := []dal.Record{invite.Record}

		userChanged := updateUserContactDetails(user, *invite.Data)

		//if user.Data.DtAccessGranted.IsZero() {
		//	user.Data.DtAccessGranted = time.Now()
		//	userChanged = true
		//}
		if invite.Data.MaxClaimsCount == 1 {
			user.Data.InvitedByUserID = invite.Data.CreatedByUserID
			userChanged = true
			counterpartyQuery := dal.From(const4contactus.ContactsCollection).
				WhereField("UserID", dal.Equal, claimedByUserID).
				WhereField("CounterpartyUserID", dal.Equal, invite.Data.CreatedByUserID).
				Limit(1).
				SelectInto(models4debtus.NewDebtusContactRecord)

			var counterpartyRecords []dal.Record
			counterpartyRecords, err = db.QueryAllRecords(c, counterpartyQuery)

			if err != nil {
				return fmt.Errorf("failed to load counterparty by CounterpartyUserID: %w", err)
			}
			if len(counterpartyRecords) == 0 {
				//counterpartyKey := NewContactIncompleteKey(tc)
				inviterUser := dbo4userus.NewUserEntry(invite.Data.CreatedByUserID)
				if err = dal4userus.GetUser(c, tx, inviterUser); err != nil {
					return fmt.Errorf("ailed to get invite creator user: %w", err)
				}

				counterparty := dal4contactus.NewContactEntryWithData("", "", &models4contactus.ContactDbo{
					ContactBase: briefs4contactus.ContactBase{
						ContactBrief: briefs4contactus.ContactBrief{
							Names: inviterUser.Data.Names,
						},
					},
				})

				if err = tx.Insert(c, counterparty.Record); err != nil {
					return fmt.Errorf("failed to insert counterparty: %w", err)
				}
			}
		}

		if userChanged {
			recordsToUpdate = append(recordsToUpdate, user.Record)
		}

		err = tx.SetMulti(tc, recordsToUpdate)
		if err != nil {
			return err
		}

		return err
	})
	if err != nil {
		return
	}
	//invite.ContactID = inviteClaimKey.StringID()
	return
}

func updateUserContactDetails(user dbo4userus.UserEntry, inviteData models4debtus.InviteData) (changed bool) {
	switch models4debtus.InviteBy(inviteData.Channel) {
	case models4debtus.InviteByEmail:
		changed = !user.Data.EmailVerified
		panic("Not implemented")
		//user.SetEmail(inviteData.ToEmail, true)
	case models4debtus.InviteBySms:
		if inviteData.ToPhoneNumber != 0 {
			panic("not implemented")
			//changed = !user.PhoneNumberConfirmed
			//user.PhoneNumber = inviteData.ToPhoneNumber
			//user.PhoneNumberConfirmed = true
		}
	}
	return
}