package models

import "github.com/strongo/decimal"

type SmsStats struct {
	SmsCount   int64               `datastore:",noindex,omitempty"`
	SmsCost    float32             `datastore:",noindex,omitempty"`
	SmsCostUSD decimal.Decimal64p2 `datastore:",noindex,omitempty"`
}
