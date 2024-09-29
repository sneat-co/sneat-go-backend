package models4debtus

import (
	"context"
	"fmt"
	"github.com/crediterra/money"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	"github.com/sneat-co/sneat-core-modules/contactus/const4contactus"
	"github.com/sneat-co/sneat-core-modules/contactus/dto4contactus"
	"github.com/sneat-co/sneat-core-modules/core/coremodels"
	"github.com/sneat-co/sneat-core-modules/spaceus/dbo4spaceus"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/const4debtus"
	"github.com/strongo/strongoapp/with"
	"reflect"
	"strings"
	"time"
)

func NewDebtusContactDbo(details dto4contactus.ContactDetails) *DebtusSpaceContactDbo {
	return &DebtusSpaceContactDbo{
		Status: const4debtus.StatusActive,
		CreatedFields: with.CreatedFields{
			CreatedAtField: with.CreatedAtField{
				CreatedAt: time.Now(),
			},
		},
		ContactDetails: details,
	}
}

type DebtusSpaceContactEntry = record.DataWithID[string, *DebtusSpaceContactDbo]

func NewDebtusSpaceContactEntry(spaceID, contactID string, dbo *DebtusSpaceContactDbo) DebtusSpaceContactEntry {
	key := dbo4spaceus.NewSpaceModuleItemKey(spaceID, const4debtus.ModuleID, const4contactus.ContactsCollection, contactID)
	if dbo == nil {
		dbo = new(DebtusSpaceContactDbo)
	}
	return record.NewDataWithID(contactID, key, dbo)
}

func NewDebtusContactKey(spaceID, contactID string) *dal.Key {
	return dbo4spaceus.NewSpaceModuleItemKey(spaceID, const4debtus.ModuleID, const4contactus.ContactsCollection, contactID)
}

func DebtusContactRecords(contacts []DebtusSpaceContactEntry) (records []dal.Record) {
	records = make([]dal.Record, len(contacts))
	for i, contact := range contacts {
		records[i] = contact.Record
	}
	return
}

func NewDebtusSpaceContacts(spaceID string, contactIDs ...string) (contacts []DebtusSpaceContactEntry) {
	contacts = make([]DebtusSpaceContactEntry, len(contactIDs))
	for i, id := range contactIDs {
		if id == "" {
			panic(fmt.Sprintf("contactIDs[%d] == 0", i))
		}
		contacts[i] = NewDebtusSpaceContactEntry(spaceID, id, nil)
	}
	return
}

func NewDebtusContactRecord() dal.Record {
	return dal.NewRecordWithIncompleteKey(const4contactus.ContactsCollection, reflect.Int64, new(DebtusSpaceContactDbo))
}

func (dbo *DebtusSpaceContactDbo) MustMatchCounterparty(counterparty DebtusSpaceContactEntry) {
	panic("not implemented")
	//if !c.Data.Balance().Equal(counterparty.Data.Balance().Reversed()) {
	//	panic(fmt.Sprintf("contact[%s].Balance() != counterpartyContact[%s].Balance(): %v != %v", c.ContactID, counterparty.ContactID, c.Data.Balance(), counterparty.Data.Balance()))
	//}
	//if c.Data.BalanceCount != counterparty.Data.BalanceCount {
	//	panic(fmt.Sprintf("contact.BalanceCount != counterpartyContact.BalanceCount:  %v != %v", c.Data.BalanceCount, counterparty.Data.BalanceCount))
	//}
}

type WithCounterpartyFields struct {
	CounterpartyUserID    string `json:"counterpartyUserID,omitempty" firestore:"counterpartyUserID,omitempty"`       // The counterparty user UserID if registered
	CounterpartySpaceID   string `json:"counterpartySpaceID,omitempty" firestore:"counterpartyUserID,omitempty"`      // The counterparty user SpaceRef if registered
	CounterpartyContactID string `json:"counterpartyContactID,omitempty" firestore:"counterpartyContactID,omitempty"` // The counterparty user ContactID if registered
}

func (v *WithCounterpartyFields) Validate() error {
	return nil
}

// DebtusSpaceContactDbo is stored in a collection at path "/teams/{teamID}/modules/debtusbot/contacts/{contactID}".
type DebtusSpaceContactDbo struct {
	with.CreatedFields
	money.Balanced
	WithCounterpartyFields
	LinkedBy string `firestore:",omitempty"`
	//
	Status DebtusContactStatus
	dto4contactus.ContactDetails
	Transfers *UserContactTransfersInfo `firestore:"transfers,omitempty"`
	coremodels.SmsStats
	//
	//TelegramChatID int

	// Decided against as we do not need it really and would require either 2 Put() instead of 1 PutMulti()
	//LastTransferID int  `firestore:",omitempty"`

	SearchName          []string `firestore:"searchName,omitempty"` // Deprecated
	NoTransferUpdatesBy []string `firestore:"noTransferUpdatesBy,omitempty"`
	SpaceIDs            []string `firestore:"spaceIDs,omitempty"`
}

func (dbo *DebtusSpaceContactDbo) String() string {
	return fmt.Sprintf("DebtusSpaceContactEntry{CounterpartyUserID: %s, CounterpartyContactID: %s, Status: %s, ContactDetails: %v, LastTransferAt: %v}", dbo.CounterpartyUserID, dbo.CounterpartyContactID, dbo.Status, dbo.ContactDetails, dbo.LastTransferAt)
}

func (dbo *DebtusSpaceContactDbo) GetTransfersInfo() (transfersInfo *UserContactTransfersInfo) {
	return dbo.Transfers
}

func (dbo *DebtusSpaceContactDbo) SetTransfersInfo(transfersInfo UserContactTransfersInfo) error {
	if err := transfersInfo.Validate(); err != nil {
		return err
	}
	dbo.Transfers = &transfersInfo
	return nil
}

func (dbo *DebtusSpaceContactDbo) Info(counterpartyID string, note, comment string) TransferCounterpartyInfo {
	return TransferCounterpartyInfo{
		ContactID: counterpartyID,
		//UserID:      dbo.UserID,
		ContactName: dbo.FullName(),
		Note:        note,
		Comment:     comment,
	}
}

//func (entity *DebtusSpaceContactDbo) UpdateSearchName() {
//	fullName := entity.GetFullName()
//	entity.SearchName = []string{strings.ToLower(fullName)}
//	if entity.Username != "" {
//		username := strings.ToLower(fullName)
//		found := false
//		for _, searchName := range entity.SearchName {
//			if searchName == username {
//				found = true
//			}
//		}
//		if !found {
//			entity.SearchName = append(entity.SearchName, username)
//		}
//	}
//}

//func (entity *DebtusSpaceContactDbo) Load(ps []datastore.Property) error {
//	p2 := make([]datastore.Property, 0, len(ps))
//	for _, p := range ps {
//		switch p.Name {
//		case "SearchName": // Ignore legacy
//		default:
//			p2 = append(p2, p)
//		}
//	}
//	if err := datastore.LoadStruct(entity, p2); err != nil {
//		return err
//	}
//	if entity.PhoneNumberIsConfirmed && !entity.PhoneNumberConfirmed {
//		entity.PhoneNumberConfirmed = true
//	}
//	return nil
//}

//var contactPropertiesToClean = map[string]gaedb.IsOkToRemove{
//	// Remove obsolete
//	"PhoneNumberIsConfirmed": gaedb.IsObsolete,
//	"SearchName":             gaedb.IsObsolete,
//	// Remove defaults
//	"CounterpartyUserID":         gaedb.IsZeroInt,
//	"CounterpartyContactID": gaedb.IsZeroInt,
//	"BalanceCount":               gaedb.IsZeroInt,
//	"BalanceJson":                gaedb.IsEmptyStringOrSpecificValue("null"), //TODO: Remove once DB cleared
//	"SmsCount":                   gaedb.IsZeroInt,
//	"SmsCost":                    gaedb.IsZeroFloat,
//	"SmsCostUSD":                 gaedb.IsZeroInt,
//	"EmailAddress":               gaedb.IsEmptyString,
//	"EmailAddressOriginal":       gaedb.IsEmptyString,
//	"TransfersJson":              gaedb.IsEmptyJSON,
//	"Nickname":                   gaedb.IsEmptyString,
//	"FirstName":                  gaedb.IsEmptyString,
//	"LastName":                   gaedb.IsEmptyString,
//	"ScreenName":                 gaedb.IsEmptyString,
//	"PhoneNumber":                gaedb.IsZeroInt,
//	"PhoneNumberConfirmed":       gaedb.IsFalse,
//	"EmailConfirmed":             gaedb.IsFalse,
//	"TelegramUserID":             gaedb.IsZeroInt,
//}

// Validate returns error if not valid. TODO: Validate DebtusSpaceContactDbo.Balanced
func (dbo *DebtusSpaceContactDbo) Validate() (err error) {
	//dbo.UpdateSearchName()
	dbo.EmailAddressOriginal = strings.TrimSpace(dbo.EmailAddressOriginal)
	dbo.EmailAddress = strings.ToLower(dbo.EmailAddressOriginal)
	if err = dbo.CreatedFields.Validate(); err != nil {
		return
	}
	if err = dbo.WithCounterpartyFields.Validate(); err != nil {
		return err
	}
	return nil
}

//func (entity *DebtusSpaceContactDbo) Save() (properties []datastore.Property, err error) {
//	if err = entity.BeforeSave(); err != nil {
//		return
//	}
//
//	if properties, err = datastore.SaveStruct(entity); err != nil {
//		return
//	}
//
//	//if properties, err = gaedb.CleanProperties(properties, contactPropertiesToClean); err != nil {
//	//	return
//	//}
//
//	//checkHasProperties(DebtusContactsCollection, properties)
//
//	return
//}

func (dbo *DebtusSpaceContactDbo) BalanceWithInterest(_ context.Context, periodEnds time.Time) (balance money.Balance, err error) {
	if transferInfo := dbo.GetTransfersInfo(); transferInfo != nil {
		err = updateBalanceWithInterest(true, dbo.Balance, transferInfo.OutstandingWithInterest, periodEnds)
	}
	return
}

func ContactsByID(contacts []DebtusSpaceContactEntry) (contactsByID map[string]*DebtusSpaceContactDbo) {
	contactsByID = make(map[string]*DebtusSpaceContactDbo, len(contacts))
	for _, contact := range contacts {
		contactsByID[contact.ID] = contact.Data
	}
	return
}
