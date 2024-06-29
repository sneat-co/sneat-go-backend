package facade

import (
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/strongo/logus"
	"time"

	"context"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/models"
)

func CheckTransferCreatorNameAndFixIfNeeded(c context.Context, tx dal.ReadwriteTransaction, transfer models.TransferEntry) (err error) {
	if transfer.Data.Creator().UserName == "" {
		user, err := User.GetUserByID(c, tx, transfer.Data.CreatorUserID)
		if err != nil {
			return err
		}

		creatorFullName := user.Data.FullName()
		if creatorFullName == "" || creatorFullName == models.NoName {
			logus.Debugf(c, "Can't fix transfers creator name as user entity has no name defined.")
			return nil
		}

		logMessage := fmt.Sprintf("Fixing transfer(%s).Creator().UserName, created: %v", transfer.ID, transfer.Data.DtCreated)
		if transfer.Data.DtCreated.After(time.Date(2017, 8, 1, 0, 0, 0, 0, time.UTC)) {
			logus.Warningf(c, logMessage)
		} else {
			logus.Infof(c, logMessage)
		}

		if transfer, err = Transfers.GetTransferByID(c, tx, transfer.ID); err != nil {
			return err
		}
		if transfer.Data.Creator().UserName == "" {
			changed := false
			switch transfer.Data.Direction() {
			case models.TransferDirectionUser2Counterparty:
				transfer.Data.From().UserName = creatorFullName
				changed = true
			case models.TransferDirectionCounterparty2User:
				transfer.Data.To().UserName = creatorFullName
				changed = true
			}
			if changed {
				if err = Transfers.SaveTransfer(c, tx, transfer); err != nil {
					return err
				}
			}
		}
		return nil
	}
	return nil
}
