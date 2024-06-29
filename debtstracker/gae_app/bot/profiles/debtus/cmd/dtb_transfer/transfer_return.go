package dtb_transfer

import (
	"fmt"
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw-store/botsfwmodels"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/crediterra/money"
	"github.com/sneat-co/debtstracker-translations/trans"
	"github.com/strongo/logus"
	"net/url"
	"strings"
	"time"

	"errors"
	"github.com/sneat-co/debtstracker-translations/emoji"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/dtdal"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/facade"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/models"
	"github.com/strongo/decimal"
	"golang.org/x/net/html"
)

//var StartReturnWizardCommand = botsfw.Command{
//	Code: "start-return-wizard",
//	Action: func(whc botsfw.WebhookContext) (m botsfw.MessageFromBot, err error) {
//	},
//}

const RETURN_WIZARD_COMMAND = "return-wizard"

var StartReturnWizardCommand = botsfw.Command{
	Code:     RETURN_WIZARD_COMMAND,
	Commands: trans.Commands(trans.COMMAND_RETURNED),
	Replies:  []botsfw.Command{AskReturnCounterpartyCommand, AskToChooseDebtToReturnCommand},
	Action: func(whc botsfw.WebhookContext) (m botsfw.MessageFromBot, err error) {
		logus.Debugf(whc.Context(), "StartReturnWizardCommand.Action()")
		whc.ChatData().SetAwaitingReplyTo(RETURN_WIZARD_COMMAND)
		return AskReturnCounterpartyCommand.Action(whc)
	},
}

func askIfReturnedInFull(whc botsfw.WebhookContext, counterparty models.ContactEntry, currency money.CurrencyCode, value decimal.Decimal64p2) (m botsfw.MessageFromBot, err error) {
	amount := money.Amount{Currency: money.CurrencyCode(currency), Value: value}
	var mt string
	switch {
	case value < 0:
		mt = trans.MESSAGE_TEXT_YOU_OWE_TO_COUNTERPARTY_SINGLE_DEBT
	case value > 0:
		mt = trans.MESSAGE_TEXT_COUNTERPARTY_OWES_YOU_SINGLE_DEBT
	case value == 0:
		errorMessage := fmt.Sprintf("ERROR: Balance for currency [%v] is: %v", currency, value)
		logus.Warningf(whc.Context(), errorMessage)
		m = whc.NewMessage(errorMessage)
		return
	}
	chatEntity := whc.ChatData()
	chatEntity.PushStepToAwaitingReplyTo(ASK_IF_RETURNED_IN_FULL_COMMAND)
	chatEntity.AddWizardParam("currency", string(currency))
	amount.Value = amount.Value.Abs()
	m = whc.NewMessage(fmt.Sprintf(
		whc.Translate(mt), html.EscapeString(counterparty.Data.FullName()), amount) +
		"\n\n" + whc.Translate(trans.MESSAGE_TEXT_IS_IT_RETURNED_IN_FULL))
	m.Format = botsfw.MessageFormatHTML
	m.Keyboard = tgbotapi.NewReplyKeyboardUsingStrings(
		[][]string{
			{whc.Translate(trans.BUTTON_TEXT_DEBT_RETURNED_FULLY)},
			{whc.Translate(trans.BUTTON_TEXT_DEBT_RETURNED_PARTIALLY)},
			{whc.Translate(trans.COMMAND_TEXT_CANCEL)},
		},
	)
	return
}

const ASK_RETURN_COUNTERPARTY_COMMAND = "ask-return-counterparty"

var AskReturnCounterpartyCommand = CreateAskTransferCounterpartyCommand(
	true,
	ASK_RETURN_COUNTERPARTY_COMMAND,
	trans.COMMAND_TEXT_RETURN,
	emoji.RETURN_BACK_ICON,
	trans.MESSAGE_TEXT_RETURN_ASK_TO_CHOOSE_COUNTERPARTY,
	[]botsfw.Command{
		AskToChooseDebtToReturnCommand,
		AskIfReturnedInFullCommand,
	},
	botsfw.Command{}, //newContactCommand - We do not allow to create a new contact on return
	func(whc botsfw.WebhookContext, counterparty models.ContactEntry) (m botsfw.MessageFromBot, err error) {
		c := whc.Context()

		logus.Debugf(c, "StartReturnWizardCommand.onCounterpartySelectedAction(counterparty.ID=%v)", counterparty.ID)
		var balanceWithInterest money.Balance
		balanceWithInterest, err = counterparty.Data.BalanceWithInterest(c, time.Now())
		if err != nil {
			err = fmt.Errorf("failed to get counterparty balance with interest: %w", err)
			return
		}
		//TODO: Display MESSAGE_TEXT_COUNTERPARTY_OWES_YOU_SINGLE_DEBT or MESSAGE_TEXT_YOU_OWE_TO_COUNTERPARTY_SINGLE_DEBT
		switch len(balanceWithInterest) {
		case 1:
			for currency, value := range balanceWithInterest {
				return askIfReturnedInFull(whc, counterparty, currency, value)
			}
		case 0:
			errorMessage := whc.Translate(trans.MESSAGE_TEXT_COUNTERPARTY_HAS_EMPTY_BALANCE, counterparty.Data.FullName())
			logus.Debugf(c, "Balance is empty: "+errorMessage)
			m = whc.NewMessage(errorMessage)
		default:
			buttons := make([][]string, len(balanceWithInterest)+1)
			var i int
			buttons[0] = []string{whc.Translate(trans.COMMAND_TEXT_CANCEL)}
			for currency, value := range balanceWithInterest {
				i += 1
				buttons[i] = []string{_debtAmountButtonText(whc, currency, value, counterparty)}
			}
			m = askToChooseDebt(whc, buttons)
		}
		return
	},
)

func askToChooseDebt(whc botsfw.WebhookContext, buttons [][]string) (m botsfw.MessageFromBot) {
	if len(buttons) > 0 {
		whc.ChatData().PushStepToAwaitingReplyTo(ASK_TO_CHOOSE_DEBT_TO_RETURN_COMMAND)
		m = whc.NewMessage(whc.Translate(trans.MESSAGE_TEXT_CHOOSE_DEBT_THAT_HAS_BEEN_RETURNED))
		m.Keyboard = tgbotapi.NewReplyKeyboardUsingStrings(buttons)
	} else {
		m = whc.NewMessage(whc.Translate(trans.MESSAGE_TEXT_NO_DEBTS_TO_RETURN))
	}
	return m
}

func _debtAmountButtonText(whc botsfw.WebhookContext, currency money.CurrencyCode, value decimal.Decimal64p2, counterparty models.ContactEntry) string {
	amount := money.Amount{Currency: currency, Value: value.Abs()}
	var mt string
	switch {
	case value > 0:
		mt = trans.BUTTON_TEXT_SOMEONE_OWES_TO_YOU_AMOUNT
	case value < 0:
		mt = trans.BUTTON_TEXT_YOU_OWE_AMOUNT_TO_SOMEONE
	default:
		mt = "ERROR (%v) - zero value: %v"
	}
	return fmt.Sprintf(whc.Translate(mt), counterparty.Data.FullName(), amount)
}

const ASK_IF_RETURNED_IN_FULL_COMMAND = "ask-if-return-in-full"

var AskIfReturnedInFullCommand = botsfw.Command{
	Code:    ASK_IF_RETURNED_IN_FULL_COMMAND,
	Replies: []botsfw.Command{AskHowMuchHaveBeenReturnedCommand},
	Action: func(whc botsfw.WebhookContext) (m botsfw.MessageFromBot, err error) {
		chatEntity := whc.ChatData()
		if chatEntity.IsAwaitingReplyTo(ASK_IF_RETURNED_IN_FULL_COMMAND) {
			switch whc.Input().(botsfw.WebhookTextMessage).Text() {
			case whc.Translate(trans.BUTTON_TEXT_DEBT_RETURNED_FULLY):
				m, err = processReturnCommand(whc, 0)
				//common.CreateTransfer(whc.Context(), whc.AppUserID(), )
			case whc.Translate(trans.BUTTON_TEXT_DEBT_RETURNED_PARTIALLY):
				m, err = AskHowMuchHaveBeenReturnedCommand.Action(whc)
			default:
				return TryToProcessHowMuchHasBeenReturned(whc)
			}
			return m, err

		} else {
			err = errors.New("AskIfReturnedInFullCommand: Not implemented yet")
			return m, err
		}
	},
}

func processReturnCommand(whc botsfw.WebhookContext, returnValue decimal.Decimal64p2) (m botsfw.MessageFromBot, err error) {
	if returnValue < 0 {
		panic(fmt.Sprintf("returnValue < 0: %v", returnValue))
	}
	c := whc.Context()
	chatEntity := whc.ChatData()
	var (
		counterpartyID string
		transferID     string
	)
	if counterpartyID, transferID, err = getReturnWizardParams(whc); err != nil {
		return m, err
	}
	counterparty, err := getCounterparty(whc, counterpartyID)
	if err != nil {
		return m, err
	}
	counterpartyBalanceWithInterest, err := counterparty.Data.BalanceWithInterest(c, time.Now())
	if err != nil {
		err = fmt.Errorf("failed to get balance with interest for contact %v: %v", counterparty.ID, err)
		return
	}
	awaitingUrl, err := url.Parse(chatEntity.GetAwaitingReplyTo())
	if err != nil {
		return m, err
	}
	currency := money.CurrencyCode(awaitingUrl.Query().Get("currency"))

	if transferID != "" && returnValue > 0 {
		var transfer models.TransferEntry
		if transfer, err = facade.Transfers.GetTransferByID(whc.Context(), nil, transferID); err != nil {
			return
		}

		returnAmount := money.NewAmount(currency, returnValue)
		if outstandingAmount := transfer.Data.GetOutstandingAmount(time.Now()); outstandingAmount.Value < returnValue {
			m.Text = whc.Translate(trans.MESSAGE_TEXT_RETURN_IS_TOO_BIG, returnAmount, outstandingAmount, outstandingAmount.Value)
			return
		}
	}

	if previousBalance, ok := counterpartyBalanceWithInterest[currency]; ok {
		if returnValue == 0 {
			returnValue = previousBalance.Abs()
		}
		previousBalance := money.Amount{Currency: currency, Value: previousBalance}
		direction, err := getReturnDirectionFromDebtValue(previousBalance)
		if err != nil {
			return m, err
		}
		return CreateReturnAndShowReceipt(whc, transferID, counterpartyID, direction, money.NewAmount(currency, returnValue))
	} else {
		return m, fmt.Errorf("ContactEntry has no currency in balance. counterpartyID=%v,  currency='%v'", counterpartyID, currency)
	}
}

const ASK_HOW_MUCH_HAVE_BEEN_RETURNED = "ask-how-much-have-been-returned"

var AskHowMuchHaveBeenReturnedCommand = botsfw.Command{
	Code: ASK_HOW_MUCH_HAVE_BEEN_RETURNED,
	Action: func(whc botsfw.WebhookContext) (m botsfw.MessageFromBot, err error) {
		c := whc.Context()
		logus.Debugf(c, "AskHowMuchHaveBeenReturnedCommand.Action()")
		chatEntity := whc.ChatData()
		if chatEntity.IsAwaitingReplyTo(ASK_HOW_MUCH_HAVE_BEEN_RETURNED) {
			return TryToProcessHowMuchHasBeenReturned(whc)
		} else {
			m = whc.NewMessage(whc.Translate(trans.MESSAGE_TEXT_ASK_HOW_MUCH_HAS_BEEN_RETURNED))
			m.Keyboard = tgbotapi.NewHideKeyboard(true)
			chatEntity.PushStepToAwaitingReplyTo(ASK_HOW_MUCH_HAVE_BEEN_RETURNED)
			return m, err
		}
	},
}

func TryToProcessHowMuchHasBeenReturned(whc botsfw.WebhookContext) (m botsfw.MessageFromBot, err error) {
	if amountValue, err := decimal.ParseDecimal64p2(whc.Input().(botsfw.WebhookTextMessage).Text()); err != nil {
		m = whc.NewMessage(whc.Translate(trans.MESSAGE_TEXT_INCORRECT_VALUE_NOT_A_NUMBER))
		return m, nil
	} else {
		if amountValue > 0 {
			return processReturnCommand(whc, amountValue)
		} else {
			m = whc.NewMessage(whc.Translate(trans.MESSAGE_TEXT_INCORRECT_VALUE_IS_NEGATIVE))
			return m, nil
		}
	}
}

const ASK_TO_CHOOSE_DEBT_TO_RETURN_COMMAND = "ask-to-choose-debt-to-return"

var AskToChooseDebtToReturnCommand = botsfw.Command{
	Code: ASK_TO_CHOOSE_DEBT_TO_RETURN_COMMAND,
	Replies: []botsfw.Command{
		AskIfReturnedInFullCommand,
	},
	Action: func(whc botsfw.WebhookContext) (m botsfw.MessageFromBot, err error) {
		c := whc.Context()
		counterpartyID, _, _ := getReturnWizardParams(whc)
		var (
			theCounterparty models.ContactEntry
			balance         money.Balance
		)
		if counterpartyID == "" {
			// Let's try to get counterpartyEntity from message text
			mt := whc.Input().(botsfw.WebhookTextMessage).Text()
			splittedBySeparator := strings.Split(mt, "|")
			counterpartyTitle := strings.Join(splittedBySeparator[:len(splittedBySeparator)-1], "|")
			counterpartyTitle = strings.TrimSpace(counterpartyTitle)
			chatEntity := whc.ChatData()
			//var botAppUser botsfwmodels.AppUserData
			//botAppUser, err = whc.AppUserData()
			//if err != nil {
			//	return m, err
			//}
			//user := botAppUser.(*models.DebutsAppUserDataOBSOLETE)
			var counterparties []models.ContactEntry
			if counterparties, err = dtdal.Contact.GetLatestContacts(whc, nil, 10, -1); err != nil {
				return m, err
			}
			var counterpartyFound bool
			now := time.Now()
			for _, counterpartyItem := range counterparties {
				counterpartyItemTitle := counterpartyItem.Data.FullName()
				if counterpartyItemTitle == counterpartyTitle {
					if balance, err = counterpartyItem.Data.BalanceWithInterest(c, now); err != nil {
						err = fmt.Errorf("failed to get balance with interest for contact %v: %w", counterpartyItem.ID, err)
						return
					}
					theCounterparty = counterpartyItem
					counterpartyFound = true
					chatEntity.AddWizardParam(WIZARD_PARAM_COUNTERPARTY, counterpartyItem.ID)
					break
				}
			}
			if !counterpartyFound {
				m = whc.NewMessageByCode(trans.MESSAGE_TEXT_UNKNOWN_COUNTERPARTY_ON_RETURN)
				return m, nil
			}
		} else {
			var counterparty models.ContactEntry
			if counterparty, err = getCounterparty(whc, counterpartyID); err != nil {
				return m, err
			}
			if balance, err = counterparty.Data.BalanceWithInterest(c, time.Now()); err != nil {
				err = fmt.Errorf("failed to get balance with interest for contact %v: %w", counterparty.ID, err)
				return
			}
			theCounterparty = counterparty
		}

		mt := whc.Input().(botsfw.WebhookTextMessage).Text()
		for currency, value := range balance {
			if mt == _debtAmountButtonText(whc, currency, value, theCounterparty) {
				return askIfReturnedInFull(whc, theCounterparty, currency, value)
			}
		}
		if m.Text == "" {
			m = whc.NewMessageByCode(trans.MESSAGE_TEXT_UNKNOWN_DEBT)
		}
		return m, err
	},
}

func CreateReturnAndShowReceipt(whc botsfw.WebhookContext, returnToTransferID string, counterpartyID string, direction models.TransferDirection, returnAmount money.Amount) (m botsfw.MessageFromBot, err error) {
	c := whc.Context()
	logus.Debugf(c, "CreateReturnAndShowReceipt(returnToTransferID=%s, counterpartyID=%s)", returnToTransferID, counterpartyID)

	if returnAmount.Value < 0 {
		logus.Warningf(c, "returnAmount.Value < 0: %v", returnAmount.Value)
		returnAmount.Value = returnAmount.Value.Abs()
	}

	creatorInfo := models.TransferCounterpartyInfo{
		UserID:    whc.AppUserID(),
		ContactID: counterpartyID,
	}

	if m, err = CreateTransferFromBot(whc, true, returnToTransferID, direction, creatorInfo, returnAmount, time.Time{}, models.NoInterest()); err != nil {
		return m, err
	}
	logus.Debugf(c, "createReturnAndShowReceipt(): %v", m)
	return m, err
}

func getReturnDirectionFromDebtValue(currentDebt money.Amount) (models.TransferDirection, error) {
	switch {
	case currentDebt.Value < 0:
		return models.TransferDirectionUser2Counterparty, nil
	case currentDebt.Value > 0:
		return models.TransferDirectionCounterparty2User, nil
	}
	return models.TransferDirection(""), fmt.Errorf("Zero value for currency: [%v]", currentDebt.Currency)
}

func getReturnWizardParams(whc botsfw.WebhookContext) (counterpartyID string, transferID string, err error) {
	awaitingReplyTo := whc.ChatData().GetAwaitingReplyTo()
	params, err := url.ParseQuery(botsfwmodels.AwaitingReplyToQuery(awaitingReplyTo))
	if err != nil {
		return counterpartyID, transferID, fmt.Errorf("failed in AwaitingReplyToQuery(): %w", err)
	}
	counterpartyID = params.Get(WIZARD_PARAM_COUNTERPARTY)
	transferID = params.Get(WIZARD_PARAM_TRANSFER)
	return
}

func getCounterparty(whc botsfw.WebhookContext, counterpartyID string) (counterparty models.ContactEntry, err error) {
	//counterparty = new(models.ContactEntry)
	if counterparty, err = facade.GetContactByID(whc.Context(), nil, counterpartyID); err != nil {
		return
	}
	return
}
