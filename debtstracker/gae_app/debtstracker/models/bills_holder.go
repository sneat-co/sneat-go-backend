package models

import (
	"fmt"
	"github.com/pquerna/ffjson/ffjson"
)

type billsHolder struct {
	OutstandingBillsCount int    `datastore:",noindex,omitempty"`
	OutstandingBillsJson  string `datastore:",noindex,omitempty"`
}

func (entity *billsHolder) GetOutstandingBills() (outstandingBills []BillJson) {
	if entity.OutstandingBillsJson == "" {
		return
	}
	if err := ffjson.Unmarshal([]byte(entity.OutstandingBillsJson), &outstandingBills); err != nil {
		panic(err)
	}
	if entity.OutstandingBillsCount != len(outstandingBills) {
		panic(fmt.Errorf("%w:len([]BillJson) != OutstandingBillsCount", ErrJsonCountMismatch))
	}
	return
}

func (entity *billsHolder) SetOutstandingBills(outstandingBills []BillJson) (changed bool, err error) {
	var data []byte
	if data, err = ffjson.Marshal(outstandingBills); err != nil {
		return
	}
	json := string(data)
	if json == "[]" {
		json = ""
	}
	entity.OutstandingBillsCount = len(outstandingBills)
	if changed = json != entity.OutstandingBillsJson; changed {
		entity.OutstandingBillsJson = json
	}
	return
}

//func init() {
//	userPropertiesToClean["OutstandingBillsJson"] = gaedb.IsEmptyJSON
//	groupPropertiesToClean["OutstandingBillsJson"] = gaedb.IsEmptyJSON
//}
