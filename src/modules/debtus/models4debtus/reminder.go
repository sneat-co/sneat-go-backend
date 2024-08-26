package models4debtus

import (
	"errors"
	"github.com/bots-go-framework/bots-fw-telegram"
	"github.com/dal-go/dalgo/record"
	"time"
)

const (
	ReminderStatusCreated           = "created"
	ReminderStatusSending           = "sending"
	ReminderStatusFailed            = "failed"
	ReminderStatusSent              = "sent"
	ReminderStatusViewed            = "viewed"
	ReminderStatusRescheduled       = "rescheduled"
	ReminderStatusUsed              = "used"
	ReminderStatusDiscarded         = "discarded"
	ReminderStatusInvalidNoTransfer = "invalid:no-transfer"
)

var ReminderStatuses = []string{
	ReminderStatusCreated,
	ReminderStatusSending,
	ReminderStatusFailed,
	ReminderStatusSent,
	ReminderStatusViewed,
	ReminderStatusRescheduled,
	ReminderStatusUsed,
	ReminderStatusDiscarded,
	ReminderStatusInvalidNoTransfer,
}

const ReminderKind = "Reminder"

//var _ datastore.PropertyLoadSaver = (*ReminderDbo)(nil)

type Reminder = record.DataWithID[string, *ReminderDbo]

//var _ db.EntityHolder = (*Reminder)(nil)

func NewReminder(id string, entity *ReminderDbo) Reminder {
	return record.NewDataWithID(id, nil, entity)
}

type ReminderDbo struct {
	ParentReminderID    int       `firestore:"parentReminderID,omitempty"`
	IsAutomatic         bool      `firestore:"isAutomatic,omitempty"`
	IsRescheduled       bool      `firestore:"isRescheduled,omitempty"`
	TransferID          string    `firestore:"transferID"`
	DtNext              time.Time `firestore:"dtNext"`
	DtScheduled         time.Time `firestore:"dtScheduled,omitempty"` // DtNext moves here once sent
	Locale              string    `firestore:"locale,omitempty"`
	ClosedByTransferIDs []string  `firestore:"closedByTransferIDs,omitempty"` // TODO: Why do we need list of IDs here?
	SentVia             string    `firestore:"sentVia,omitempty"`
	Status              string    `firestore:"status"`
	UserID              string    `firestore:"userID"`
	CounterpartyID      string    // If this field != 0 then r is to a counterparty
	DtCreated           time.Time `firestore:"dtCreated,omitempty"`
	DtUpdated           time.Time `firestore:"dtUpdated,omitempty"`
	DtSent              time.Time `firestore:"dtSent,omitempty"`
	DtUsed              time.Time `firestore:"dtUsed,omitempty"` // When a user clicks "Yes/no returned"
	DtViewed            time.Time `firestore:"dtViewed,omitempty"`
	DtDiscarded         time.Time `firestore:"dtDiscarded,omitempty"`
	BotID               string    `firestore:"botID,omitempty"`
	ChatIntID           int64     `firestore:"chatIntID,omitempty"`
	MessageIntID        int64     `firestore:"messageIntID,omitempty"`
	MessageStrID        string    `firestore:"messageStrID,omitempty"`
	ErrDetails          string    `firestore:"errDetails,omitempty"`
}

//func (r *ReminderDbo) Save() (properties []datastore.Property, err error) {
//	if err = r.validate(); err != nil {
//		return nil, err
//	}
//	r.DtUpdated = time.Now()
//	if properties, err = datastore.SaveStruct(r); err != nil {
//		return
//	}
//
//	if properties, err = gaedb.CleanProperties(properties, map[string]gaedb.IsOkToRemove{
//		"DtDiscarded":      gaedb.IsZeroTime,
//		"DtNext":           gaedb.IsZeroTime,
//		"DtScheduled":      gaedb.IsZeroTime,
//		"DtSent":           gaedb.IsZeroTime,
//		"DtUsed":           gaedb.IsZeroTime,
//		"DtViewed":         gaedb.IsZeroTime,
//		"ErrDetails":       gaedb.IsEmptyString,
//		"IsAutomatic":      gaedb.IsFalse,
//		"IsRescheduled":    gaedb.IsFalse,
//		"Locale":           gaedb.IsEmptyString,
//		"MessageIntID":     gaedb.IsZeroInt,
//		"MessageStrID":     gaedb.IsEmptyString,
//		"ParentReminderID": gaedb.IsZeroInt,
//		"SentVia":          gaedb.IsEmptyString,
//	}); err != nil {
//		return
//	}
//
//	return
//}

func (r *ReminderDbo) Validate() (err error) {
	if err = validateString("Unknown reminder.Status", r.Status, ReminderStatuses); err != nil {
		return err
	}
	if r.TransferID == "" {
		return errors.New("reminder.TransferID == 0")
	}
	if r.SentVia == "" {
		return errors.New("reminder.SentVia is empty")
	}
	if r.DtCreated.IsZero() {
		return errors.New("reminder.DtCreated.IsZero()")
	}
	if !r.DtSent.IsZero() && r.DtSent.Before(r.DtCreated) {
		return errors.New("reminder.DtSent.Before(n.DtCreated)")
	}
	if !r.DtViewed.IsZero() && r.DtViewed.Before(r.DtSent) {
		return errors.New("reminder.DtViewed.Before(n.DtSent)")
	}
	if r.ChatIntID != 0 && r.BotID == "" || r.ChatIntID == 0 && r.BotID != "" {
		return errors.New("r.TgChatID != 0 && r.TgBot == '' || r.TgChatID == 0 && r.TgBot != ''")
	}
	return nil
}

func NewReminderViaTelegram(botID string, chatID int64, userID, transferID string, isAutomatic bool, next time.Time) (reminder *ReminderDbo) {
	return &ReminderDbo{
		Status:      ReminderStatusCreated,
		SentVia:     telegram.PlatformID,
		BotID:       botID,
		ChatIntID:   chatID,
		UserID:      userID,
		TransferID:  transferID,
		DtCreated:   time.Now(),
		IsAutomatic: isAutomatic,
		DtNext:      next,
	}
}

func (r *ReminderDbo) ScheduleNextReminder(parentReminderID int, next time.Time) *ReminderDbo {
	reminder := *r
	reminder.ParentReminderID = parentReminderID
	reminder.Status = ReminderStatusRescheduled

	reminder.DtCreated = time.Now()
	reminder.DtNext = next
	reminder.Status = ReminderStatusCreated
	zero := time.Time{}
	reminder.DtSent = zero
	reminder.DtDiscarded = zero
	reminder.DtViewed = zero
	reminder.MessageStrID = ""
	reminder.MessageIntID = 0

	r.IsRescheduled = true
	return &reminder
}
