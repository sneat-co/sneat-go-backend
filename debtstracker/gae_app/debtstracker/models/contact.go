package models

import (
	"context"
	"fmt"
	"github.com/crediterra/money"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	"github.com/pquerna/ffjson/ffjson"
	"github.com/strongo/strongoapp/with"
	"reflect"
	"strings"
	"time"
)

func NewDebtusContactData(userID string, details ContactDetails) *DebtusContactData {
	return &DebtusContactData{
		Status: STATUS_ACTIVE,
		UserID: userID,
		CreatedFields: with.CreatedFields{
			CreatedAtField: with.CreatedAtField{
				CreatedAt: time.Now(),
			},
			CreatedByField: with.CreatedByField{
				CreatedBy: userID,
			},
		},
		ContactDetails: details,
	}
}

const DebtusContactsCollection = "contacts"

type Contact = record.DataWithID[string, *DebtusContactData]

func NewDebtusContactKey(contactID string) *dal.Key {
	if contactID == "" {
		panic("NewDebtusContactKey(): contactID is required parameter")
	}
	return dal.NewKeyWithID(DebtusContactsCollection, contactID)
}

func DebtusContactRecords(contacts []Contact) (records []dal.Record) {
	records = make([]dal.Record, len(contacts))
	for i, contact := range contacts {
		records[i] = contact.Record
	}
	return
}

func NewDebtusContacts(ids ...string) (contacts []Contact) {
	contacts = make([]Contact, len(ids))
	for i, id := range ids {
		if id == "" {
			panic(fmt.Sprintf("ids[%d] == 0", i))
		}
		contacts[i] = NewDebtusContact(id, nil)
	}
	return
}

func NewDebtusContact(id string, data *DebtusContactData) Contact {
	key := NewDebtusContactKey(id)
	if data == nil {
		data = new(DebtusContactData)
	}
	return Contact{
		WithID: record.NewWithID(id, key, data),
		Data:   data,
	}
}

func NewDebtusContactRecord() dal.Record {
	return dal.NewRecordWithIncompleteKey(DebtusContactsCollection, reflect.Int64, new(DebtusContactData))
}

//var _ db.EntityHolder = (*Contact)(nil)

//func (Contact) Kind() string {
//	return DebtusContactsCollection
//}

//func (c *Contact) Entity() interface{} {
//	return c.Data
//}

//func (Contact) NewEntity() interface{} {
//	return new(DebtusContactData)
//}

//func (c *Contact) SetEntity(entity interface{}) {
//	if entity == nil {
//		c.Data = nil
//	} else {
//		c.Data = entity.(*DebtusContactData)
//	}
//}

func (entity *DebtusContactData) MustMatchCounterparty(counterparty Contact) {
	panic("not implemented")
	//if !c.Data.Balance().Equal(counterparty.Data.Balance().Reversed()) {
	//	panic(fmt.Sprintf("contact[%s].Balance() != counterpartyContact[%s].Balance(): %v != %v", c.ID, counterparty.ID, c.Data.Balance(), counterparty.Data.Balance()))
	//}
	//if c.Data.BalanceCount != counterparty.Data.BalanceCount {
	//	panic(fmt.Sprintf("contact.BalanceCount != counterpartyContact.BalanceCount:  %v != %v", c.Data.BalanceCount, counterparty.Data.BalanceCount))
	//}
}

// DebtusContactData is stored in a collection at path "/teams/{teamID}/modules/debtus/contacts/{contactID}".
type DebtusContactData struct {
	with.CreatedFields
	money.Balanced
	UserID                     string // owner can not be in parent key as we have problem with filtering transfers then
	CounterpartyUserID         string // The counterparty user ID if registered
	CounterpartyCounterpartyID string
	LinkedBy                   string `datastore:",noindex"`
	//
	Status string
	ContactDetails
	TransfersJson string `datastore:",noindex"`
	SmsStats
	//
	//TelegramChatID int
	//
	//LasttransferID int  `datastore:",noindex"` - Decided against as we do not need it really and would require either 2 Put() instead of 1 PutMulti()
	SearchName          []string `datastore:",noindex"` // Deprecated
	NoTransferUpdatesBy []string `datastore:",noindex"`
	GroupIDs            []string `datastore:",noindex"`
}

func (entity *DebtusContactData) String() string {
	return fmt.Sprintf("Contact{UserID: %v, CounterpartyUserID: %v, CounterpartyCounterpartyID: %v, Status: %v, ContactDetails: %v, Balance: '%v', LastTransferAt: %v}", entity.UserID, entity.CounterpartyUserID, entity.CounterpartyCounterpartyID, entity.Status, entity.ContactDetails, entity.BalanceJson, entity.LastTransferAt)
}

func (entity *DebtusContactData) GetTransfersInfo() (transfersInfo *UserContactTransfersInfo) {
	if entity.TransfersJson == "" {
		return &UserContactTransfersInfo{}
	}
	transfersInfo = new(UserContactTransfersInfo)
	if err := ffjson.Unmarshal([]byte(entity.TransfersJson), transfersInfo); err != nil {
		panic(err)
	}
	return
}

func (entity *DebtusContactData) SetTransfersInfo(transfersInfo UserContactTransfersInfo) error {
	if data, err := ffjson.Marshal(transfersInfo); err != nil {
		return err
	} else {
		entity.TransfersJson = string(data)
		return nil
	}
}

func (entity *DebtusContactData) Info(counterpartyID string, note, comment string) TransferCounterpartyInfo {
	return TransferCounterpartyInfo{
		ContactID:   counterpartyID,
		UserID:      entity.UserID,
		ContactName: entity.FullName(),
		Note:        note,
		Comment:     comment,
	}
}

//func (entity *DebtusContactData) UpdateSearchName() {
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

//func (entity *DebtusContactData) Load(ps []datastore.Property) error {
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
//	"CounterpartyCounterpartyID": gaedb.IsZeroInt,
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

func (entity *DebtusContactData) Validate() (err error) {
	//entity.UpdateSearchName()
	entity.EmailAddressOriginal = strings.TrimSpace(entity.EmailAddressOriginal)
	entity.EmailAddress = strings.ToLower(entity.EmailAddressOriginal)
	return nil
}

//func (entity *DebtusContactData) Save() (properties []datastore.Property, err error) {
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

func (entity *DebtusContactData) BalanceWithInterest(c context.Context, periodEnds time.Time) (balance money.Balance, err error) {
	balance = entity.Balance()
	if transferInfo := entity.GetTransfersInfo(); transferInfo != nil {
		err = updateBalanceWithInterest(true, balance, transferInfo.OutstandingWithInterest, periodEnds)
	}
	return
}

func ContactsByID(contacts []Contact) (contactsByID map[string]*DebtusContactData) {
	contactsByID = make(map[string]*DebtusContactData, len(contacts))
	for _, contact := range contacts {
		contactsByID[contact.ID] = contact.Data
	}
	return
}
