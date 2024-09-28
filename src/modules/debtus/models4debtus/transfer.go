package models4debtus

import (
	"errors"
	"fmt"
	"github.com/crediterra/money"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	"github.com/sneat-co/sneat-go-backend/src/core/coremodels"
	"github.com/sneat-co/sneat-go-backend/src/coremodules/auth/models4auth"
	"github.com/sneat-co/sneat-go-backend/src/coremodules/contactus/dto4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/gae_app/general"
	"github.com/strongo/decimal"
	"reflect"
	"time"
)

const MaxTransferAmount = decimal.Decimal64p2(^uint64(0) >> 8)

type TransferDirection string

func (d TransferDirection) Reverse() TransferDirection {
	switch d {
	case TransferDirectionUser2Counterparty:
		return TransferDirectionCounterparty2User
	case TransferDirectionCounterparty2User:
		return TransferDirectionUser2Counterparty
	default:
		panic("Reverse not supported for %v" + string(d))
	}
}

const ( // TransferEntry directions
	TransferDirectionUser2Counterparty = "u2c"
	TransferDirectionCounterparty2User = "c2u"
	TransferDirection3dParty           = "3d-party"
)

func IsKnownTransferDirection(direction TransferDirection) bool {
	switch direction {
	case TransferDirectionUser2Counterparty, TransferDirectionCounterparty2User, TransferDirection3dParty:
		return true
	}
	return false
}

const ( // TransferEntry statuses
	// TransferViewed   = "viewed" // TODO: use the status

	// TransferAccepted for api4transfers that have been accepted by the counterparty
	TransferAccepted = "accepted"

	// TransferDeclined for api4transfers that have been declined by the counterparty
	TransferDeclined = "declined"
)

const TransfersCollection = "transfers"

//var _ datastore.PropertyLoadSaver = (*TransferData)(nil)

type TransferEntry = record.DataWithID[string, *TransferData]

func NewTransfers(transferIDs []string) []TransferEntry {
	transfers := make([]TransferEntry, len(transferIDs))
	for i, transferID := range transferIDs {
		transfers[i] = NewTransfer(transferID, nil)
	}
	return transfers
}

func TransferFromRecord(r dal.Record) (transfer TransferEntry) {
	return TransferEntry{
		WithID: record.NewWithID(r.Key().ID.(string), r.Key(), r.Data),
		Data:   r.Data().(*TransferData),
	}
}

func TransfersFromRecords(records []dal.Record) (transfers []TransferEntry) {
	transfers = make([]TransferEntry, len(records))
	for i, r := range records {
		transfers[i] = TransferFromRecord(r)
	}
	return
}

func TransferRecords(transfers []TransferEntry) []dal.Record {
	records := make([]dal.Record, len(transfers))
	for i, transfer := range transfers {
		records[i] = transfer.Record
	}
	return records
}

func NewTransferKey(id string) *dal.Key {
	if id == "" {
		panic("id == 0")
	}
	return dal.NewKeyWithID(TransfersCollection, id)
}

var NewTransferRecord = func() dal.Record {
	return NewTransferWithIncompleteKey(nil).Record
}

func NewTransferWithIncompleteKey(data *TransferData) TransferEntry {
	key := dal.NewIncompleteKey(TransfersCollection, reflect.String, nil)
	if data == nil {
		data = new(TransferData)
	}
	return TransferEntry{
		WithID: record.NewWithID("", key, data),
		Data:   data,
	}
}

func NewTransfer(id string, data *TransferData) TransferEntry {
	key := NewTransferKey(id)
	if data == nil {
		data = new(TransferData)
	}
	return TransferEntry{
		WithID: record.WithID[string]{
			ID:     id,
			Record: dal.NewRecordWithData(key, data),
		},
		Data: data,
	}
}

//var _ db.EntityHolder = (*TransferEntry)(nil)

//func (TransferEntry) Kind() string {
//	return TransfersCollection
//}

//func (t TransferEntry) IntID() int64 {
//	return t.ContactID
//}

//func (t *TransferEntry) Entity() interface{} {
//	return t.TransferData
//}

//func (TransferEntry) NewEntity() interface{} {
//	return new(TransferData)
//}

//func (t *TransferEntry) SetEntity(entity interface{}) {
//	if entity == nil {
//		t.Data = nil
//	} else {
//		t.Data = entity.(*TransferData)
//	}
//}

func (t *TransferData) HasObsoleteProps() bool {
	return t.hasObsoleteProps
}

func (t *TransferData) GetStartDate() time.Time {
	return t.DtCreated // TODO: Change to DtStart?
}

func (t *TransferData) GetLendingValue() decimal.Decimal64p2 {
	return t.AmountInCents
}

type TransferData struct {
	hasObsoleteProps bool
	general.CreatedOn
	from *TransferCounterpartyInfo
	to   *TransferCounterpartyInfo

	BillIDs []string

	coremodels.SmsStats
	// DirectionObsoleteProp string `firestore:"direction,omitempty"`

	// We need it as it is not always possible to identify original transfer (think multiple & partial api4transfers)
	IsReturn bool `firestore:",omitempty"`

	// List of transfer to which this debt is a return. Should be populated only if IsReturn=True
	ReturnToTransferIDs []string `firestore:",omitempty"` // TODO: to make it obsolete - move to ReturnsJson
	//
	returns      TransferReturns // Deserialized cache
	ReturnsJson  string          `firestore:",omitempty"`
	ReturnsCount int             `firestore:",omitempty"`
	// ReturntransferIDs []int `firestore:",omitempty"` // Obsolete - replaced with ReturnsJson List of api4transfers that return money to this debts
	//
	CreatorUserID           string `firestore:",omitempty"` // Do not delete, is NOT obsolete!
	CreatorCounterpartyID   int    `firestore:",omitempty"` // TODO: Replace with <From|To>ContactID
	CreatorCounterpartyName string `firestore:",omitempty"` // TODO: Replace with <From|To>ContactName
	CreatorNote             string `firestore:",omitempty"` // TODO: Replace with <From|To>Note
	CreatorComment          string `firestore:",omitempty"` // TODO: Replace with <From|To>Comment

	CreatorTgReceiptByTgMsgID int64 `firestore:",omitempty"` // TODO: Move to ReceiptEntry ?
	//
	// CreatorTgBotID       string `firestore:",omitempty"` // TODO: Migrated to TransferCounterpartyInfo
	// CreatorTgChatID      int64  `firestore:",omitempty"` // TODO: Migrated to TransferCounterpartyInfo
	// CounterpartyTgBotID  string `firestore:",omitempty"` // TODO: Migrated to TransferCounterpartyInfo
	// CounterpartyTgChatID int64  `firestore:",omitempty"` // TODO: Migrated to TransferCounterpartyInfo
	//
	// CreatorAutoRemindersDisabled bool   `firestore:",omitempty"`
	// CreatorReminderID      int64 `firestore:",omitempty"` // obsolete
	// CounterpartyReminderID int64 `firestore:",omitempty"` // obsolete
	//
	//CounterpartyUserID           int64  `firestore:",omitempty"` // Replaced with <From|To>UserID
	//CounterpartyContactID   int64  `firestore:",omitempty"` // Replaced with <From|To>ContactID
	//CounterpartyCounterpartyName string `firestore:",omitempty"` // Replaced with <From|To>ContactName
	//CounterpartyNote             string `firestore:",omitempty"` // Replaced with <From|To>Note
	//CounterpartyComment          string `firestore:",omitempty"` // Replaced with <From|To>Comment
	// CounterpartyAutoRemindersDisabled bool   `firestore:",omitempty"`
	// CounterpartyTgReceiptInlineMessageID string    `firestore:",omitempty"` - not useful as we can edit message just once on callback

	FromJson string `firestore:"C_From,omitempty"`
	ToJson   string `firestore:"C_To,omitempty"`

	// ** New properties to replace Creator/DebtusSpaceContactEntry set of props **
	// FromUserID           int64  `firestore:",omitempty"`
	// FromUserName         string `firestore:",omitempty"`
	// FromCounterpartyID   int64  `firestore:",omitempty"`
	// FromCounterpartyName string `firestore:",omitempty"`
	// FromComment          string `firestore:",omitempty"`
	// FromNote             string `firestore:",omitempty"`
	// ToUserID             int64  `firestore:",omitempty"`
	// ToUserName           string `firestore:",omitempty"`
	// ToCounterpartyID     int64  `firestore:",omitempty"`
	// ToCounterpartyName   string `firestore:",omitempty"`
	// ToComment            string `firestore:",omitempty"`
	// ToNote               string `firestore:",omitempty"`

	AcknowledgeStatus string    `firestore:",omitempty"`
	AcknowledgeTime   time.Time `firestore:",omitempty"`

	// This 2 fields are used in conjunction with .Order("-DtCreated")
	BothUserIDs         []string // This is necessary to show transactions by user regardless who created
	BothCounterpartyIDs []string // This is necessary to show transactions by counterparty regardless who created
	//
	DtCreated time.Time
	DtDueOn   time.Time `firestore:"dtDueOn,omitempty"`

	AmountInCents         decimal.Decimal64p2 `firestore:"amountInCents"`
	AmountInCentsReturned decimal.Decimal64p2 `firestore:"amountInCentsReturned,omitempty"`
	AmountInCentsInterest decimal.Decimal64p2 `firestore:"amountInCentsInterest,omitempty"`
	// AmountInCentsOutstanding decimal.Decimal64p2 `firestore:",omitempty"` // Removed

	TransferInterest

	IsOutstanding bool               `firestore:"isOutstanding,omitempty"`
	Currency      money.CurrencyCode `firestore:"currency"` // Should be indexed for loading outstanding api4transfers
	//
	ReceiptsSentCount int      `firestore:"receiptsSentCount,omitempty"`
	ReceiptIDs        []string `firestore:"receiptIDs,omitempty"`
}

// AmountReturned returns amount returned to counterparty
func (t *TransferData) AmountReturned() decimal.Decimal64p2 {
	if t.AmountInCentsReturned > 0 {
		return t.AmountInCentsReturned
	}
	if t.IsReturn && t.AmountInCentsReturned == 0 {
		return t.AmountInCents
	}
	return 0
}

func (t *TransferData) String() string {
	return fmt.Sprintf(
		"TransferData{DtCreated: %v, Direction: %v, GetAmount(): %v, AmoutInCentsReturned: %v, IsReturn: %v, ReturnToTransferIDs: %v, CreatorUserID: %s, Creator: %v, DebtusSpaceContactEntry: %v, BothUserIDs: %v, BothCounterpartyIDs: %v, From: %v, To: %v}",
		t.DtCreated, t.Direction(), t.GetAmount(), t.AmountInCentsReturned, t.IsReturn, t.ReturnToTransferIDs, t.CreatorUserID, t.Creator(), t.Counterparty(), t.BothUserIDs, t.BothCounterpartyIDs, t.From(), t.To())
}

func (t *TransferData) Direction() TransferDirection {
	// if t.DirectionObsoleteProp != "" {
	// 	return TransferDirection(t.DirectionObsoleteProp)
	// }
	switch t.CreatorUserID {
	case "":
		panic("CreatorUserID == 0")
	case t.From().UserID:
		return TransferDirectionUser2Counterparty
	case t.To().UserID:
		return TransferDirectionCounterparty2User
	}
	return TransferDirection3dParty
}

func (t *TransferData) DirectionForUser(userID string) TransferDirection {
	switch userID {
	case t.From().UserID:
		return TransferDirectionUser2Counterparty
	case t.To().UserID:
		return TransferDirectionCounterparty2User
	case t.CreatorUserID:
		return TransferDirection3dParty
	default:
		panic(t.transferIsNotAssociatedWithUser(userID))
	}
}

func (t *TransferData) IsReverseDirection(t2 *TransferData) bool {
	return t.DirectionForUser(t.CreatorUserID) == t2.DirectionForUser(t.CreatorUserID).Reverse()
}

// DirectionForContact
func (t *TransferData) DirectionForContact(contactID string) TransferDirection {
	switch contactID {
	case t.From().ContactID:
		return TransferDirectionCounterparty2User
	case t.To().ContactID:
		return TransferDirectionUser2Counterparty
	default:
		panic(t.transferIsNotAssociatedWithContact(contactID))
	}
}

func (t *TransferData) transferIsNotAssociatedWithUser(userID string) string {
	return fmt.Sprintf(
		"TransferEntry is not associated with userID=%s  (FromUserID=%s, ToUserID=%s)",
		userID, t.From().UserID, t.To().UserID,
	)
}

func (t *TransferData) transferIsNotAssociatedWithContact(contactID string) string {
	return fmt.Sprintf(
		"TransferEntry is not associated with contactID=%s  (FromContactID=%s, ToContactID=%s)",
		contactID, t.From().ContactID, t.To().ContactID,
	)
}

func (t *TransferData) transferIsNotRelatedToCreator() string {
	return ErrTransferNotRelatedToCreator.Error() + fmt.Sprintf(
		"\nDirection(): %v, CreatorUserID: %s, From: %v, To: %v",
		t.Direction(), t.CreatorUserID, t.FromJson, t.ToJson,
	)
}

func (t *TransferData) ReturnDirectionForUser(userID string) TransferDirection {
	switch userID {
	case "":
		panic("userID == 0")
	case t.From().UserID:
		return TransferDirectionCounterparty2User
	case t.To().UserID:
		return TransferDirectionUser2Counterparty
	default:
		panic(t.transferIsNotAssociatedWithUser(userID))
	}
}

var ErrTransferNotRelatedToCreator = errors.New("TransferEntry is not related to creator")

func (t *TransferData) Creator() *TransferCounterpartyInfo { // TODO: Same as t.Creator()
	if t.CreatorUserID == "" {
		panic("CreatorUserID == 0")
	}
	if counterparty := t.From(); counterparty.UserID == t.CreatorUserID {
		return counterparty
	} else if counterparty = t.To(); counterparty.UserID == t.CreatorUserID {
		return counterparty
	}
	panic(t.transferIsNotRelatedToCreator())
}

func (t *TransferData) Counterparty() *TransferCounterpartyInfo {
	// return TransferCounterpartyInfo{
	// 	UserID:         t.CounterpartyUserID,
	// 	ContactID: t.CreatorCounterpartyID,
	// 	ContactName:           t.CreatorCounterpartyName,
	// 	Note:           t.CreatorNote,
	// 	Comment:        t.CreatorComment,
	// }
	switch t.Direction() {
	case TransferDirectionUser2Counterparty:
		return t.To()
	case TransferDirectionCounterparty2User:
		return t.From()
	default:
		panic(t.transferIsNotRelatedToCreator())
	}
}

func (t *TransferData) CounterpartyInfoByUserID(userID string) *TransferCounterpartyInfo {
	switch userID {
	case t.From().UserID:
		return t.To()
	case t.To().UserID:
		return t.From()
	default:
		panic(t.transferIsNotAssociatedWithUser(userID))
	}
}

func (t *TransferData) UserInfoByUserID(userID string) *TransferCounterpartyInfo {
	switch userID {
	case t.From().UserID:
		return t.from
	case t.To().UserID:
		return t.to
	default:
		panic(t.transferIsNotAssociatedWithUser(userID))
	}
}

// const TRANSFER_REMINDERS_DISABLED = "disabled"
//
// func (t *TransferEntry) IsRemindersDisabled(userID int64) bool {
// 	switch userID {
// 	case t.CreatorUserID:
// 		return t.CreatorAutoRemindersDisabled
// 	case t.CounterpartyUserID:
// 		return t.CounterpartyAutoRemindersDisabled
// 	default:
// 		panic("Attempt to check reminders for a user not related to the transfer")
// 	}
// }
//
// // Returns true if value have been changed and false if unchanged.
// func (t *TransferEntry) setAutoRemindersDisabled(userID int64, value bool) bool {
// 	switch userID {
// 	case t.CreatorUserID:
// 		if t.CreatorAutoRemindersDisabled != value {
// 			t.CreatorAutoRemindersDisabled = value
// 			return true
// 		}
// 	case t.CounterpartyUserID:
// 		if t.CounterpartyAutoRemindersDisabled != value {
// 			t.CounterpartyAutoRemindersDisabled = value
// 			return true
// 		}
// 	default:
// 		panic("Attempt to set remindersDisabled for a user not related to the transfer")
// 	}
// 	return false
// }
//
// // Returns true if value have been changed and false if unchanged.
// func (t *TransferEntry) EnableAutoReminders(userID int64) bool {
// 	return t.setAutoRemindersDisabled(userID, false)
// }
//
// // Returns true if value have been changed and false if unchanged.
// func (t *TransferEntry) DisableAutoReminders(userID int64) bool {
// 	return t.setAutoRemindersDisabled(userID, true)
// }

//func (t *TransferData) Load(ps []datastore.Property) error {
//	// Load I and J as usual.
//	p2 := make([]datastore.Property, 0, len(ps))
//	var creationPlatform string
//	var ( // TODO: obsolete props migrated to TransferCounterpartyJson
//		creatorReminderID, counterpartyReminderID         int
//		creatorTgChatID, counterpartyTgChatID             int64
//		creatorTgBotID, counterpartyTgBotID               string
//		creatorContactName, counterpartyContactName       string
//		creatorNote, counterpartyNote                     string
//		creatorComment, counterpartyComment               string
//		creatorUserID, counterpartyUserID                 int64
//		creatorCounterpartyID, counterpartyCounterpartyID int64
//		// creatorTgReceiptByTgMsgID, counterpartyTgReceiptByTgMsgID int64
//	)
//	for _, p := range ps {
//		switch p.Name {
//		// case "AmountInCentsOutstanding": // Ignore legacy
//		// 	t.hasObsoleteProps = true
//		case "CounterpartyAutoRemindersDisabled": // Ignore legacy
//			t.hasObsoleteProps = true
//		case "CreatorAutoRemindersDisabled": // Ignore legacy
//			t.hasObsoleteProps = true
//		case "ReturnTransferIDs": // Ignore legacy
//			t.hasObsoleteProps = true
//		case "IsDue2Notify": // Ignore legacy
//			t.hasObsoleteProps = true
//		case "DtDueNext": // Ignore legacy
//			t.hasObsoleteProps = true
//		case "CounterpartyNotifications": // Ignore legacy
//			t.hasObsoleteProps = true
//		case "CreatorNotifications": // Ignore legacy
//			t.hasObsoleteProps = true
//		case "CounterpartyTgReceiptInlineMessageID": // Ignore legacy
//			t.hasObsoleteProps = true
//		case "CreationPlatform":
//			t.hasObsoleteProps = true
//			creationPlatform = p.Value.(string)
//
//		case "CreatorCounterpartyName":
//			t.hasObsoleteProps = true
//			counterpartyContactName = p.Value.(string)
//		case "CounterpartyCounterpartyName":
//			t.hasObsoleteProps = true
//			counterpartyContactName = p.Value.(string)
//
//		case "CreatorNote":
//			t.hasObsoleteProps = true
//			creatorNote = p.Value.(string)
//		case "CounterpartyNote":
//			t.hasObsoleteProps = true
//			counterpartyNote = p.Value.(string)
//
//		case "CreatorComment":
//			t.hasObsoleteProps = true
//			creatorComment = p.Value.(string)
//		case "CounterpartyComment":
//			t.hasObsoleteProps = true
//			counterpartyComment = p.Value.(string)
//
//		case "CreatorUserID": // Is NOT obsolete!
//			// t.hasObsoleteProps = true
//			p2 = append(p2, p)
//			creatorUserID = p.Value.(int64)
//		case "CounterpartyUserID":
//			t.hasObsoleteProps = true
//			counterpartyUserID = p.Value.(int64)
//
//		case "CreatorCounterpartyID":
//			t.hasObsoleteProps = true
//			creatorCounterpartyID = p.Value.(int64)
//		case "CounterpartyContactID":
//			t.hasObsoleteProps = true
//			counterpartyCounterpartyID = p.Value.(int64)
//
//			// case "FromUserID": // TODO: Ignore legacy, temporary
//			// case "FromUserName": // TODO: Ignore legacy, temporary
//			// case "FromCounterpartyID": // TODO: Ignore legacy, temporary
//			// case "FromCounterpartyName": // TODO: Ignore legacy, temporary
//			// case "FromComment": // TODO: Ignore legacy, temporary
//			// case "FromNote": // TODO: Ignore legacy, temporary
//			// case "ToUserID": // TODO: Ignore legacy, temporary
//			// case "ToUserName": // TODO: Ignore legacy, temporary
//			// case "ToCounterpartyID": // TODO: Ignore legacy, temporary
//			// case "ToCounterpartyName": // TODO:  Ignore legacy, temporary
//			// case "ToComment": // TODO: Ignore legacy, temporary
//			// case "ToNote": // TODO: Ignore legacy, temporary
//
//		case "CreatorReminderID":
//			t.hasObsoleteProps = true
//			creatorReminderID = p.Value.(int)
//		case "CounterpartyReminderID":
//			t.hasObsoleteProps = true
//			counterpartyReminderID = p.Value.(int)
//
//		case "CreatorTgBotID":
//			t.hasObsoleteProps = true
//			creatorTgBotID = p.Value.(string)
//		case "CounterpartyTgBotID":
//			t.hasObsoleteProps = true
//			counterpartyTgBotID = p.Value.(string)
//
//		case "CreatorTgChatID":
//			t.hasObsoleteProps = true
//			creatorTgChatID = p.Value.(int64)
//		case "CounterpartyTgChatID":
//			t.hasObsoleteProps = true
//			counterpartyTgChatID = p.Value.(int64)
//		case "Amount", "AmountReturned", "AmountOutstanding":
//			t.hasObsoleteProps = true
//			if v, isFloat := p.Value.(float64); isFloat {
//				if v != 0 {
//					p.Name = strings.Replace(p.Name, "Amount", "AmountInCents", 1)
//					if v < 0.01 { // Fix very small amounts
//						v = 0.01
//					}
//					val := int64(decimal.NewDecimal64p2FromFloat64(v))
//					if val < 0 && v > 0 { // value is too big so we are getting overflow
//						val = int64(MaxTransferAmount)
//					}
//					p.Value = val
//					p2 = append(p2, p)
//				}
//			} else {
//				return fmt.Errorf("obsolete property '%v' expected to be of type float64, got: %T=%v", p.Name, p.Value, p.Value)
//			}
//		default:
//			p2 = append(p2, p)
//		}
//	}
//
//	if err := datastore.LoadStruct(t, p2); err != nil {
//		return err
//	}
//
//	// t.hasObsoleteProps = t.hasObsoleteProps || t.DirectionObsoleteProp != ""
//
//	if t.CreatedOnPlatform == "" && creationPlatform != "" {
//		t.CreatedOnPlatform = creationPlatform
//	}
//
//	// switch t.DirectionObsoleteProp {
//	// case "from":
//	// 	t.DirectionObsoleteProp = TransferDirectionUser2Counterparty
//	// case "to":
//	// 	t.DirectionObsoleteProp = TransferDirectionCounterparty2User
//	// }
//
//	if t.GetOutstandingValue(time.Now()) > 0 && !t.IsOutstanding {
//		t.IsOutstanding = true
//	}
//
//	{ // TODO: Get rid once all api4transfers migrated - Moves properties to JSON
//		migrateToCounterpartyInfo := func(
//			counterparty *TransferCounterpartyInfo,
//			userID int64,
//			contactID int64,
//			reminderID int,
//			tgChatID int64,
//			tgBotID, contactName, note, comment string,
//		) {
//			if userID != 0 && counterparty.UserID == 0 {
//				counterparty.UserID = userID
//			}
//			if contactID != 0 && counterparty.ContactID == 0 {
//				counterparty.ContactID = contactID
//			}
//			if reminderID != 0 {
//				counterparty.ReminderID = reminderID
//			}
//			if tgChatID != 0 {
//				counterparty.TgChatID = tgChatID
//			}
//			if tgBotID != "" {
//				counterparty.TgBotID = tgBotID
//			}
//			if contactName != "" && counterparty.ContactName == "" {
//				counterparty.ContactName = contactName
//			}
//			if note != "" && counterparty.Note == "" {
//				counterparty.Note = note
//			}
//			if comment != "" && counterparty.Comment == "" {
//				counterparty.Comment = comment
//			}
//		}
//
//		if creatorUserID != 0 { // TODO: temporary workaround to fix migration
//			migrateToCounterpartyInfo(t.Creator(), creatorUserID, counterpartyCounterpartyID, creatorReminderID, creatorTgChatID, creatorTgBotID, creatorContactName, creatorNote, creatorComment)
//			migrateToCounterpartyInfo(t.Counterparty(), counterpartyUserID, creatorCounterpartyID, counterpartyReminderID, counterpartyTgChatID, counterpartyTgBotID, counterpartyContactName, counterpartyNote, counterpartyComment)
//		}
//	}
//
//	return nil
//}

//var transferPropertiesToClean = map[string]gaedb.IsOkToRemove{
//	// Remove obsolete properties
//	"Amount":            gaedb.IsObsolete,
//	"AmountTotal":       gaedb.IsObsolete,
//	"AmountReturned":    gaedb.IsObsolete,
//	"AmountOutstanding": gaedb.IsObsolete,
//	//
//
//	// Remove defaults
//	"SmsCount":          gaedb.IsZeroInt,
//	"SmsCost":           gaedb.IsZeroFloat,
//	"SmsCostUSD":        gaedb.IsZeroInt,
//	"ReceiptsSentCount": gaedb.IsZeroInt,
//	// "CreatorReminderID":         gaedb.IsZeroInt,
//	// "CounterpartyReminderID":    gaedb.IsZeroInt,
//	// "CreatorTgChatID":           gaedb.IsZeroInt,
//	// "CounterpartyTgChatID":      gaedb.IsZeroInt,
//	"CreatorTgReceiptByTgMsgID": gaedb.IsZeroInt,
//	// "CounterpartyTgBotID":       gaedb.IsEmptyString,
//	// "CreatorTgBotID":            gaedb.IsEmptyString,
//	"Direction":                gaedb.IsEmptyString,
//	"BillID":                   gaedb.IsEmptyString,
//	"AmountInCentsOutstanding": gaedb.IsZeroInt,
//	"AmountInCentsReturned":    gaedb.IsZeroInt,
//	"AcknowledgeStatus":        gaedb.IsEmptyString,
//	"AcknowledgeTime":          gaedb.IsZeroTime,
//	"DtDueOn":                  gaedb.IsZeroTime,
//	"IsOutstanding":            gaedb.IsFalse,
//	"IsReturn":                 gaedb.IsFalse,
//}

func (t *TransferData) Validate() (err error) {
	if t.CreatorUserID == "" {
		err = errors.New("*TransferData.CreatorUserID == 0")
		return
	}

	if t.AmountInCents == 0 { // Should be always presented
		err = errors.New("*TransferData.AmountInCents == 0")
		return
	}

	if t.AmountInCents > MaxTransferAmount {
		err = fmt.Errorf("*TransferData.AmountInCents is too big, expected to be less then %v, got %v", MaxTransferAmount, t.AmountInCents)
		return
	}

	if t.Currency == "" { // Should be always presented
		err = errors.New("*TransferData.Currency is empty string")
		return
	}

	if t.AmountInCentsReturned < 0 {
		err = fmt.Errorf("*TransferData.AmountInCentsReturned:%v < 0", t.AmountInCentsReturned)
		return
	}

	if err = t.validateTransferInterestAndReturns(); err != nil {
		return
	}

	if t.IsOutstanding {
		switch t.HasInterest() {
		case true:
			// Can we simply check for zero outstanding value?
			// What if there is complex interest rule that allocate interest after grace period?
			if t.GetOutstandingValue(time.Now()) == 0 {
				t.IsOutstanding = false
			}
		case false:
			if t.AmountInCents == t.AmountInCentsReturned {
				t.IsOutstanding = false
			}
		}
	}

	// t.onSaveMigrateUserProps()

	// switch t.Direction() { // TODO: Delete later!
	// case "":
	// 	if t.BillID == "" && t.From().UserID == 0 && t.To().UserID == 0 {
	// 		err = errors.New("t.Direction is empty string")
	// 		return
	// 	}
	// case TransferDirectionUser2Counterparty:
	// case TransferDirectionCounterparty2User:
	// default:
	// 	err = errors.New("Unknown direction: " + t.DirectionObsoleteProp)
	// 	return
	// }

	// if t.AmountInCentsOutstanding < 0 {
	// 	err = fmt.Errorf("*TransferData.AmountInCentsOutstanding:%v < 0", t.AmountInCentsOutstanding)
	// 	return
	// }

	// if t.AmountInCentsReturned > t.AmountInCents {
	// 	err = fmt.Errorf("*TransferData.AmountInCentsReturned:%v > AmountInCents:%v", t.AmountInCentsReturned, t.AmountInCents)
	// 	return
	// }

	// if t.AmountInCentsOutstanding > t.AmountInCents {
	// 	err = fmt.Errorf("*TransferData.AmountInCentsOutstanding:%v > AmountInCents:%v", t.AmountInCentsOutstanding, t.AmountInCents)
	// 	return
	// }
	//
	// if t.AmountInCentsReturned+t.AmountInCentsOutstanding > t.AmountInCents {
	// 	err = fmt.Errorf("*TransferData.AmountInCentsReturned:%v + AmountInCentsOutstanding:%v > AmountInCents:%v", t.AmountInCentsReturned, t.AmountInCentsOutstanding, t.AmountInCents)
	// 	return
	// }

	if t.IsReturn {
		return errors.New("not implemented: temporally disabled this on 11 May 2018 to fix migration mapreduce")
		// TODO: Temporally commented just this if on 11 May 2018 to fix migration mapreduce
		// if len(t.ReturnToTransferIDs) == 0 {
		// 	err = errors.New("*TransferData: IsReturn == true && len(ReturnToTransferIDs) == 0")
		// 	return
		// }

		// if (t.AmountInCentsReturned != 0 || t.AmountInCentsOutstanding != 0) && t.AmountInCents != t.AmountInCentsReturned+t.AmountInCentsOutstanding {
		// 	err = fmt.Errorf("*TransferData: IsReturn == true && AmountInCents != AmountInCentsReturned + AmountInCentsOutstanding: %v != %v + %v", t.AmountInCents, t.AmountInCentsReturned, t.AmountInCentsOutstanding)
		// 	return
		// }
		// } else {
		// 	if t.AmountInCents != t.AmountInCentsReturned+t.AmountInCentsOutstanding {
		// 		err = fmt.Errorf("*TransferData: IsReturn == false && AmountInCents != AmountInCentsReturned + AmountInCentsOutstanding: %v != %v + %v", t.AmountInCents, t.AmountInCentsReturned, t.AmountInCentsOutstanding)
		// 		return
		// 	}
	}

	if t.CreatorUserID <= "" { // Should be always presented
		err = fmt.Errorf("*TransferData.CreatorUserID:%s <= 0", t.CreatorUserID)
		return
	}

	from := t.From()
	to := t.To()
	if from.UserName == dto4contactus.NoName {
		from.UserName = ""
	}
	if to.UserName == dto4contactus.NoName {
		to.UserName = ""
	}

	if from.ContactID == "" && to.ContactID == "" {
		err = errors.New("from.ContactID == 0 && to.ContactID == 0")
		return
	} else { // Always store 2 values, even if 1 is zero, so we can query such records.
		t.BothCounterpartyIDs = []string{from.ContactID, to.ContactID}
	}

	if from.UserID == "" && to.UserID == "" {
		if len(t.BillIDs) == 0 {
			err = errors.New("t.BillIDs is empty && t.From().UserID == 0 && t.To().UserID == 0")
			return
		}
		t.BothUserIDs = []string{}
	} else { // Always store 2 values, even if 1 is zero, so we can query such records.
		t.BothUserIDs = []string{from.UserID, to.UserID}
	}

	if from.UserID != t.CreatorUserID && from.ContactName == "" && from.UserName == "" { // Should be always presented
		err = errors.New("either FromCounterpartyName or FromUserName should be presented")
		return
	}
	if to.UserID != t.CreatorUserID && to.ContactName == "" && to.UserName == "" { // Should be always presented
		err = errors.New("either ToCounterpartyName or ToUserName should be presented")
		return
	}

	if isFixed, s := models4auth.FixContactName(from.ContactName); isFixed {
		from.ContactName = s
	}

	if isFixed, s := models4auth.FixContactName(to.ContactName); isFixed {
		to.ContactName = s
	}

	if err = t.onSaveSerializeJson(); err != nil {
		return
	}

	if t.FromJson == "" {
		err = errors.New("FromJson is empty")
		return
	}

	if t.ToJson == "" {
		err = errors.New("ToJson is empty")
		return
	}

	return
}

//func (*TransferData) movedToJson(propName string) bool {
//	return propName == "CounterpartyUserID" || (strings.HasPrefix(propName, "Creator") || strings.HasPrefix(propName, "Counterparty")) && (strings.HasSuffix(propName, "CounterpartyID") ||
//		strings.HasSuffix(propName, "CounterpartyName") ||
//		strings.HasSuffix(propName, "Note") ||
//		strings.HasSuffix(propName, "Comment") ||
//		strings.HasSuffix(propName, "TgBotID") ||
//		strings.HasSuffix(propName, "TgChatID"))
//}

//func (t *TransferData) Save() (properties []datastore.Property, err error) {
//	if err = t.BeforeSave(); err != nil {
//		return
//	}
//
//	// Serialize from struct to list of properties
//	if properties, err = datastore.SaveStruct(t); err != nil {
//		return properties, err
//	}
//
//	// To optimize storage we filter out default values
//	//if properties, err = gaedb.CleanProperties(properties, transferPropertiesToClean); err != nil {
//	//	return
//	//}
//
//	// { // Obsolete properties that were moved to JSON also should be removed
//	// 	if migratedToJson := t.FromJson != "" && t.ToJson != ""; migratedToJson {
//	// 		if t.DirectionObsoleteProp != "" {
//	// 			t.DirectionObsoleteProp = ""
//	// 		}
//	// 		properties2 := make([]datastore.Property, 0, len(properties))
//	// 		for _, p := range properties {
//	// 			if t.movedToJson(p.Name) {
//	// 				continue
//	// 			}
//	// 			properties2 = append(properties2, p)
//	// 		}
//	// 		properties = properties2
//	// 	}
//	// }
//
//	// Make general application-wide checks and call hooks if any
//	//checkHasProperties(TransfersCollection, properties)
//
//	return
//}

func NewTransferData(creatorUserID string, isReturn bool, amount money.Amount, from *TransferCounterpartyInfo, to *TransferCounterpartyInfo) *TransferData {
	if creatorUserID == "" {
		panic("creatorUserID == 0")
	}
	if from == nil {
		panic("from == nil")
	}
	if to == nil {
		panic("to == nil")
	}
	if amount.Value == 0 {
		panic("amount.Value == 0")
	}
	if amount.Currency == "" {
		panic("amount.Currency is empty")
	}
	transfer := &TransferData{
		CreatorUserID: creatorUserID,
		IsReturn:      isReturn,
		//
		from: from,
		to:   to,

		DtCreated: time.Now(),
		//
		// DirectionObsoleteProp: string(direction),
		AmountInCents: amount.Value,
		Currency:      amount.Currency,
	}
	if !isReturn {
		// transfer.AmountInCentsOutstanding = amount.Value
		transfer.IsOutstanding = true
	}
	return transfer
}

func (t *TransferData) GetAmount() money.Amount {
	return money.Amount{Currency: t.Currency, Value: t.AmountInCents}
}

func (t *TransferData) GetReturnedAmount() money.Amount {
	return money.Amount{Currency: t.Currency, Value: t.AmountReturned()}
}

//func ReverseTransfers(t []TransferEntry) {
//	last := len(t) - 1
//	for i := 0; i < len(t)/2; i++ {
//		t[i], t[last-i] = t[last-i], t[i]
//	}
//}
