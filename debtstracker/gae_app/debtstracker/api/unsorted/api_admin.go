package unsorted

import (
	"context"
	"errors"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/api"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/auth"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/common"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/dtdal"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/facade"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/facade/dto"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/models"
	"github.com/strongo/delaying"
	"github.com/strongo/log"
	"net/http"
	"reflect"
)

func HandleAdminFindUser(c context.Context, w http.ResponseWriter, r *http.Request, _ auth.AuthInfo) {

	if userID := r.URL.Query().Get("userID"); userID != "" {
		if user, err := dtdal.User.GetUserByStrID(c, userID); err != nil {
			log.Errorf(c, fmt.Errorf("failed to get user by ID=%v: %w", userID, err).Error())
		} else {
			api.JsonToResponse(c, w, []dto.ApiUserDto{{ID: userID, Name: user.Data.FullName()}})
		}
		return
	} else {
		tgUserText := r.URL.Query().Get("tgUser")

		if tgUserText == "" {
			api.BadRequestMessage(c, w, "tgUser is empty string")
			return
		}

		tgUsers, err := dtdal.TgUser.FindByUserName(c, nil, tgUserText)

		if err != nil {
			api.InternalError(c, w, err)
			return
		}

		users := make([]dto.ApiUserDto, len(tgUsers))

		for i, tgUser := range tgUsers {
			users[i] = dto.ApiUserDto{
				ID:   tgUser.GetAppUserID(),
				Name: tgUser.BaseData().UserName,
			}
		}

		api.JsonToResponse(c, w, users)
	}
}

func HandleAdminMergeUserContacts(c context.Context, w http.ResponseWriter, r *http.Request, _ auth.AuthInfo) {
	keepID := api.GetStrID(c, w, r, "keepID")
	if keepID == "" {
		return
	}
	deleteID := api.GetStrID(c, w, r, "deleteID")
	if deleteID == "" {
		return
	}

	log.Infof(c, "keepID: %d, deleteID: %d", keepID, deleteID)

	db, err := facade.GetDatabase(c)
	if err != nil {
		api.ErrorAsJson(c, w, http.StatusInternalServerError, err)
		return
	}

	if err := db.RunReadwriteTransaction(c, func(c context.Context, tx dal.ReadwriteTransaction) error {
		contacts, err := facade.GetContactsByIDs(c, tx, []string{keepID, deleteID})
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
		if contactToDelete.Data.CounterpartyUserID != "" && contactToKeep.Data.CounterpartyUserID == "" {
			return errors.New("contactToDelete.CounterpartyUserID != 0 && contactToKeep.CounterpartyUserID == 0")
		}
		user, err := facade.User.GetUserByID(c, tx, contactToKeep.Data.UserID)
		if err != nil {
			return err
		}
		if user.ID != "" {
			return errors.New("not implemented yet: Need to update counterparty & user balances + last transfer info")
		}
		if userChanged := user.Data.RemoveContact(deleteID); userChanged {
			if err = facade.User.SaveUser(c, tx, user); err != nil {
				return err
			}
		}
		if err := delayChangeTransfersCounterparty.EnqueueWork(c, delaying.With(common.QUEUE_SUPPORT, "changeTransfersCounterparty", 0), deleteID, keepID, ""); err != nil {
			return err
		}
		if err := tx.Delete(c, models.NewDebtusContactKey(deleteID)); err != nil {
			return err
		} else {
			log.Warningf(c, "Contact %d has been deleted from DB (non revocable)", deleteID)
		}
		return nil
	}); err != nil {
		api.ErrorAsJson(c, w, http.StatusInternalServerError, err)
		return
	}
}

func DelayedChangeTransfersCounterparty(c context.Context, oldID, newID int64, cursor string) (err error) {
	log.Debugf(c, "delayedChangeTransfersCounterparty(oldID=%d, newID=%d)", oldID, newID)

	var q = dal.From(models.TransferKind).
		WhereField("BothCounterpartyIDs", dal.Equal, oldID).
		Limit(100).
		SelectKeysOnly(reflect.Int)

	var reader dal.Reader
	if reader, err = facade.DB().QueryReader(c, q); err != nil {
		return err
	}
	transferIDs, err := dal.SelectAllIDs[int](reader, dal.WithLimit(q.Limit()))
	if err != nil {
		return err
	}

	log.Infof(c, "Loaded %d transferIDs", len(transferIDs))
	args := make([][]interface{}, len(transferIDs))
	for i, id := range transferIDs {
		args[i] = []interface{}{id, oldID, newID, ""}
	}
	return delayChangeTransferCounterparty.EnqueueWorkMulti(c, delaying.With(common.QUEUE_SUPPORT, "changeTransferCounterparty", 0), args...)
}

func DelayedChangeTransferCounterparty(c context.Context, transferID string, oldID, newID string, cursor string) (err error) {
	log.Debugf(c, "delayedChangeTransferCounterparty(oldID=%d, newID=%d, cursor=%v)", oldID, newID, cursor)
	if _, err = facade.GetContactByID(c, nil, newID); err != nil {
		return err
	}
	var db dal.DB
	if db, err = facade.GetDatabase(c); err != nil {
		return err
	}
	err = db.RunReadwriteTransaction(c, func(c context.Context, tx dal.ReadwriteTransaction) error {
		transfer, err := facade.Transfers.GetTransferByID(c, tx, transferID)
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
			err = facade.Transfers.SaveTransfer(c, tx, transfer)
		}
		return err
	})
	return err
}
