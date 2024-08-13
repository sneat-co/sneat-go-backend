package facade4splitus

import (
	// "github.com/sneat-co/sneat-go-backend/debtus/gae_app/debtus/dtdal"
	// "github.com/sneat-co/sneat-go-backend/debtus/gae_app/debtus/dtdal/dalmocks"
	"fmt"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/gae_app/debtstracker/dtmocks"
	"github.com/sneat-co/sneat-go-backend/src/modules/splitus/briefs4splitus"
	"github.com/sneat-co/sneat-go-backend/src/modules/splitus/models4splitus"

	// "errors"
	"strings"
	"testing"

	"context"
	"github.com/strongo/decimal"
)

func TestCreateBillPanicOnNilContext(t *testing.T) {
	t.Skip("TODO: fix")
	defer func() {
		if r := recover(); r == nil {
			t.Error("No panic")
		} else {
			err := r.(string)
			if !strings.Contains(err, "context.Context") {
				t.Errorf("Error does not mention context: %v", err)
			}
		}
	}()
	_, _ = CreateBill(context.Background(), nil, spaceID, nil)
}

func TestCreateBillPanicOnNilBill(t *testing.T) {
	t.Skip("TODO: fix")
	defer func() {
		if r := recover(); r == nil {
			t.Error("No panic")
		} else {
			err := r.(string)
			if !strings.Contains(err, "*models.BillDbo") {
				t.Errorf("Error does not mention bill: %v", err)
			}
		}
	}()
	_, _ = CreateBill(context.Background(), nil, spaceID, nil)
}

func TestCreateBillErrorNoMembers(t *testing.T) {
	// dtdal.BillEntry = dalmocks.NewBillDalMock()
	// billEntity := createGoodBillSplitByPercentage(t)
	// billEntity.setBillMembers([]models.BillMemberBrief{})
	// bill, err := BillEntry.CreateBill(context.Background(), context.Background(), billEntity)
	// if err != nil {
	// 	if !strings.Contains(err.Error(), "members") {
	// 		t.Error("Error does not mention members:", err)
	// 		return
	// 	}
	// }
	// if bill.ContactID != 1 {
	// 	t.Error("Unexpected bill ContactID:", bill.ContactID)
	// 	return
	// }

}

//const mockBillID = "1"

const spaceID = "testSpace1"

func TestCreateBillAmountZeroError(t *testing.T) {
	t.Skip("TODO: fix")
	dtmocks.SetupMocks(context.Background())
	billEntity := createGoodBillSplitByPercentage(t)
	billEntity.AmountTotal = 0
	billEntity.Currency = "EUR"
	bill, err := CreateBill(context.Background(), nil, spaceID, billEntity)
	if err == nil {
		t.Error("Error expected")
	}
	errText := err.Error()
	if !strings.Contains(errText, "== 0") || !strings.Contains(errText, "AmountTotal") {
		t.Error("Unexpected error text:", errText)
	}
	if bill.ID != "" {
		t.Error("bill.ContactID != empty string")
	}
}

func TestCreateBillAmountNegativeError(t *testing.T) {
	t.Skip("TODO: fix")
	dtmocks.SetupMocks(context.Background())
	billEntity := createGoodBillSplitByPercentage(t)
	billEntity.AmountTotal = -5
	bill, err := CreateBill(context.Background(), nil, spaceID, billEntity)
	if err == nil {
		t.Error("Error expected")
	}
	errText := err.Error()
	if !strings.Contains(errText, "< 0") || !strings.Contains(errText, "AmountTotal") {
		t.Error("Unexpected error text:", errText)
	}
	if bill.ID != "" {
		t.Error("bill.ContactID != empty string")
	}
}

func TestCreateBillAmountError(t *testing.T) {
	dtmocks.SetupMocks(context.Background())
	billEntity := createGoodBillSplitByPercentage(t)
	members := billEntity.GetBillMembers()
	billEntity.AmountTotal += 5
	members[0].Paid += 5
	// billEntity.setBillMembers(members)
	// bill, err := BillEntry.CreateBill(context.Background(), context.Background(), billEntity)
	// if err == nil {
	// 	t.Error("Error expected")
	// }
	// errText := err.Error()
	// if !strings.Contains(errText, "totalOwedByMembers != billEntity.AmountTotal") {
	// 	t.Error("Unexpected error text:", errText)
	// }
	// if bill.ContactID != 0 {
	// 	t.Error("bill.ContactID != 0")
	// }
}

func TestCreateBillStatusMissingError(t *testing.T) {
	t.Skip("TODO: fix")
	dtmocks.SetupMocks(context.Background())
	billEntity := createGoodBillSplitByPercentage(t)
	billEntity.Status = ""
	bill, err := CreateBill(context.Background(), nil, spaceID, billEntity)
	if err == nil {
		t.Error("Error expected")
		return
	}
	errText := err.Error()
	if !strings.Contains(errText, "required") || !strings.Contains(errText, "Status") {
		t.Error("Unexpected error text:", errText)
	}
	if bill.ID != "" {
		t.Error("bill.ContactID != empty string")
	}
}

func TestCreateBillStatusUnknownError(t *testing.T) {
	t.Skip("TODO: fix")
	dtmocks.SetupMocks(context.Background())
	billEntity := createGoodBillSplitByPercentage(t)
	billEntity.Status = "bogus"
	bill, err := CreateBill(context.Background(), nil, spaceID, billEntity)
	if err == nil {
		t.Error("Error expected")
		return
	}
	errText := err.Error()
	if !strings.Contains(errText, "Invalid status") || !strings.Contains(errText, "expected one of") {
		t.Error("Unexpected error text:", errText)
	}
	if bill.ID != "" {
		t.Error("bill.ContactID != empty string")
	}
}

func TestCreateBillMemberNegativeAmountError(t *testing.T) {
	t.Skip("TODO: fix")
	dtmocks.SetupMocks(context.Background())
	billEntity := createGoodBillSplitByPercentage(t)
	members := billEntity.GetBillMembers()
	members[3].Owes *= -1
	// billEntity.setBillMembers(members)
	// billEntity.AmountTotal += members[3].Owes
	// bill, err := BillEntry.CreateBill(context.Background(), context.Background(), billEntity)
	// if err == nil {
	// 	t.Error("Error expected")
	// }
	// errText := err.Error()
	// if !strings.Contains(errText, "negative") || !strings.Contains(errText, "members[3]") {
	// 	t.Error("Unexpected error text:", errText)
	// }
	// if bill.ContactID != 0 {
	// 	t.Error("bill.ContactID != 0")
	// }
}

func TestCreateBillTooManyMembersError(t *testing.T) {
	t.Skip("TODO: fix")
	dtmocks.SetupMocks(context.Background())
	billEntity := createGoodBillSplitByPercentage(t)
	members := billEntity.GetBillMembers()
	members[0].Paid = billEntity.AmountTotal / 2
	members[1].Paid = billEntity.AmountTotal / 2
	members[2].Paid = billEntity.AmountTotal / 2
	if err := billEntity.SetBillMembers(members); err != nil {
		t.Error(err)
	}
	c := context.Background()
	bill, err := CreateBill(c, nil, spaceID, billEntity)
	if err == nil {
		t.Error("Error expected")
	}
	errText := err.Error()
	if errText != "bill has too many payers" {
		t.Error("Unexpected error text:", errText)
	}
	if bill.ID == "0" {
		t.Error("bill.ContactID is empty string")
	}
}

func TestCreateBillMembersOverPaid(t *testing.T) {
	t.Skip("TODO: fix")
	dtmocks.SetupMocks(context.Background())
	billEntity := createGoodBillSplitByPercentage(t)
	members := billEntity.GetBillMembers()
	members[0].Paid = billEntity.AmountTotal + 10
	if err := billEntity.SetBillMembers(members); err != nil {
		t.Fatal(err)
	}
	bill, err := CreateBill(context.Background(), nil, spaceID, billEntity)
	if err == nil {
		t.Fatal("Error expected")
	}
	errText := err.Error()
	if !strings.Contains(errText, "Total paid") || !strings.Contains(errText, "equal") {
		t.Error("Unexpected error text:", errText)
	}
	if bill.ID != "" {
		t.Error("bill.ContactID != empty string")
	}
}

var verifyMemberUserID = func(t *testing.T, members []*briefs4splitus.BillMemberBrief, i int, expectedUserID string) {
	member := members[i]
	if member.UserID != expectedUserID {
		t.Errorf("members[%d].UserID == %v, expected: %v, member: %+v", i, member.UserID, expectedUserID, member)
	}
}

var verifyMemberOwes = func(t *testing.T, members []*briefs4splitus.BillMemberBrief, i int, expecting decimal.Decimal64p2) {
	member := members[i]
	if member.Owes != expecting {
		t.Errorf("members[%d].Owes:%v == %v", i, member.Owes, expecting)
	}
}

func TestCreateBillSuccess(t *testing.T) {
	t.Skip("TODO: fix")
	c := context.Background()
	dtmocks.SetupMocks(c)
	billDbo := createGoodBillSplitByPercentage(t)

	//members := billDbo.GetBillMembers()

	bill, err := CreateBill(c, nil, spaceID, billDbo)
	if err != nil {
		t.Error(err)
		return
	}
	if bill.ID == "" {
		t.Error("Unexpected bill ContactID", bill.ID)
		return
	}

	members := billDbo.GetBillMembers()
	if err != nil {
		t.Error(err)
		return
	}
	if len(members) != len(billDbo.Members) {
		t.Error("len(members) != billDbo.MembersCount")
	}

	verifyMemberUserID(t, members, 0, "1")
	verifyMemberUserID(t, members, 1, "3")
	verifyMemberUserID(t, members, 2, "5")
	verifyMemberUserID(t, members, 3, "")

	// if len(mockDB.BillMock.Bills) != 1 {
	// 	t.Errorf("Expected to have 1 bill in DB, got: %d", len(mockDB.BillMock.Bills))
	// }
}

func createGoodBillSplitByPercentage(t *testing.T) (billEntity *models4splitus.BillDbo) {
	billEntity = new(models4splitus.BillDbo)
	billEntity.Status = models4splitus.BillStatusOutstanding
	billEntity.SplitMode = models4splitus.SplitModePercentage
	billEntity.CreatorUserID = "1"
	billEntity.AmountTotal = 848
	billEntity.Currency = "EUR"

	percent := 25
	if err := billEntity.SetBillMembers([]*briefs4splitus.BillMemberBrief{
		{Percent: 2500, MemberBrief: briefs4splitus.MemberBrief{ID: "1", Shares: percent, UserID: "1", Name: "First member"}, Paid: billEntity.AmountTotal},
		{Percent: 2500, MemberBrief: briefs4splitus.MemberBrief{ID: "2", Shares: percent, UserID: "3", Name: "Second contact", ContactByUser: briefs4splitus.MemberContactBriefsByUserID{"1": briefs4splitus.MemberContactBrief{ContactID: "2", ContactName: "Second contact"}}}},
		{Percent: 2500, MemberBrief: briefs4splitus.MemberBrief{ID: "3", Shares: percent, UserID: "5", Name: "Fifth user", ContactByUser: briefs4splitus.MemberContactBriefsByUserID{"1": briefs4splitus.MemberContactBrief{ContactID: "4", ContactName: "Forth contact"}}}},
		{Percent: 2500, MemberBrief: briefs4splitus.MemberBrief{ID: "4", Shares: percent, Name: "12th contact", ContactByUser: briefs4splitus.MemberContactBriefsByUserID{"5": briefs4splitus.MemberContactBrief{ContactID: "12", ContactName: "12th contact"}}}},
	}); err != nil {
		t.Error(fmt.Errorf("%w: Failed to set members", err))
		return
	}
	return
}

func createGoodBillSplitEqually(t *testing.T) (billEntity *models4splitus.BillDbo, err error) {
	billEntity = new(models4splitus.BillDbo)
	billEntity.Status = models4splitus.BillStatusOutstanding
	billEntity.SplitMode = models4splitus.SplitModeEqually
	billEntity.CreatorUserID = "1"
	billEntity.AmountTotal = 637
	billEntity.Currency = "EUR"

	if err = billEntity.SetBillMembers([]*briefs4splitus.BillMemberBrief{
		{Owes: 213, MemberBrief: briefs4splitus.MemberBrief{ID: "1", UserID: "1", Name: "First user"}, Paid: billEntity.AmountTotal},
		{Owes: 212, MemberBrief: briefs4splitus.MemberBrief{ID: "2", Name: "Second", ContactByUser: briefs4splitus.MemberContactBriefsByUserID{"1": briefs4splitus.MemberContactBrief{ContactID: "2"}}}},
		{Owes: 212, MemberBrief: briefs4splitus.MemberBrief{ID: "3", Name: "Forth", ContactByUser: briefs4splitus.MemberContactBriefsByUserID{"1": briefs4splitus.MemberContactBrief{ContactID: "4"}}}},
	}); err != nil {
		err = fmt.Errorf("%w: Failed to set members", err)
		return
	}
	return
}

func createGoodBillSplitEquallyWithAdjustments(t *testing.T) (billEntity *models4splitus.BillDbo, err error) {
	t.Helper()

	if billEntity, err = createGoodBillSplitEqually(t); err != nil {
		return
	}

	members := billEntity.GetBillMembers()
	members[1].Adjustment = 10
	members[2].Adjustment = 20
	if err = billEntity.SetBillMembers(members); err != nil {
		t.Fatal(err)
	}
	members = billEntity.GetBillMembers()
	if len(members) != 3 {
		t.Fatal("len(members) != 3")
	}
	/*
		637 - 30 = 607
		607 / 3 = 202
	*/
	validateOwes := func(i int, expecting decimal.Decimal64p2) {
		if members[i].Owes != expecting {
			t.Fatalf("members[%d].Owes:%v != %v", i, members[0].Owes, expecting)
		}
	}
	validateOwes(0, 203)
	validateOwes(1, 212)
	validateOwes(2, 222)
	return
}

func createGoodBillSplitByShare(t *testing.T) (billEntity *models4splitus.BillDbo, err error) {
	billEntity = new(models4splitus.BillDbo)
	billEntity.Status = models4splitus.BillStatusOutstanding
	billEntity.SplitMode = models4splitus.SplitModeShare
	billEntity.CreatorUserID = "1"
	billEntity.AmountTotal = 636
	billEntity.Currency = "EUR"

	if err = billEntity.SetBillMembers([]*briefs4splitus.BillMemberBrief{
		{MemberBrief: briefs4splitus.MemberBrief{ID: "1", Shares: 2, UserID: "1", Name: "First user"}, Paid: billEntity.AmountTotal},
		{MemberBrief: briefs4splitus.MemberBrief{ID: "2", Shares: 1, Name: "Second", ContactByUser: briefs4splitus.MemberContactBriefsByUserID{"1": briefs4splitus.MemberContactBrief{ContactID: "2"}}}},
		{MemberBrief: briefs4splitus.MemberBrief{ID: "3", Shares: 3, Name: "Forth", ContactByUser: briefs4splitus.MemberContactBriefsByUserID{"1": briefs4splitus.MemberContactBrief{ContactID: "4"}}}},
	}); err != nil {
		t.Error(fmt.Errorf("%w: Failed to set members", err))
		return
	}
	members := billEntity.GetBillMembers()
	verifyMemberOwes(t, members, 0, 212)
	verifyMemberOwes(t, members, 1, 106)
	verifyMemberOwes(t, members, 2, 318)
	return
}

// There is no way to check as we do not expose membser publicly
// func TestCreateBillEquallyTooManyAmountsError(t *testing.T) {
// 	c := context.Background()
// 	dtmocks.SetupMocks(c)
// 	billEntity, err := createGoodBillSplitEqually(t)
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}
// 	members := billEntity.GetBillMembers()
// 	members[1].Owes -= decimal.NewDecimal64p2FromFloat64(0.01)
// 	t.Logf("memebers: %v", members)
// 	if err = billEntity.SetBillMembers(members); err != nil {
// 		t.Fatal(err)
// 	}
// 	bill, err := BillEntry.CreateBill(c, c, billEntity)
// 	if err == nil {
// 		t.Fatal("Error expected")
// 	}
// 	errText := err.Error()
// 	if !strings.Contains(errText, "len(amountsCountByValue) > 2") {
// 		t.Error("Unexpected error text:", errText)
// 	}
// 	if bill.ContactID == "" {
// 		t.Error("bill.ContactID is empty string")
// 	}
// }

// func TestCreateBillEquallyAmountDeviateTooMuchError(t *testing.T) {
// 	c := context.Background()
// 	dtmocks.SetupMocks(c)
// 	billEntity, err := createGoodBillSplitEqually(t)
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}
// 	members := billEntity.GetBillMembers()
// 	members[0].Owes += decimal.NewDecimal64p2FromFloat64(0.01)
// 	members[1].Owes -= decimal.NewDecimal64p2FromFloat64(0.01)
// 	if err = billEntity.SetBillMembers(members); err != nil {
// 		t.Fatal(err)
// 	}
// 	bill, err := BillEntry.CreateBill(context.Background(), context.Background(), billEntity)
// 	if err == nil {
// 		t.Error("Error expected")
// 		return
// 	}
// 	errText := err.Error()
// 	if !strings.Contains(errText, "deviated too much") {
// 		t.Error("Unexpected error text:", errText)
// 	}
// 	if bill.ContactID == "" {
// 		t.Error("bill.ContactID is empty string")
// 	}
// }

func TestCreateBillEquallySuccess(t *testing.T) {
	t.Skip("TODO: fix")
	c := context.Background()
	dtmocks.SetupMocks(c)
	billEntity, err := createGoodBillSplitEqually(t)
	if err != nil {
		t.Error(err)
		return
	}
	bill, err := CreateBill(c, nil, spaceID, billEntity)
	if err != nil {
		t.Error(err)
		return
	}
	if bill.ID == "" {
		t.Error("bill.ContactID is empty string")
	}
}

func TestCreateBillAdjustmentSuccess(t *testing.T) {
	t.Skip("TODO: fix")
	c := context.Background()
	dtmocks.SetupMocks(c)
	billEntity, err := createGoodBillSplitEquallyWithAdjustments(t)
	if err != nil {
		t.Fatal(err)
	}
	bill, err := CreateBill(c, nil, spaceID, billEntity)
	if err != nil {
		t.Error(err)
		return
	}
	if bill.ID == "" {
		t.Error("bill.ContactID is empty string")
	}
}

func TestCreateBillAdjustmentTotalAdjustmentIsTooBigError(t *testing.T) {
	t.Skip("TODO: fix")
	dtmocks.SetupMocks(context.Background())
	billEntity, err := createGoodBillSplitEquallyWithAdjustments(t)
	if err != nil {
		return
	}
	members := billEntity.GetBillMembers()
	members[1].Adjustment += decimal.NewDecimal64p2FromFloat64(4.15)
	members[2].Adjustment += decimal.NewDecimal64p2FromFloat64(3.16)
	// billEntity.setBillMembers(members)
	// bill, err := BillEntry.CreateBill(context.Background(), context.Background(), billEntity)
	// if err == nil {
	// 	t.Error("Error expected")
	// 	return
	// }
	// errText := err.Error()
	// if !strings.Contains(errText, "totalAdjustmentByMembers > billEntity.AmountTotal") {
	// 	t.Error("Unexpected error text:", errText)
	// }
	// if bill.ContactID != 0 {
	// 	t.Error("bill.ContactID != 0")
	// }
}

func TestCreateBillAdjustmentMemberAdjustmentIsTooBigError(t *testing.T) {
	t.Skip("TODO: fix")
	dtmocks.SetupMocks(context.Background())
	billEntity, err := createGoodBillSplitEquallyWithAdjustments(t)
	if err != nil {
		return
	}
	members := billEntity.GetBillMembers()
	members[1].Adjustment += decimal.NewDecimal64p2FromFloat64(7.19)
	// billEntity.setBillMembers(members)
	// bill, err := BillEntry.CreateBill(context.Background(), context.Background(), billEntity)
	// if err == nil {
	// 	t.Error("Error expected")
	// 	return
	// }
	// errText := err.Error()
	// if !strings.Contains(errText, "members[1].Adjustment > billEntity.AmountTotal") {
	// 	t.Error("Unexpected error text:", errText)
	// }
	// if bill.ContactID != 0 {
	// 	t.Error("bill.ContactID != 0")
	// }
}

func TestCreateBillAdjustmentAmountDeviateTooMuchError(t *testing.T) {
	t.Skip("TODO: fix")
	dtmocks.SetupMocks(context.Background())
	billEntity, err := createGoodBillSplitEquallyWithAdjustments(t)
	if err != nil {
		return
	}
	members := billEntity.GetBillMembers()
	members[1].Adjustment += decimal.NewDecimal64p2FromFloat64(0.10)
	// billEntity.setBillMembers(members)
	// bill, err := BillEntry.CreateBill(context.Background(), context.Background(), billEntity)
	// if err == nil {
	// 	t.Error("Error expected")
	// 	return
	// }
	// errText := err.Error()
	// if !strings.Contains(errText, "deviated too much") {
	// 	t.Error("Unexpected error text:", errText)
	// }
	// if bill.ContactID != 0 {
	// 	t.Error("bill.ContactID != 0")
	// }
}

func TestCreateBillShareSuccess(t *testing.T) {
	t.Skip("TODO: fix")
	dtmocks.SetupMocks(context.Background())
	billEntity, err := createGoodBillSplitByShare(t)
	if err != nil {
		return
	}
	bill, err := CreateBill(context.Background(), nil, spaceID, billEntity)
	if err != nil {
		t.Error(err)
		return
	}
	if bill.ID == "" {
		t.Error("bill.ContactID is empty string")
	}
}

func TestCreateBillShareAmountDeviateTooMuchError(t *testing.T) {
	t.Skip("TODO: fix")
	//mockDB := dtmocks.SetupMocks(context.Background())
	billEntity, err := createGoodBillSplitEquallyWithAdjustments(t)
	if err != nil {
		return
	}
	members := billEntity.GetBillMembers()
	members[1].Owes += decimal.NewDecimal64p2FromFloat64(0.10)
	members[2].Owes -= decimal.NewDecimal64p2FromFloat64(0.10)
	//billEntity.setBillMembers(members)
	//bill, err := BillEntry.CreateBill(context.Background(), context.Background(), billEntity)
	//if err == nil {
	//	t.Error("Error expected")
	//	return
	//}
	//errText := err.Error()
	//if !strings.Contains(errText, "deviated too much") {
	//	t.Error("Unexpected error text:", errText)
	//}
	//if bill.ContactID != 0 {
	//	t.Error("bill.ContactID != 0")
	//}
	//if len(mockDB.BillMock.Bills) != 0 {
	//	t.Errorf("Expected to have 0 bills in database, got: %d", len(mockDB.BillMock.Bills))
	//}
}
