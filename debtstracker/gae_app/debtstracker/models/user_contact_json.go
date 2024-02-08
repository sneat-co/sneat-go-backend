package models

//go:generate ffjson $GOFILE

import (
	"encoding/json"
	"fmt"
	"time"

	"context"
	"errors"
	"github.com/crediterra/money"
	"github.com/pquerna/ffjson/ffjson"
	"github.com/strongo/decimal"
)

type LastTransfer struct {
	ID string    `datastore:"LastTransferID,noindex"`
	At time.Time `datastore:"LastTransferAt,noindex"`
}

type TransferWithInterestJson struct {
	TransferInterest
	TransferID string
	Direction  TransferDirection
	Starts     time.Time
	Currency   money.CurrencyCode `json:",omitempty"` // TODO: will be obsolete once we group outstanding by currency
	Amount     decimal.Decimal64p2
	Returns    TransferReturns `json:",omitempty"`
}

func (t TransferWithInterestJson) Equal(t2 TransferWithInterestJson) bool {
	if t.TransferInterest != t2.TransferInterest ||
		t.TransferID != t2.TransferID ||
		t.Starts != t2.Starts ||
		t.Currency != t2.Currency ||
		t.Amount != t2.Amount ||
		len(t.Returns) != len(t2.Returns) {
		return false
	}
	for i, r := range t.Returns {
		if r != t2.Returns[i] {
			return false
		}
	}
	return true
}

var _ = TransferInterestCalculable(TransferWithInterestJson{})

func (t TransferWithInterestJson) GetLendingValue() decimal.Decimal64p2 {
	return t.Amount
}

func (t TransferWithInterestJson) GetStartDate() time.Time {
	return t.Starts
}

func (t TransferWithInterestJson) GetReturns() (returns TransferReturns) {
	return t.Returns
}

type UserContactTransfersInfo struct {
	Count                   int                        `json:",omitempty"`
	Last                    LastTransfer               `json:",omitempty"`
	OutstandingWithInterest []TransferWithInterestJson `json:",omitempty"`
}

func (o *UserContactTransfersInfo) Equal(o2 *UserContactTransfersInfo) bool {
	if o2 == nil || o.Count != o2.Count || o.Last != o2.Last || len(o.OutstandingWithInterest) != len(o2.OutstandingWithInterest) {
		return false
	}
	for i, t := range o.OutstandingWithInterest {
		if !t.Equal(o2.OutstandingWithInterest[i]) {
			return false
		}
	}
	return true
}

type UserContactJson struct {
	ID          string
	Name        string
	Status      string                    `json:",omitempty"`
	UserID      string                    `json:",omitempty"` // TODO: new prop, update in map reduce and change code!
	TgUserID    int64                     `json:",omitempty"`
	BalanceJson *json.RawMessage          `json:"Balance,omitempty"`
	Transfers   *UserContactTransfersInfo `json:",omitempty"`
}

func (o UserContactJson) Equal(o2 UserContactJson) bool {
	return o.ID == o2.ID &&
		o.Name == o2.Name &&
		o.Status == o2.Status &&
		o.UserID == o2.UserID &&
		o.TgUserID == o2.TgUserID &&
		((o.BalanceJson == nil && o2.BalanceJson == nil) || (o.BalanceJson != nil && o2.BalanceJson != nil && string(*o.BalanceJson) == string(*o2.BalanceJson))) &&
		((o.Transfers == nil && o2.Transfers == nil) || (o.Transfers != nil && o2.Transfers != nil && o.Transfers.Equal(o2.Transfers)))
}

func (o UserContactJson) Balance() (balance money.Balance) {
	balance = make(money.Balance)
	if o.BalanceJson == nil {
		return
	}
	if err := ffjson.Unmarshal(*o.BalanceJson, &balance); err != nil { // TODO: Migrate to ffjson.UnmarshalFast() ?
		panic(err)
	}
	return
}

func (o *UserContactJson) SetBalance(balance money.Balance) (err error) {
	for c, v := range balance {
		if v == 0 {
			return errors.New("balance is zero for currency: " + string(c))
		}
	}
	if data, err := ffjson.Marshal(balance); err != nil {
		return err
	} else {
		rawJson := json.RawMessage(data)
		o.BalanceJson = &rawJson
	}
	return
}

func (o UserContactJson) BalanceWithInterest(c context.Context, periodEnds time.Time) (balance money.Balance, err error) {
	balance = o.Balance()
	if o.Transfers != nil {
		if err = updateBalanceWithInterest(false, balance, o.Transfers.OutstandingWithInterest, periodEnds); err != nil {
			return
		}
	}
	return
}

func NewUserContactJson(counterpartyID string, status, name string, balanced money.Balanced) UserContactJson {
	result := UserContactJson{
		ID:     counterpartyID,
		Status: status,
		Name:   name,
	}
	if balanced.BalanceJson != "" {
		balance := json.RawMessage(balanced.BalanceJson)
		result.BalanceJson = &balance
	}
	if balanced.CountOfTransfers != 0 {
		if balanced.LastTransferID == "" {
			panic(fmt.Sprintf("balanced.CountOfTransfers:%v != 0 && balanced.LastTransferID == 0", balanced.CountOfTransfers))
		}
		if balanced.LastTransferAt.IsZero() {
			panic(fmt.Sprintf("balanced.CountOfTransfers:%v != 0 && balanced.LastTransferAt.IsZero():true", balanced.CountOfTransfers))
		}
		result.Transfers = &UserContactTransfersInfo{
			Count: balanced.CountOfTransfers,
			Last: LastTransfer{
				ID: balanced.LastTransferID,
				At: balanced.LastTransferAt,
			},
		}
	}
	return result
}
