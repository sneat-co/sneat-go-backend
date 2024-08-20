package models4debtus

import (
	"github.com/dal-go/dalgo/record"
	"time"
)

const RefererKind = "Referer"

type Referer = record.DataWithID[string, *RefererDbo]

type RefererDbo struct {
	Platform   string    `firestore:"p"`
	ReferredTo string    `firestore:"to"`
	ReferredBy string    `firestore:"by"`
	DtCreated  time.Time `firestore:"t"`
}
