package gaedal

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/dtdal"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/facade"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/models"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/general"
	"github.com/strongo/log"
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

func (InviteDalGae) GetInvite(c context.Context, tx dal.ReadSession, inviteCode string) (invite models.Invite, err error) {
	if tx == nil {
		if tx, err = facade.GetDatabase(c); err != nil {
			return
		}
	}
	invite = models.NewInvite(inviteCode, nil)
	return invite, tx.Get(c, invite.Record)
}

// ClaimInvite claims invite by user - TODO compare with ClaimInvite2 and get rid of one of them
func (InviteDalGae) ClaimInvite(c context.Context, userID string, inviteCode, claimedOn, claimedVia string) (err error) {
	var db dal.DB
	if db, err = facade.GetDatabase(c); err != nil {
		return
	}
	err = db.RunReadwriteTransaction(c, func(tc context.Context, tx dal.ReadwriteTransaction) error {
		invite := models.NewInvite(inviteCode, nil)

		if err = tx.Get(tc, invite.Record); err != nil {
			return err
		}
		log.Debugf(c, "Invite found")
		// TODO: Check invite.For
		//invite.ClaimedCount += 1
		inviteClaim := models.NewInviteClaimWithoutID(models.NewInviteClaimData(inviteCode, userID, claimedOn, claimedVia))
		user := models.NewAppUser(userID, nil)
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
		log.Debugf(c, "inviteClaimKey.IntegerID(): %v", inviteClaimID)
		return DelayUpdateInviteClaimedCount(tc, inviteClaimID)
	})
	return
}

const (
	AUTO_GENERATE_INVITE_CODE = ""
	INVITE_CODE_LENGTH        = 5
	PERSONAL_INVITE           = 1
)

func (InviteDalGae) CreatePersonalInvite(ec strongoapp.ExecutionContext, userID string, inviteBy models.InviteBy, inviteToAddress, createdOnPlatform, createdOnID, related string) (models.Invite, error) {
	return createInvite(ec, models.InviteTypePersonal, userID, inviteBy, inviteToAddress, createdOnPlatform, createdOnID, INVITE_CODE_LENGTH, AUTO_GENERATE_INVITE_CODE, related, PERSONAL_INVITE)
}

func (InviteDalGae) CreateMassInvite(ec strongoapp.ExecutionContext, userID string, inviteCode string, maxClaimsCount int32, createdOnPlatform string) (invite models.Invite, err error) {
	invite, err = createInvite(ec, models.InviteTypePublic, userID, "", "", createdOnPlatform, "", uint8(len(inviteCode)), inviteCode, "", maxClaimsCount)
	return
}

func createInvite(ec strongoapp.ExecutionContext, inviteType models.InviteType, userID string, inviteBy models.InviteBy, inviteToAddress, createdOnPlatform, createdOnID string, inviteCodeLen uint8, inviteCode, related string, maxClaimsCount int32) (invite models.Invite, err error) {
	if inviteCode != AUTO_GENERATE_INVITE_CODE && !dtdal.InviteCodeRegex.Match([]byte(inviteCode)) {
		err = fmt.Errorf("Invalid invite code: %v", inviteCode)
		return
	}
	if related != "" && len(strings.Split(related, "=")) != 2 {
		panic(fmt.Sprintf("Invalid format for related: %v", related))
	}
	c := ec.Context()

	dtCreated := time.Now()
	invite = models.NewInvite(inviteCode, &models.InviteData{
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
	case models.InviteByEmail:
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
	case models.InviteBySms:
		var phoneNumber int64
		phoneNumber, err = strconv.ParseInt(inviteToAddress, 10, 64)
		if err != nil {
			return
		}
		invite.Data.ToPhoneNumber = phoneNumber
	}
	var db dal.DB
	if db, err = facade.GetDatabase(c); err != nil {
		return
	}
	err = db.RunReadwriteTransaction(c, func(tc context.Context, tx dal.ReadwriteTransaction) error {
		if inviteCode != AUTO_GENERATE_INVITE_CODE {
			//inviteKey = models.NewInviteKey(inviteCode)
		} else {
			for {
				if inviteCodeLen == 0 {
					inviteCodeLen = INVITE_CODE_LENGTH
				}
				inviteCode = dtdal.RandomCode(inviteCodeLen)
				existingInvite := models.NewInvite(inviteCode, nil)

				if err := tx.Get(c, existingInvite.Record); dal.IsNotFound(err) {
					//log.Debugf(c, "New invite code: %v", inviteCode)
					break
				} else {
					log.Warningf(c, "Already existing invite code: %v", inviteCode)
				}
			}
		}
		return tx.Set(c, invite.Record)
	}, nil)
	if err == nil {
		log.Infof(c, "Invite created with code: %v", inviteCode)
	} else {
		log.Errorf(c, "Failed to create invite with code: %v", err)
	}
	return
}

// ClaimInvite2 claims invite by user - TODO compare with ClaimInvite and get rid of one of them
func (InviteDalGae) ClaimInvite2(c context.Context, inviteCode string, invite models.Invite, claimedByUserID string, claimedOn, claimedVia string) (err error) {
	var db dal.DB
	if db, err = facade.GetDatabase(c); err != nil {
		return fmt.Errorf("failed to create database: %w", err)
	}

	err = db.RunReadwriteTransaction(c, func(tc context.Context, tx dal.ReadwriteTransaction) error {
		//userKey := models.NewAppUserKey(claimedByUserID)
		user := models.NewAppUser(claimedByUserID, nil)
		if err = tx.GetMulti(tc, []dal.Record{invite.Record, user.Record}); err != nil {
			return err
		}

		invite.Data.ClaimedCount += 1
		if invite.Data.MaxClaimsCount > 0 && invite.Data.ClaimedCount > invite.Data.MaxClaimsCount {
			return fmt.Errorf("invite.ClaimedCount > invite.MaxClaimsCount: %v > %v", invite.Data.ClaimedCount, invite.Data.MaxClaimsCount)
		}
		inviteClaimData := models.NewInviteClaimData(inviteCode, claimedByUserID, claimedOn, claimedVia)
		inviteClaim := models.NewInviteClaim(0, inviteClaimData)
		if err = tx.Insert(c, inviteClaim.Record); err != nil {
			return err
		}
		recordsToUpdate := []dal.Record{invite.Record}

		userChanged := updateUserContactDetails(user.Data, *invite.Data)

		if user.Data.DtAccessGranted.IsZero() {
			user.Data.DtAccessGranted = time.Now()
			userChanged = true
		}
		if invite.Data.MaxClaimsCount == 1 {
			user.Data.InvitedByUserID = invite.Data.CreatedByUserID
			userChanged = true
			counterpartyQuery := dal.From(models.DebtusContactsCollection).
				WhereField("UserID", dal.Equal, claimedByUserID).
				WhereField("CounterpartyUserID", dal.Equal, invite.Data.CreatedByUserID).
				Limit(1).
				SelectInto(models.NewDebtusContactRecord)

			counterpartyRecords, err := db.QueryAllRecords(c, counterpartyQuery)

			if err != nil {
				return fmt.Errorf("failed to load counterparty by CounterpartyUserID: %w", err)
			}
			if len(counterpartyRecords) == 0 {
				//counterpartyKey := NewContactIncompleteKey(tc)
				inviteCreator, err := facade.User.GetUserByID(c, tx, invite.Data.CreatedByUserID)
				if err != nil {
					return fmt.Errorf("ailed to get invite creator user: %w", err)
				}

				counterparty := models.NewDebtusContact("", models.NewDebtusContactData(claimedByUserID, models.ContactDetails{
					FirstName:    inviteCreator.Data.FirstName,
					LastName:     inviteCreator.Data.LastName,
					Username:     inviteCreator.Data.Username,
					EmailContact: inviteCreator.Data.EmailContact,
					PhoneContact: inviteCreator.Data.PhoneContact,
				}))

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
	//invite.ID = inviteClaimKey.StringID()
	return
}

func updateUserContactDetails(user *models.DebutsAppUserDataOBSOLETE, inviteData models.InviteData) (changed bool) {
	switch models.InviteBy(inviteData.Channel) {
	case models.InviteByEmail:
		changed = !user.EmailConfirmed
		user.SetEmail(inviteData.ToEmail, true)
	case models.InviteBySms:
		if inviteData.ToPhoneNumber != 0 {
			changed = !user.PhoneNumberConfirmed
			user.PhoneNumber = inviteData.ToPhoneNumber
			user.PhoneNumberConfirmed = true
		}
	}
	return
}
