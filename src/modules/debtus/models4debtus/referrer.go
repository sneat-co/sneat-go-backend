package models4debtus

import (
	"github.com/dal-go/dalgo/record"
	"time"
)

const RefererKind = "Referer"

type Referer = record.DataWithID[string, *RefererDbo]

type RefererDbo struct {
	Platform   string    `datastore:"p"`
	ReferredTo string    `datastore:"to"`
	ReferredBy string    `datastore:"by"`
	DtCreated  time.Time `datastore:"t"`
}
