package models

import (
	"errors"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/general"
	"reflect"
	"time"
)

const (
	ReceiptKind = "Receipt"

	ReceiptStatusCreated      = "created"
	ReceiptStatusSending      = "sending"
	ReceiptStatusSent         = "sent"
	ReceiptStatusViewed       = "viewed"
	ReceiptStatusAcknowledged = "acknowledged"
)

var ReceiptStatuses = [4]string{
	ReceiptStatusCreated,
	ReceiptStatusSent,
	ReceiptStatusViewed,
	ReceiptStatusAcknowledged,
}

type Receipt = record.DataWithID[string, *ReceiptData]

func NewReceiptKey(id string) *dal.Key {
	if id == "" {
		return NewReceiptIncompleteKey()
	}
	return dal.NewKeyWithID(ReceiptKind, id)
}

func NewReceiptWithoutID(data *ReceiptData) Receipt {
	key := NewReceiptIncompleteKey()
	return record.NewDataWithID("", key, data)
}

func NewReceipt(id string, data *ReceiptData) Receipt {
	key := NewReceiptKey(id)
	return record.NewDataWithID(id, key, data)
}

const (
	ReceiptForFrom = "from"
	ReceiptForTo   = "to"
)

type ReceiptFor string

type ReceiptData struct {
	Status               string
	TransferID           string
	CreatorUserID        string     // IMPORTANT: Can be different from transfer.CreatorUserID (usually same). Think of 3d party bills
	For                  ReceiptFor `datastore:",noindex"` // TODO: always fill. If receipt.CreatorUserID != transfer.CreatorUserID then receipt.For must be set to either "from" or "to"
	ViewedByUserIDs      []string
	CounterpartyUserID   string // TODO: Is it always equal to AcknowledgedByUserID?
	AcknowledgedByUserID string // TODO: Is it always equal to CounterpartyUserID?
	general.CreatedOn
	TgInlineMsgID  string `datastore:",noindex"`
	DtCreated      time.Time
	DtSent         time.Time
	DtFailed       time.Time
	DtViewed       time.Time
	DtAcknowledged time.Time
	SentVia        string
	SentTo         string
	Lang           string `datastore:",noindex"`
	Error          string `datastore:",noindex"` //TODO: Need a comment on when it is used
}

func NewReceiptIncompleteKey() *dal.Key {
	return dal.NewIncompleteKey(ReceiptKind, reflect.Int, nil)
}

func NewReceiptEntity(creatorUserID, transferID, counterpartyUserID, lang, sentVia, sentTo string, createdOn general.CreatedOn) *ReceiptData {
	if creatorUserID == counterpartyUserID {
		panic("creatorUserID == counterpartyUserID")
	}
	if transferID == "" {
		panic("transferID == 0")
	}
	if createdOn.CreatedOnID == "" {
		panic("CreatedOnID is empty")
	}
	if createdOn.CreatedOnPlatform == "" {
		panic("CreatedOnPlatform is empty")
	}
	return &ReceiptData{
		CreatorUserID:      creatorUserID,
		CounterpartyUserID: counterpartyUserID,
		TransferID:         transferID,
		CreatedOn:          createdOn,
		DtCreated:          time.Now(),
		Lang:               lang,
		SentVia:            sentVia,
		SentTo:             sentTo,
		Status:             ReceiptStatusCreated,
	}
}

//func (r *ReceiptData) Load(ps []datastore.Property) error {
//	return datastore.LoadStruct(r, ps)
//}

func (r *ReceiptData) Validate() (err error) {
	if r.TransferID == "" {
		return errors.New("receipt.TransferID == 0")
	}
	if err = validateString("Unknown receipt.Status", r.Status, ReceiptStatuses[:]); err != nil {
		return err
	}
	if r.CreatorUserID == "" {
		err = errors.New("ReceiptData.CreatorUserID == 0")
		return
	}
	if r.CounterpartyUserID == r.CreatorUserID {
		err = errors.New("ReceiptData.CounterpartyUserID == ReceiptData.CreatorUserID")
		return
	}
	if r.CreatedOn.CreatedOnID == "" {
		err = errors.New("ReceiptData.CreatedOnID is empty")
		return
	}
	if r.CreatedOn.CreatedOnPlatform == "" {
		err = errors.New("ReceiptData.CreatedOnPlatform is empty")
		return
	}
	if r.Lang == "" {
		err = errors.New("ReceiptData.Lang is empty")
		return
	}
	if r.Status == "" {
		err = errors.New("ReceiptData.Status is empty")
		return
	}

	if r.DtCreated.IsZero() {
		r.DtCreated = time.Now()
	}

	//if properties, err = datastore.SaveStruct(r); err != nil {
	//	return
	//}
	//
	//if properties, err = gaedb.CleanProperties(properties, map[string]gaedb.IsOkToRemove{
	//	"TgInlineMsgID":        gaedb.IsEmptyString,
	//	"AcknowledgedByUserID": gaedb.IsZeroInt,
	//	"CounterpartyUserID":   gaedb.IsZeroInt,
	//	"DtAcknowledged":       gaedb.IsZeroTime,
	//	"DtFailed":             gaedb.IsZeroTime,
	//	"DtSent":               gaedb.IsZeroTime,
	//	"DtViewed":             gaedb.IsZeroTime,
	//	"Error":                gaedb.IsEmptyString,
	//	"For":                  gaedb.IsEmptyString,
	//	"SentTo":               gaedb.IsEmptyString,
	//	"SentVia":              gaedb.IsEmptyString,
	//}); err != nil {
	//	return
	//}

	return
}
