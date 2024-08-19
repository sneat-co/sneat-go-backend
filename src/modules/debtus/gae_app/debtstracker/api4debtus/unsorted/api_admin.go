package unsorted

import (
	"context"
	"errors"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/auth"
	"github.com/sneat-co/sneat-go-backend/src/auth/facade4auth"
	"github.com/sneat-co/sneat-go-backend/src/core/queues"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/dal4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/facade4debtus"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/facade4debtus/dto4debtus"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/gae_app/debtstracker/api4debtus"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/models4debtus"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/dal4userus"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/dbo4userus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/delaying"
	"github.com/strongo/logus"
	"github.com/strongo/validation"
	"net/http"
	"reflect"
)

func HandleAdminFindUser(c context.Context, w http.ResponseWriter, r *http.Request, _ auth.AuthInfo) {

	if userID := r.URL.Query().Get("userID"); userID != "" {
		appUser := dbo4userus.NewUserEntry(userID)
		if err := dal4userus.GetUser(c, nil, appUser); err != nil {
			logus.Errorf(c, fmt.Errorf("failed to get user by ContactID=%v: %w", userID, err).Error())
		} else {
			api4debtus.JsonToResponse(c, w, []dto4debtus.ApiUserDto{{ID: userID, Name: appUser.Data.GetFullName()}})
		}
		return
	} else {
		tgUserText := r.URL.Query().Get("tgUser")

		if tgUserText == "" {
			api4debtus.BadRequestMessage(c, w, "tgUser is empty string")
			return
		}

		tgUsers, err := facade4auth.TgUser.FindByUserName(c, nil, tgUserText)

		if err != nil {
			api4debtus.InternalError(c, w, err)
			return
		}

		users := make([]dto4debtus.ApiUserDto, len(tgUsers))

		for i, tgUser := range tgUsers {
			users[i] = dto4debtus.ApiUserDto{
				ID:   tgUser.GetAppUserID(),
				Name: tgUser.BaseData().UserName,
			}
		}

		api4debtus.JsonToResponse(c, w, users)
	}
}

func HandleAdminMergeUserContacts(c context.Context, w http.ResponseWriter, r *http.Request, _ auth.AuthInfo) {
	keepID := api4debtus.GetStrID(c, w, r, "keepID")
	if keepID == "" {
		return
	}
	deleteID := api4debtus.GetStrID(c, w, r, "deleteID")
	if deleteID == "" {
		return
	}
	spaceID := api4debtus.GetStrID(c, w, r, "spaceID")
	if spaceID == "" {
		api4debtus.BadRequestError(c, w, validation.NewErrRequestIsMissingRequiredField("spaceID"))
		return
	}

	logus.Infof(c, "keepID: %s, deleteID: %s", keepID, deleteID)

	if err := facade.RunReadwriteTransaction(c, func(c context.Context, tx dal.ReadwriteTransaction) error {

		contacts, err := dal4contactus.GetContactsByIDs(c, tx, spaceID, []string{keepID, deleteID})
		if err != nil {
			return err
		}
		if len(contacts) < 2 {
			return fmt.Errorf("len(contacts):%d < 2", len(contacts))
		}
		contactToKeep := contacts[0]
		contactToDelete := contacts[1]
		if contactToKeep.Data.UserID != contactToDelete.Data.UserID {
			return errors.New("contactToKeep.UserID != contactToDelete.UserID")
		}
		if contactToDelete.Data.UserID != "" && contactToKeep.Data.UserID == "" {
			return errors.New("contactToDelete.CounterpartyUserID != 0 && contactToKeep.CounterpartyUserID == 0")
		}
		contactusSpace := dal4contactus.NewContactusSpaceEntry(spaceID)

		if err = dal4contactus.GetContactusSpace(c, tx, contactusSpace); err != nil {
			return err
		}

		if contactusSpace.Data.HasContact(deleteID) {
			update := contactusSpace.Data.RemoveContact(deleteID)
			if err = tx.Update(c, contactusSpace.Key, []dal.Update{update}); err != nil {
				return err
			}
		}
		if err := delayChangeTransfersCounterparty.EnqueueWork(c, delaying.With(queues.QueueSupport, "changeTransfersCounterparty", 0), deleteID, keepID, ""); err != nil {
			return err
		}
		if err := tx.Delete(c, models4debtus.NewDebtusContactKey(spaceID, deleteID)); err != nil {
			return err
		} else {
			logus.Warningf(c, "DebtusSpaceContactEntry %s has been deleted from DB (non revocable)", deleteID)
		}
		return nil
	}); err != nil {
		api4debtus.ErrorAsJson(c, w, http.StatusInternalServerError, err)
		return
	}
}

func DelayedChangeTransfersCounterparty(c context.Context, oldID, newID int64, cursor string) (err error) {
	logus.Debugf(c, "delayedChangeTransfersCounterparty(oldID=%d, newID=%d)", oldID, newID)

	var q = dal.From(models4debtus.TransfersCollection).
		WhereField("BothCounterpartyIDs", dal.Equal, oldID).
		Limit(100).
		SelectKeysOnly(reflect.Int)

	var db dal.DB

	if db, err = facade.GetDatabase(c); err != nil {
		return err
	}

	var reader dal.Reader
	if reader, err = db.QueryReader(c, q); err != nil {
		return err
	}
	transferIDs, err := dal.SelectAllIDs[int](reader, dal.WithLimit(q.Limit()))
	if err != nil {
		return err
	}

	logus.Infof(c, "Loaded %d transferIDs", len(transferIDs))
	args := make([][]interface{}, len(transferIDs))
	for i, id := range transferIDs {
		args[i] = []interface{}{id, oldID, newID, ""}
	}
	return delayChangeTransferCounterparty.EnqueueWorkMulti(c, delaying.With(queues.QueueSupport, "changeTransferCounterparty", 0), args...)
}

func DelayedChangeTransferCounterparty(c context.Context, spaceID, transferID, oldID, newID string, cursor string) (err error) {
	logus.Debugf(c, "delayedChangeTransferCounterparty(spaceID=%s, oldID=%s, newID=%s, cursor=%s)", spaceID, oldID, newID, cursor)
	if _, err = facade4debtus.GetDebtusSpaceContactByID(c, nil, spaceID, newID); err != nil {
		return err
	}
	err = facade.RunReadwriteTransaction(c, func(c context.Context, tx dal.ReadwriteTransaction) error {
		transfer, err := facade4debtus.Transfers.GetTransferByID(c, tx, transferID)
		if err != nil {
			return err
		}
		changed := false
		for i, contactID := range transfer.Data.BothCounterpartyIDs {
			if contactID == oldID {
				transfer.Data.BothCounterpartyIDs[i] = newID
				changed = true
				break
			}
		}
		if changed {
			if from := transfer.Data.From(); from.ContactID == oldID {
				from.ContactID = newID
			} else if to := transfer.Data.To(); to.ContactID == oldID {
				to.ContactID = newID
			}
			err = facade4debtus.Transfers.SaveTransfer(c, tx, transfer)
		}
		return err
	})
	return err
}