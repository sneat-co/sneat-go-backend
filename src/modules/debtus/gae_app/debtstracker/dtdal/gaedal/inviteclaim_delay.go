package gaedal

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/core/queues"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/gae_app/debtstracker/dtdal"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/models4debtus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/delaying"
	"github.com/strongo/logus"
)

func DelayUpdateInviteClaimedCount(c context.Context, claimID int64) error {
	return delayerUpdateInviteClaimedCount.EnqueueWork(c, delaying.With(queues.QueueInvites, "UpdateInviteClaimedCount", 0), claimID)
}

func delayedUpdateInviteClaimedCount(c context.Context, claimID int64) (err error) {
	logus.Debugf(c, "delayerUpdateInviteClaimedCount(claimID=%v)", claimID)
	err = facade.RunReadwriteTransaction(c, func(tc context.Context, tx dal.ReadwriteTransaction) (err error) {
		claim := models4debtus.NewInviteClaim(claimID, nil)
		err = tx.Get(c, claim.Record)
		if err != nil {
			if dal.IsNotFound(err) {
				logus.Errorf(c, "Claim not found by id: %v", claimID)
				return nil
			}
			return fmt.Errorf("failed to get InviteClaimData by id=%v: %w", claimID, err)
		}
		invite, err := dtdal.Invite.GetInvite(c, tx, claim.Data.InviteCode)
		if err != nil {
			if dal.IsNotFound(err) {
				logus.Errorf(c, "Invite not found by code: %v", claim.Data.InviteCode)
				return nil // Internationally return NIL to avoid retrying
			}
			return err
		}
		for _, cid := range invite.Data.LastClaimIDs {
			if cid == claimID {
				logus.Infof(c, "Invite already has been updated for this claim (claimID=%v, inviteCode=%v).", claimID, claim.Data.InviteCode)
				return nil
			}
		}
		invite.Data.ClaimedCount += 1
		if invite.Data.LastClaimed.Before(claim.Data.DtClaimed) {
			invite.Data.LastClaimed = claim.Data.DtClaimed
		}
		invite.Data.LastClaimIDs = append(invite.Data.LastClaimIDs, claimID)
		if len(invite.Data.LastClaimIDs) > 10 {
			invite.Data.LastClaimIDs = invite.Data.LastClaimIDs[len(invite.Data.LastClaimIDs)-10:]
		}

		if err = tx.Set(tc, invite.Record); err != nil {
			return fmt.Errorf("failed to save invite to DB: %w", err)
		}
		return err
	})
	if err != nil {
		logus.Errorf(c, "Failed to update Invite.ClaimedCount for claimID=%v", claimID)
	}
	return err
}
