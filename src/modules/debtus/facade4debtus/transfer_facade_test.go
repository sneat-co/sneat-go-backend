package facade4debtus

import (
	"context"
	"fmt"
	"github.com/bots-go-framework/bots-fw-telegram"
	"github.com/crediterra/go-interest"
	"github.com/crediterra/money"
	"github.com/sneat-co/sneat-go-backend/src/coremodules/spaceus/dto4spaceus"
	"github.com/sneat-co/sneat-go-backend/src/coremodules/userus/dbo4userus"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/gae_app/debtstracker/dtdal"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/gae_app/debtstracker/dtmocks"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/models4debtus"
	"github.com/strongo/decimal"
	"github.com/strongo/strongoapp"
	"runtime/debug"
	"strings"
	"testing"
	"time"
)

type assertHelper struct {
	t *testing.T
}

func (assert assertHelper) OutputIsNilIfErr(output CreateTransferOutput, err error) (CreateTransferOutput, error) {
	assert.t.Helper()
	if err != nil {
		if output.Transfer.ID != "" {
			assert.t.Errorf("Returned transfer.ContactID != 0 with error: (ContactID=%v), error: %v", output.Transfer.ID, err)
		}
		if output.Transfer.Data != nil {
			assert.t.Errorf("Returned a non nil transfer entity with error: %v", err)
		}
		// if counterparty != nil {
		// 	t.Errorf("Returned a counterparty with error: %v", counterparty)
		// }
		// if creatorUser != nil {
		// 	t.Errorf("Returned creatorUser with error: %v", creatorUser)
		// 	return
		// }
	}
	return output, err
}

func TestCreateTransfer(t *testing.T) {
	t.Skip("TODO: fix")
	// c, done, err := aetest.NewContext()
	// if err != nil {
	// 	t.Fatal(err.Error())
	// }
	// defer done()
	dtmocks.SetupMocks(context.Background())
	c := context.Background()
	assert := assertHelper{t: t}

	source := dtdal.NewTransferSourceBot(telegram.PlatformID, "test-bot", "444")

	currency := money.CurrencyEUR

	const (
		userID         = "1"
		counterpartyID = "2"
	)

	/* Test CreateTransfer that should succeed  - new counterparty by name */
	{
		// counterparty, err := dtdal.DebtusSpaceContactEntry.CreateContact(c, user.ContactID, 0, 0, models.ContactDetails{
		// 	FirstName: "First",
		// 	LastName:  "DebtusSpaceContactEntry",
		// })

		from := &models4debtus.TransferCounterpartyInfo{
			UserID:  userID,
			Note:    "Note 1",
			Comment: "Comment 1",
		}

		to := &models4debtus.TransferCounterpartyInfo{
			ContactID: counterpartyID,
		}

		creatorUser := dbo4userus.NewUserEntry(userID)

		noInterest := models4debtus.NoInterest()

		dueOn := time.Now().Add(time.Minute)

		request := CreateTransferRequest{
			SpaceRequest: dto4spaceus.SpaceRequest{},
			Amount:       money.NewAmount(currency, 10),
			DueOn:        &dueOn,
			Interest:     &noInterest,
			IsReturn:     false,
		}
		newTransfer := NewTransferInput(
			strongoapp.LocalHostEnv,
			source,
			creatorUser,
			request,
			from, to,
		)

		output, err := assert.OutputIsNilIfErr(Transfers.CreateTransfer(c, newTransfer))
		if err != nil {
			t.Error("Should not fail:", err)
			return
		}
		transfer := output.Transfer
		fromUser, toUser := output.From.User, output.To.User
		fromCounterparty, toCounterparty := output.From.Contact, output.To.Contact

		if output.Transfer.ID == "" {
			t.Error("transfer.ContactID == 0")
			return
		}
		if output.Transfer.Data == nil {
			t.Error("transfer.TransferData == nil")
			return
		}

		if toCounterparty.Data == nil {
			t.Error("toCounterparty.DebtusSpaceContactDbo == nil")
			return
		}
		if fromUser.Data == nil {
			t.Error("fromUser.DebutsAppUserDataOBSOLETE == nil")
			return
		}
		if transfer.Data.CreatorUserID != userID {
			t.Errorf("transfer.CreatorUserID:%v != userID:%v", transfer.Data.CreatorUserID, userID)
		}
		if transfer.Data.Counterparty().ContactID != to.ContactID {
			t.Errorf("transfer.DebtusSpaceContactEntry().ContactID:%v != to.ContactID:%v", transfer.Data.Counterparty().ContactID, to.ContactID)
		}
		if transfer.Data.Counterparty().ContactName == "" {
			t.Error("transfer.DebtusSpaceContactEntry().ContactName is empty string")
		}

		transfer2, err := Transfers.GetTransferByID(c, nil, transfer.ID)
		if err != nil {
			t.Error(fmt.Errorf("failed to get transfer by id=%v: %w", transfer.ID, err))
			return
		}
		if transfer.Data == nil {
			t.Error("transfer == nil")
			return
		}
		if len(transfer2.Data.BothUserIDs) != 2 {
			t.Errorf("len(transfer2.BothUserIDs):%v != 2", len(transfer2.Data.BothUserIDs))
			return
		}
		if fromUser.ID != "" && fromUser.Data == nil {
			t.Error("fromUser.ContactID != 0 && fromUser.DebutsAppUserDataOBSOLETE == nil")
		}
		if toUser.ID != "" && toUser.Data == nil {
			t.Error("toUser.ContactID != 0 && toUser.DebutsAppUserDataOBSOLETE == nil")
		}
		if toCounterparty.ID != "" && toCounterparty.Data == nil {
			t.Error("fromCounterparty.DebtusSpaceContactDbo == nil")
		}
		if fromCounterparty.ID != "" && fromCounterparty.Data == nil {
			t.Error("fromCounterparty.ContactID != 0 && fromCounterparty.DebtusSpaceContactDbo == nil")
		}
	}
}

func Test_removeClosedTransfersFromOutstandingWithInterest(t *testing.T) {
	transfersWithInterest := []models4debtus.TransferWithInterestJson{
		{TransferID: "1"},
		{TransferID: "2"},
		{TransferID: "3"},
		{TransferID: "4"},
		{TransferID: "5"},
	}
	transfersWithInterest = removeClosedTransfersFromOutstandingWithInterest(transfersWithInterest, []string{"2", "3"})
	if len(transfersWithInterest) != 3 {
		t.Fatalf("len(transfersWithInterest) != 3: %v", transfersWithInterest)
	}
	for i, transferID := range []string{"1", "4", "5"} {
		if transfersWithInterest[i].TransferID != transferID {
			t.Fatalf("transfersWithInterest[%v].TransferID: %v != %v", i, transfersWithInterest[i].TransferID, transferID)
		}
	}
}

type createTransferTestCase struct {
	name  string
	steps []createTransferStep
}

type createTransferStepInput struct {
	direction          models4debtus.TransferDirection
	isReturn           bool
	returnToTransferID string
	amount             decimal.Decimal64p2
	time               time.Time
	models4debtus.TransferInterest
}

type createTransferStepExpects struct {
	isExpectedError func(error) bool // in case of error rest are no validated

	balance        decimal.Decimal64p2
	amountReturned decimal.Decimal64p2
	isOutstanding  bool

	// we should have both `returns` and `api4transfers` properties to assert partial returns
	returns   models4debtus.TransferReturns
	transfers []transferExpectedState
}

type createTransferStep struct {
	input             createTransferStepInput
	createdTransferID string // Save transfer ContactID for checking transfer state in later steps
	expects           createTransferStepExpects
}

type transferExpectedState struct {
	stepIndex      int
	isOutstanding  bool
	amountReturned decimal.Decimal64p2
}

func TestCreateTransfers(t *testing.T) {
	dayHour := func(d, h int) time.Time {
		return time.Date(2000, 01, d, h, 0, 0, 0, time.Local)
	}

	expects := func(
		balance decimal.Decimal64p2,
		amountReturned decimal.Decimal64p2,
		isOutstanding bool,
		returns models4debtus.TransferReturns,
		transfers []transferExpectedState,
	) createTransferStepExpects {
		return createTransferStepExpects{
			balance:        balance,
			amountReturned: amountReturned,
			isOutstanding:  isOutstanding,
			returns:        returns,
			transfers:      transfers,
		}
	}

	expectsError := func(f func(error) bool) createTransferStepExpects {
		return createTransferStepExpects{
			isExpectedError: f,
		}
	}

	testCases := []createTransferTestCase{
		{
			name: "error on attempt to make a return with interest",
			steps: []createTransferStep{
				{
					input: createTransferStepInput{
						direction: models4debtus.TransferDirectionUser2Counterparty,
						amount:    10,
						time:      dayHour(1, 1),
					},
					expects: expects(10, 0, true, nil, nil),
				},
				{
					input: createTransferStepInput{
						direction:        models4debtus.TransferDirectionCounterparty2User,
						amount:           10,
						time:             dayHour(1, 1),
						TransferInterest: models4debtus.NewInterest(interest.FormulaSimple, 2.00, interest.RatePeriodDaily),
					},
					expects: expectsError(func(err error) bool {
						if err == nil {
							return false
						}
						errMsg := err.Error()

						return strings.Contains(errMsg, "interest") && strings.Contains(errMsg, "outstanding")
					}),
				},
			},
		},
		{
			name: "Same time, full return",
			steps: []createTransferStep{
				{
					input: createTransferStepInput{
						direction: models4debtus.TransferDirectionUser2Counterparty,
						amount:    11,
						time:      dayHour(1, 1),
					},
					expects: expects(11, 0, true, nil, nil),
				},
				{
					input: createTransferStepInput{
						direction: models4debtus.TransferDirectionCounterparty2User,
						amount:    11,
						time:      dayHour(1, 2),
					},
					expects: expects(0, 11, false,
						models4debtus.TransferReturns{
							{Amount: 11},
						},
						[]transferExpectedState{
							{stepIndex: 0, isOutstanding: false, amountReturned: 11},
						},
					),
				},
			},
		},
		{
			name: "No interest: 2_gives, 1 under-return, 1 over-return",
			steps: []createTransferStep{
				{
					input: createTransferStepInput{
						direction: models4debtus.TransferDirectionUser2Counterparty,
						amount:    10,
						time:      dayHour(1, 1),
					},
					expects: expects(10, 0, true, nil, nil),
				},
				{
					input: createTransferStepInput{
						direction: models4debtus.TransferDirectionUser2Counterparty,
						amount:    15,
						time:      dayHour(1, 2),
					},
					expects: expects(25, 0, true, nil,
						[]transferExpectedState{
							{isOutstanding: true, amountReturned: 0},
						}),
				},
				{
					input: createTransferStepInput{
						isReturn:  false,
						direction: models4debtus.TransferDirectionCounterparty2User,
						amount:    7,
						time:      dayHour(1, 3),
					},
					expects: expects(18, 7, false,
						models4debtus.TransferReturns{
							{Amount: 7},
						},
						[]transferExpectedState{
							{stepIndex: 0, amountReturned: 7, isOutstanding: true},
						}),
				},
				{
					input: createTransferStepInput{
						isReturn:  false,
						direction: models4debtus.TransferDirectionCounterparty2User,
						amount:    30,
						time:      dayHour(1, 4),
					},
					expects: expects(-12, 18, true,
						models4debtus.TransferReturns{
							{Amount: 3},
							{Amount: 15},
						},
						[]transferExpectedState{
							{stepIndex: 0, amountReturned: 10, isOutstanding: false},
							{stepIndex: 1, amountReturned: 15, isOutstanding: false},
						}),
				},
			},
		},
	}

	for i, testCase := range testCases {
		if testCase.name == "" {
			t.Fatalf("Test case #%v has no name", i+1)
		}
		if len(testCase.steps) == 0 {
			t.Fatalf("Test case #%v has no steps", i+1)
		}
		t.Run(testCase.name, func(t *testing.T) {
			testCreateTransfer(t, testCase)
		})
	}
}

func testCreateTransfer(t *testing.T, testCase createTransferTestCase) {
	t.Skip("TODO: fix")
	c := context.TODO()
	dtmocks.SetupMocks(c)
	assert := assertHelper{t: t}
	currency := money.CurrencyEUR

	source := dtdal.NewTransferSourceBot(telegram.PlatformID, "test-bot", "444")

	const (
		userID    = "u1"
		contactID = "c2"
		spaceID   = "s3"
	)

	creatorUser := dbo4userus.NewUserEntry(userID)

	debtusSpace := models4debtus.NewDebtusSpaceEntry(spaceID)

	tUser := &models4debtus.TransferCounterpartyInfo{
		UserID: userID,
	}

	tContact := &models4debtus.TransferCounterpartyInfo{
		ContactID: contactID,
	}

	var (
		i    int
		step createTransferStep
	)

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Paniced at step #%v: %v\nStack: %v", i+1, r, string(debug.Stack()))
		}
	}()

	for i, step = range testCase.steps {
		var from, to *models4debtus.TransferCounterpartyInfo
		switch step.input.direction {
		case models4debtus.TransferDirectionUser2Counterparty:
			from = tUser
			to = tContact
		case models4debtus.TransferDirectionCounterparty2User:
			from = tContact
			to = tUser
		}

		request := CreateTransferRequest{
			SpaceRequest:       dto4spaceus.SpaceRequest{},
			IsReturn:           step.input.isReturn,
			ReturnToTransferID: step.input.returnToTransferID,
			Amount:             money.NewAmount(currency, step.input.amount),
			DueOn:              &step.input.time,
			Interest:           &step.input.TransferInterest,
		}
		newTransfer := NewTransferInput(
			strongoapp.LocalHostEnv,
			source,
			creatorUser,
			request,
			from, to,
		)

		// =============================================================
		output, err := Transfers.CreateTransfer(c, newTransfer)
		// =============================================================
		if output, err = assert.OutputIsNilIfErr(output, err); err != nil {
			if step.expects.isExpectedError(err) {
				// t.Logf("step #%v: got expected error: %v", i+1, err)
				continue
			}
			t.Error(err)
			return
		}

		// Save transfer ContactID for checking transfer state in later steps
		testCase.steps[i].createdTransferID = output.Transfer.ID

		var (
			contact models4debtus.DebtusSpaceContactEntry
		)
		switch step.input.direction {
		case models4debtus.TransferDirectionUser2Counterparty:
			creatorUser = output.From.User
			contact = output.To.DebtusContact
		case models4debtus.TransferDirectionCounterparty2User:
			creatorUser = output.To.User
			contact = output.From.DebtusContact
		}
		if output.Transfer.Data.IsOutstanding != step.expects.isOutstanding {
			t.Errorf("step #%v: Expected transfer.IsOutstanding does not match actual: expected:%v != got:%v",
				i+1, step.expects.isOutstanding, output.Transfer.Data.IsOutstanding)
		}
		if balance := debtusSpace.Data.Balance[currency]; balance != step.expects.balance {
			t.Errorf("step #%v: Expected user balance does not match actual: expected:%v != got:%v",
				i+1, step.expects.balance, balance)
		}
		userContact := debtusSpace.Data.Contacts[contact.ID]
		if balance := userContact.Balance[currency]; balance != step.expects.balance {
			t.Errorf("step #%v: Expected userContact balance does not match actual: expected:%v != got:%v",
				i+1, step.expects.balance, balance)
		}
		if balance := contact.Data.Balance[currency]; balance != step.expects.balance {
			t.Errorf("step #%v: Expected Contact balance does not match actual: expected:%v != got:%v",
				i+1, step.expects.balance, balance)
		}
		{ // verify balance counts
			var expectedBalanceCount int
			if step.expects.balance != 0 {
				expectedBalanceCount = 1
			}

			if balanceCount := len(userContact.Balance); balanceCount != expectedBalanceCount {
				t.Errorf("step #%v: Expected userContact len(balance) does not match actual: expected:%v != got:%v",
					i+1, expectedBalanceCount, balanceCount)
			}
			if balanceCount := len(debtusSpace.Data.Balance); balanceCount != expectedBalanceCount {
				t.Errorf("step #%v: Expected creatorUser len(balance) does not match actual: expected:%v != got:%v",
					i+1, expectedBalanceCount, balanceCount)
			}
		}

		var dbContact models4debtus.DebtusSpaceContactEntry
		if dbContact, err = GetDebtusSpaceContactByID(c, nil, spaceID, contact.ID); err != nil {
			t.Errorf("step #%v: %v", i+1, err)
			break
		}
		if balance := dbContact.Data.Balance[currency]; balance != step.expects.balance {
			t.Errorf("step #%v: Expected Contact balance does not match actual dbContact: expected:%v != got:%v",
				i+1, step.expects.balance, balance)
		}
		if output.Transfer.Data.AmountInCentsReturned != step.expects.amountReturned {
			t.Errorf("step #%v: Expected transfer.AmountInCentsReturned does not match actual: expected:%v != got:%v",
				i+1, step.expects.amountReturned, output.Transfer.Data.AmountInCentsReturned)
		}
		if output.Transfer.Data.ReturnsCount != len(step.expects.returns) {
			t.Errorf("step #%v: Expected transfer.ReturnsCount does not match actual: expected:%v != got:%v",
				i+1, len(step.expects.transfers), output.Transfer.Data.ReturnsCount)
		}
		{ // Verify returns and previous api4transfers state
			returns := output.Transfer.Data.GetReturns()
			if len(returns) != output.Transfer.Data.ReturnsCount {
				t.Errorf("step #%v: transfer.ReturnsCount is not equal to len(transfer.GetReturns()): %v != %v",
					i+1, output.Transfer.Data.ReturnsCount, len(returns))
				continue
			}
			for j, r := range returns {
				isPreviousTransfer := false
				for k := 0; k < i; k++ {
					if r.TransferID == testCase.steps[k].createdTransferID {
						isPreviousTransfer = true
						break
					}
				}
				if !isPreviousTransfer {
					t.Errorf("step #%v: transfer.Returns[%v] references unknown transfer with ContactID=%v",
						i+1, j, r.TransferID)
				}
				if r.Amount != step.expects.returns[j].Amount {
					t.Errorf("step #%v: Expected transfer.Returns[%v] does not match actual: expected:%v != got:%v",
						i+1, j, step.expects.returns[j].Amount, r.Amount)
					break
				}
			}
			for j, expectedTransfer := range step.expects.transfers {
				var previousTransfer models4debtus.TransferEntry
				if previousTransfer, err = Transfers.GetTransferByID(c, nil, testCase.steps[j].createdTransferID); err != nil {
					t.Fatalf("step #%v: %v",
						i+1, err)
					break
				}
				if previousTransfer.Data.IsOutstanding != expectedTransfer.isOutstanding {
					t.Errorf("step #%v: previous transfer from step #%v is expected to have IsOutstanding=%v but got %v",
						i+1, j+1, expectedTransfer.isOutstanding, previousTransfer.Data.IsOutstanding)
					break
				}
				if previousTransfer.Data.AmountInCentsReturned != expectedTransfer.amountReturned {
					t.Errorf("step #%v: previous transfer from step #%v is expected to have AmountInCentsReturned=%v but got %v:\nTransfer: %+v",
						i+1, j+1, expectedTransfer.amountReturned, previousTransfer.Data.AmountInCentsReturned, previousTransfer)
					break
				}
			}
		}
		if t.Failed() {
			break
		}
	}
}
