package gaedal

import (
	"context"
	"fmt"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/coremodules/contactus/const4contactus"
	"github.com/sneat-co/sneat-go-backend/src/coremodules/contactus/dal4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/const4debtus"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/gae_app/debtstracker/dtdal"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/models4debtus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/delaying"
	"github.com/strongo/logus"
	"strings"
)

type ContactDalGae struct {
}

func NewContactDalGae() ContactDalGae {
	return ContactDalGae{}
}

var _ dtdal.ContactDal = (*ContactDalGae)(nil)

func (contactDalGae ContactDalGae) DeleteContact(ctx context.Context, tx dal.ReadwriteTransaction, spaceID, contactID string) (err error) {
	logus.Debugf(ctx, "ContactDalGae.DeleteContact(spaceID=%s, contactID=%s)", spaceID, contactID)
	if err = tx.Delete(ctx, models4debtus.NewDebtusContactKey(spaceID, contactID)); err != nil {
		return
	}
	if err = delayDeleteContactTransfers(ctx, contactID, ""); err != nil { // TODO: Move to facade4debtus!
		return
	}
	return
}

const DeleteContactTransfersFuncKey = "DeleteContactTransfers"

func delayDeleteContactTransfers(ctx context.Context, contactID string, cursor string) error {
	if err := delayerDeleteContactTransfersDelayFunc.EnqueueWork(ctx, delaying.With(const4debtus.QueueTransfers, DeleteContactTransfersFuncKey, 0), contactID, cursor); err != nil {
		return err
	}
	return nil
}

func delayedDeleteContactTransfers(ctx context.Context, contactID string, cursor string) (err error) {
	logus.Debugf(ctx, "delayedDeleteContactTransfers(contactID=%s, cursor=%v", contactID, cursor)
	const limit = 100
	var transferIDs []string
	transferIDs, cursor, err = dtdal.Transfer.LoadTransferIDsByContactID(ctx, contactID, limit, cursor)
	if err != nil {
		return
	}
	keys := make([]*dal.Key, len(transferIDs))
	for i, transferID := range transferIDs {
		keys[i] = models4debtus.NewTransferKey(transferID)
	}
	if err = facade.RunReadwriteTransaction(ctx, func(ctx context.Context, tx dal.ReadwriteTransaction) (err error) {
		if err = tx.DeleteMulti(ctx, keys); err != nil {
			return err
		}
		if len(transferIDs) == limit {
			if err = delayDeleteContactTransfers(ctx, contactID, cursor); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return
	}
	return
}

func (ContactDalGae) SaveContact(ctx context.Context, tx dal.ReadwriteTransaction, contact models4debtus.DebtusSpaceContactEntry) error {
	if err := tx.Set(ctx, contact.Record); err != nil {
		return fmt.Errorf("failed to SaveContact(): %w", err)
	}
	return nil
}

func newUserActiveContactsQuery(userID string) dal.QueryBuilder {
	return newUserContactsQuery(userID).WhereField("Status", dal.Equal, const4debtus.StatusActive)
}

func newUserContactsQuery(userID string) dal.QueryBuilder {
	return dal.From(const4contactus.ContactsCollection).WhereField("UserID", dal.Equal, userID)
}

func (ContactDalGae) GetContactsWithDebts(ctx context.Context, tx dal.ReadSession, spaceID, userID string) (counterparties []models4debtus.DebtusSpaceContactEntry, err error) {
	query := newUserContactsQuery(userID).
		WhereField("BalanceCount", dal.GreaterThen, 0).
		SelectInto(models4debtus.NewDebtusContactRecord)
	//var (
	//	counterpartyEntities []*models.DebtusSpaceContactDbo
	//)
	records, err := tx.QueryAllRecords(ctx, query)
	counterparties = make([]models4debtus.DebtusSpaceContactEntry, len(records))
	for i, record := range records {
		counterparties[i] = models4debtus.NewDebtusSpaceContactEntry(spaceID, record.Key().ID.(string), record.Data().(*models4debtus.DebtusSpaceContactDbo))
	}
	return
}

func (ContactDalGae) GetLatestContacts(whc botsfw.WebhookContext, tx dal.ReadSession, spaceID string, limit, totalCount int) (contacts []models4debtus.DebtusSpaceContactEntry, err error) {
	ctx := whc.Context()
	appUserID := whc.AppUserID()
	query := newUserActiveContactsQuery(appUserID).
		OrderBy(dal.DescendingField("LastTransferAt")).
		Limit(limit).
		SelectInto(models4debtus.NewDebtusContactRecord)
	if tx == nil {
		if tx, err = facade.GetSneatDB(ctx); err != nil {
			return
		}
	}
	var records []dal.Record
	records, err = tx.QueryAllRecords(ctx, query)
	var contactsCount = len(records)
	logus.Debugf(ctx, "GetLatestContacts(limit=%v, totalCount=%v): %v", limit, totalCount, contactsCount)
	if (limit == 0 && contactsCount < totalCount) || (limit > 0 && totalCount > 0 && contactsCount < limit && contactsCount < totalCount) {
		logus.Debugf(ctx, "Querying contacts without index -LastTransferAt")
		query = newUserActiveContactsQuery(appUserID).
			Limit(limit).
			SelectInto(models4debtus.NewTransferRecord)
		if records, err = tx.QueryAllRecords(ctx, query); err != nil {
			return
		}
	}
	contacts = make([]models4debtus.DebtusSpaceContactEntry, len(records))
	for i, record := range records {
		contactID := record.Key().ID.(string)
		dbo := record.Data().(*models4debtus.DebtusSpaceContactDbo)
		contacts[i] = models4debtus.NewDebtusSpaceContactEntry(spaceID, contactID, dbo)
	}
	return
}

func (contactDalGae ContactDalGae) GetContactIDsByTitle(ctx context.Context, tx dal.ReadSession, spaceID, userID string, title string, caseSensitive bool) (contactIDs []string, err error) {
	contactusSpace := dal4contactus.NewContactusSpaceEntry(spaceID)
	if err = dal4contactus.GetContactusSpace(ctx, tx, contactusSpace); err != nil {
		return
	}
	if caseSensitive {
		for id, contact := range contactusSpace.Data.Contacts {
			if contact.Names.GetFullName() == title {
				contactIDs = append(contactIDs, id)
			}
		}
	} else {
		title = strings.ToLower(title)
		for id, contact := range contactusSpace.Data.Contacts {
			if strings.ToLower(contact.Names.GetFullName()) == title {
				contactIDs = append(contactIDs, id)
			}
		}
	}
	return
}

//func zipCounterparty(keys []*datastore.Key, entities []*models.DebtusSpaceContactDbo) (contacts []models.DebtusSpaceContactEntry) {
//	if len(keys) != len(entities) {
//		panic(fmt.Sprintf("len(keys):%d != len(entities):%d", len(keys), len(entities)))
//	}
//	contacts = make([]models.DebtusSpaceContactEntry, len(entities))
//	for i, entity := range entities {
//		contacts[i] = models.NewDebtusSpaceContactEntry(keys[i].IntID(), entity)
//	}
//	return
//}

func (contactDalGae ContactDalGae) InsertContact(ctx context.Context, tx dal.ReadwriteTransaction, contactEntity *models4debtus.DebtusSpaceContactDbo) (
	contact models4debtus.DebtusSpaceContactEntry, err error,
) {
	contact.Data = contactEntity
	if err = tx.Insert(ctx, contact.Record); err != nil {
		return
	}
	contact.ID = contact.Key.ID.(string)
	return
}
