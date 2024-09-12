package facade4splitus

import "errors"

var (
	ErrBillHasNoPayer       = errors.New("bill has no payer")
	ErrBillHasTooManyPayers = errors.New("bill has too many payers")
)
