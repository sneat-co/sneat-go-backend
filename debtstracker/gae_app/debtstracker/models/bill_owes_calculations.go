package models

import (
	"fmt"

	"github.com/strongo/decimal"
)

func (entity *BillCommon) updateMemberOwes(members []BillMemberJson) (err error) {

	switch entity.SplitMode {
	case SplitModeEqually:
		updateMemberOwesForEqualSplit(entity.AmountTotal, entity.CreatorUserID, members)
	case SplitModeShare:
		updateMemberOwesForSplitByShares(entity.AmountTotal, entity.CreatorUserID, members)
	case SplitModePercentage:
		updateMemberOwesForSplitByPercentage(entity.AmountTotal, entity.CreatorUserID, members)
	default:
		err = fmt.Errorf("%w: %s", ErrUnknownSplitMode, string(entity.SplitMode))
	}
	return
}

func updateMemberOwesForEqualSplit(amountTotal decimal.Decimal64p2, creatorUserID string, members []BillMemberJson) {
	membersCount := int64(len(members))
	if membersCount == 0 {
		return
	}
	adjustedTotal := amountTotal
	for i := range members {
		if adjustedTotal -= members[i].Adjustment; adjustedTotal < 0 {
			panic("adjustedTotal < 0")
		}
	}
	perMember := decimal.Decimal64p2(int64(adjustedTotal) / membersCount)

	getRemainder := func() decimal.Decimal64p2 {
		return adjustedTotal - decimal.Decimal64p2(int64(perMember)*membersCount)
	}

	remainder := getRemainder()

	// if we need this add comment on why
	// for remainder > 1 || remainder < -1 {
	// 	switch {
	// 	case remainder > 1:
	// 		perMember += 1
	// 	case remainder < -1:
	// 		perMember -= 1
	// 	}
	// 	remainder = getRemainder()
	// }

	if remainder > 1 || remainder < -1 {
		panic("remainder > 1 || remainder < -1")
	}

	creatorIndex := -1
	for i, member := range members {
		members[i].Owes = perMember + member.Adjustment
		if creatorUserID != "" && member.UserID == creatorUserID {
			if creatorIndex == -1 {
				creatorIndex = i
			} else {
				panic(fmt.Sprintf("creator is listed as member twice with UserID=%v", creatorUserID))
			}
		}
	}

	if remainder != 0 {
		i := creatorIndex
		if i < 0 {
			i = 0
		}
		members[i].Owes += remainder
	}
	// fixTotal(adjustedTotal, members)
}

func updateMemberOwesForSplitByPercentage(amountTotal decimal.Decimal64p2, creatorUserID string, members []BillMemberJson) {
	updateMemberOwesForSplitByShares(amountTotal, creatorUserID, members)
}

func updateMemberOwesForSplitByShares(amountTotal decimal.Decimal64p2, creatorUserID string, members []BillMemberJson) {
	membersCount := len(members)
	if membersCount == 0 {
		return
	}

	totalShares := 0

	for _, m := range members {
		if m.Shares < 0 {
			panic(fmt.Sprintf("m.Shares < 0: %v", m.Shares))
		}
		totalShares += m.Shares
	}

	if totalShares == 0 {
		totalShares = 10 * membersCount
		for i := range members {
			members[i].Shares = 10
		}
	}

	perShareBy100 := float64(amountTotal) / float64(totalShares) * 100

	getRemainder := func() decimal.Decimal64p2 {
		return amountTotal - decimal.Decimal64p2(int64(perShareBy100)*int64(totalShares)/100)
	}

	remainder := getRemainder()

	for remainder > 1 || remainder < -1 {
		switch {
		case remainder > 1:
			perShareBy100 += 1
		case remainder < -1:
			perShareBy100 -= 1
		}
		remainder = getRemainder()
	}

	creatorIndex := -1
	for i := range members {
		members[i].Owes = decimal.Decimal64p2(perShareBy100 * float64(members[i].Shares) / 100)
		if creatorUserID != "" && members[i].UserID == creatorUserID {
			creatorIndex = i
		}
	}
	if remainder != 0 {
		i := creatorIndex
		if i < 0 {
			i = 0
		}
		members[i].Owes += remainder
	}
	fixTotal(amountTotal, members)
}

func fixTotal(amountTotal decimal.Decimal64p2, members []BillMemberJson) {
	var total decimal.Decimal64p2
	for _, member := range members {
		total += member.Owes
	}
	switch amountTotal - total {
	case 0:
	case 1:
		// Let's ad remainder to a members with smallest amount
		var (
			idx                     = -1
			min decimal.Decimal64p2 = 1<<63 - 1
		)
		for i, m := range members {
			if m.Owes < min {
				idx = i
			}
		}
		members[idx].Owes += 1
	case -1:
		// Let's ad deduct remainder from a members with largest amount
		var (
			idx = -1
			max decimal.Decimal64p2
		)
		for i, m := range members {
			if m.Owes > max {
				idx = i
			}
		}
		members[idx].Owes += 1
	default:
		panic(fmt.Sprintf("Remainder is too big: %v", amountTotal-total))
	}
}
