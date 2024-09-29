package models4splitus

import (
	"errors"
	"github.com/sneat-co/sneat-mod-debtus-go/debtus/const4debtus"
)

type SplitMode string

const (
	SplitModeEqually     SplitMode = "equally"
	SplitModeExactAmount SplitMode = "exact-amount"
	SplitModePercentage  SplitMode = "percentage"
	SplitModeShare       SplitMode = "shares"

	// SplitModeAdjustment  SplitMode = "adjustment"
)

var ErrUnknownSplitMode = errors.New("Unknown split mode")

type PayMode string

const (
	PayModePrepay  = "prepay"
	PayModeBillpay = "billpay"
)

const (
	BillKind = "BillEntry"
)

const (
	BillStatusDraft       = const4debtus.StatusDraft
	BillStatusDeleted     = const4debtus.StatusDeleted
	BillStatusOutstanding = "outstanding"
	BillStatusSettled     = "settled"
)

var (
	BillStatuses = [3]string{
		BillStatusDraft,
		BillStatusOutstanding,
		BillStatusSettled,
	}
	BillSplitModes = [4]SplitMode{
		// SplitModeAdjustment,
		SplitModeEqually,
		SplitModeExactAmount,
		SplitModePercentage,
		SplitModeShare,
	}
)

func IsValidBillSplit(split SplitMode) bool {
	for _, v := range BillSplitModes {
		if split == v {
			return true
		}
	}
	return false
}

func IsValidBillStatus(status string) bool {
	for _, v := range BillStatuses {
		if status == v {
			return true
		}
	}
	return false
}

const (
	BillsHistoryCollection = "bill_history"
	SplitsCollection       = "splits"
)
