package dtb_transfer

import (
	"fmt"
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw-store/botsfwmodels"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/crediterra/money"
	"github.com/sneat-co/debtstracker-translations/emoji"
	"github.com/sneat-co/debtstracker-translations/trans"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/common4debtus"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/models4debtus"
	"github.com/strongo/logus"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var BorrowingWizardCompletedCommand = TransferWizardCompletedCommand("transfer-to-completed")
var LendingWizardCompletedCommand = TransferWizardCompletedCommand("transfer-from-completed")

func CreateStartTransferWizardCommand(code, messageText string, commands []string, askTransferAmountCommand botsfw.Command) botsfw.Command {
	return botsfw.Command{
		Code:     code,
		Commands: commands,
		Replies:  []botsfw.Command{askTransferAmountCommand},
		Matcher: func(c botsfw.Command, whc botsfw.WebhookContext) bool {
			if m, ok := whc.Input().(botsfw.WebhookTextMessage); ok && IsCurrencyIcon(m.Text()) && whc.ChatData().GetAwaitingReplyTo() == code {
				return true
			}
			return false
		},
		Action: func(whc botsfw.WebhookContext) (m botsfw.MessageFromBot, err error) {
			c := whc.Context()
			logus.Debugf(c, "CreateStartTransferWizardCommand(code=%v).Action()", code)
			mt := strings.TrimSpace(whc.Input().(botsfw.WebhookTextMessage).Text())
			chatEntity := whc.ChatData()
			switch {
			case money.HasCurrencyPrefix(mt) || IsCurrencyIcon(mt):
				currency := money.CleanupCurrency(mt)
				chatEntity.AddWizardParam("currency", string(currency))
				return askTransferAmountCommand.Action(whc)
			case mt == "...":
				m = whc.NewMessage(whc.Translate(trans.MESSAGE_TEXT_NOT_IMPLEMENTED_YET) + " " + emoji.FLUSHED_FACE)
			case mt == whc.Translate(trans.COMMAND_TEXT_CANCEL):
				return cancelTransferWizardCommandAction(whc)
			default:
				isMainMenuCommand := strings.HasPrefix(mt, "/")
				if !isMainMenuCommand { // TODO: This sucks, do we need to move the asking for amount in dedicated command?
					mtLower := strings.ToLower(mt)
					for _, command := range commands {
						if strings.HasPrefix(mtLower, command) {
							isMainMenuCommand = true
							break
						}
					}
				}
				if isMainMenuCommand {
					whc.ChatData().SetAwaitingReplyTo(code)
					m = whc.NewMessageByCode(messageText)
					m.Text += "\n\n" + strings.Replace(whc.Translate(trans.MESSAGE_TEXT_CHOOSE_CURRENCY), "<a>", fmt.Sprintf(
						`<a href="%v">`,
						common4debtus.GetChooseCurrencyUrlForUser(
							whc.AppUserID(), whc.Locale(), whc.BotPlatform().ID(), whc.GetBotCode(),
							"tg-chat="+botsfwmodels.NewChatID(whc.GetBotCode(), whc.MustBotChatID()),
						),
					), 1)
					buttons := AskTransferCurrencyButtons(whc)
					keyboard := tgbotapi.NewReplyKeyboardUsingStrings(buttons)
					keyboard.OneTimeKeyboard = true
					m.Keyboard = keyboard
				} else if _, err = strconv.ParseFloat(mt, 64); err == nil {
					return whc.NewMessageByCode(trans.MESSAGE_TEXT_CURRENCY_NAME_IS_NUMBER), nil
					// User entered a number
				} else { // TODO: Document why we allow this!?
					//err = nil // Ignore error from strconv.ParseFloat()
					if strings.ToLower(mt) == "euro" {
						mt = "EUR"
					} else if len(mt) == 3 {
						currencyCode := strings.ToUpper(mt)
						if currencyCode != mt && money.CurrencyCode(mt).IsMoney() {
							mt = currencyCode
						}
					}
					currency := money.CurrencyCode(mt)
					chatEntity.AddWizardParam("currency", string(currency))
					return askTransferAmountCommand.Action(whc)
				}
			}
			return m, nil
		},
	}
}

var AskLendingAmountCommand = AskTransferAmountCommand("ask-lending-amount", trans.MESSAGE_TEXT_ASK_LENDING_AMOUNT,
	CreateAskTransferCounterpartyCommand(false, "ask-lending-counterparty", "", "", trans.MESSAGE_TEXT_ASK_LENDING_COUNTERPARTY,
		[]botsfw.Command{TransferAskDueDateReturnToUser},
		NewCounterpartyCommand(TransferAskDueDateReturnToUser),
		_transferAskDueDate(TransferAskDueDateReturnToUser),
	),
)

var StartLendingWizardCommand = CreateStartTransferWizardCommand(
	"start-lending-wizard",
	trans.MESSAGE_TEXT_ASK_LENDING_TYPE,
	trans.Commands(trans.COMMAND_TEXT_GAVE, trans.COMMAND_GAVE, emoji.GIVE_ICON),
	AskLendingAmountCommand,
)

var AskBorrowingAmountCommand = AskTransferAmountCommand("ask-borrowing-amount", trans.MESSAGE_TEXT_ASK_BORROWING_AMOUNT,
	CreateAskTransferCounterpartyCommand(false, "ask-borrowing-counterparty", "", "", trans.MESSAGE_TEXT_ASK_BORROWING_COUNTERPARTY,
		[]botsfw.Command{TransferAskDueDateReturnByUser},
		NewCounterpartyCommand(TransferAskDueDateReturnByUser),
		_transferAskDueDate(TransferAskDueDateReturnByUser),
	),
)

var StartBorrowingWizardCommand = CreateStartTransferWizardCommand(
	"start-borrowing-wizard",
	trans.MESSAGE_TEXT_ASK_BORROWING_TYPE,
	trans.Commands(trans.COMMAND_TEXT_GOT, trans.COMMAND_GOT, emoji.TAKE_ICON),
	AskBorrowingAmountCommand,
)

const SET_DUE_DATE_COMMAND = "set-due-date"

var SetDueDateCommand = botsfw.Command{
	Code: SET_DUE_DATE_COMMAND,
	Action: func(whc botsfw.WebhookContext) (botsfw.MessageFromBot, error) {
		whc.ChatData().SetAwaitingReplyTo("")
		m := whc.NewMessage("Due date to be saved")
		return m, nil
	},
}

const ASK_DUE_DATE_COMMAND = "ask-due-date"

//var AskDueDateCommand = botsfw.Command{
//	Code:    ASK_DUE_DATE_COMMAND,
//	Replies: []botsfw.Command{SetDueDateCommand},
//	Action: func(whc botsfw.WebhookContext) (botsfw.MessageFromBot, error) {
//		m := whc.NewMessage(`<strong>When is the due date?</strong>
//
//Recognized inputs:
//	<code>in N days|weeks|months</code>
//   	<i>Example: in 10 days, in 1 month</i>
//	<code>on DD.MM.YYYY</code>
//		<i>Example: on 30.12.2016</i>
//	<code>on MM/DD/YYYY</code>
//		<i>Example: on 12/30/2016</i>
//	<code>on DD Month</code>
//		<i>Example: on 12 March</i>`)
//		chatEntity := whc.ChatData()
//		awaitingReplyTo := chatEntity.GetAwaitingReplyTo()
//		transferID := strings.Split(awaitingReplyTo, ":")[1]
//		chatEntity.SetAwaitingReplyTo(fmt.Sprintf("asked-for-deadline:transferID=%v", transferID))
//		m.Format = botsfw.MessageFormatHTML
//		m.IsReplyToInputMessage = true
//		m.Keyboard = tgbotapi.ForceReply{ForceReply: true}
//		return m, nil
//	},
//}

func _transferAskDueDate(c botsfw.Command) _onContactSelectedAction {
	return func(whc botsfw.WebhookContext, counterparty models4debtus.DebtusSpaceContactEntry) (m botsfw.MessageFromBot, err error) {
		return c.Action(whc)
	}
}

var reDate = regexp.MustCompile(`(\d{1,2})(\.|/)(\d{1,2})(\.|/)(\d+)`)

func TransferAskDueDateCommand(code string, nextCommand botsfw.Command) botsfw.Command {
	return botsfw.Command{
		Code: code,
		Replies: []botsfw.Command{
			nextCommand,
		},
		Action: func(whc botsfw.WebhookContext) (botsfw.MessageFromBot, error) {

			c := whc.Context()
			logus.Infof(c, "TransferAskDueDateCommand(code=%v).Action()", code)
			m := whc.NewMessageByCode(trans.MESSAGE_TEXT_ASK_DUE)
			chatEntity := whc.ChatData()
			if chatEntity.IsAwaitingReplyTo(code) {
				mt := strings.TrimSpace(whc.Input().(botsfw.WebhookTextMessage).Text())
				logus.Debugf(c, "Chat is awating reply to %v", code)
				var duration time.Duration
				switch mt {
				case whc.Translate(trans.COMMAND_TEXT_IN_FEW_MINUTES):
					duration, _ = time.ParseDuration("1m")
				case whc.Translate(trans.COMMAND_TEXT_TOMORROW):
					duration, _ = time.ParseDuration("24h")
				case whc.Translate(trans.COMMAND_TEXT_DAY_AFTER_TOMORROW):
					duration, _ = time.ParseDuration("48h")
				case whc.Translate(trans.COMMAND_TEXT_IN_1_WEEK):
					duration, _ = time.ParseDuration("168h")
				case whc.Translate(trans.COMMAND_TEXT_IN_1_MONTH):
					duration, _ = time.ParseDuration("720h")
				case whc.Translate(trans.COMMAND_TEXT_IT_CAN_BE_RETURNED_ANYTIME):
					return nextCommand.Action(whc)
				default:
					if strings.Contains(mt, emoji.CALENDAR_ICON) {
						m = whc.NewMessageByCode(trans.MESSAGE_TEXT_ASK_DUE_DATE)
						return m, nil
					} else {
						m, date, err := processSetDate(whc)
						if !date.IsZero() {
							chatEntity.AddWizardParam("due", date.Format(TRANSFER_WIZARD_DUE_DATE_FORMAT))
							return nextCommand.Action(whc)
						}
						return m, err
					}
				}
				if duration > 0 {
					chatEntity.AddWizardParam("due", duration.String())
				}
				return nextCommand.Action(whc)
			} else {
				logus.Debugf(c, "Chat is NOT awating reply to %v", code)
				chatEntity.PushStepToAwaitingReplyTo(code)
				keyboard := tgbotapi.NewReplyKeyboardUsingStrings([][]string{
					{
						whc.Translate(trans.COMMAND_TEXT_IN_FEW_MINUTES),
						whc.Translate(trans.COMMAND_TEXT_IT_CAN_BE_RETURNED_ANYTIME),
					},
					{
						whc.Translate(trans.COMMAND_TEXT_TOMORROW),
						whc.Translate(trans.COMMAND_TEXT_DAY_AFTER_TOMORROW),
					},
					{
						whc.Translate(trans.COMMAND_TEXT_IN_1_WEEK),
						emoji.CALENDAR_ICON + " " + whc.Translate(trans.COMMAND_TEXT_SET_DATE),
					},
				})
				keyboard.OneTimeKeyboard = true
				m.Keyboard = keyboard
				return m, nil
			}
		},
	}
}

var TransferAskDueDateReturnByUser = TransferAskDueDateCommand(
	"ask-due-date-return-by-user",
	TransferFromUserAskNoteOrCommentCommand,
)
var TransferAskDueDateReturnToUser = TransferAskDueDateCommand(
	"ask-due-date-return-to-user",
	TransferToUserAskNoteOrCommentCommand,
)

func processSetDate(whc botsfw.WebhookContext) (m botsfw.MessageFromBot, date time.Time, err error) {
	c := whc.Context()

	mt := strings.TrimSpace(whc.Input().(botsfw.WebhookTextMessage).Text())
	if match := reDate.FindStringSubmatch(mt); len(match) > 0 {
		if match[2] != match[4] {
			m = whc.NewMessageByCode(trans.MESSAGE_TEXT_INVALID_DATE)
			return m, date, nil
		}
		var dateFormat string
		yearStr := match[5]
		switch len(yearStr) {
		case 4:
			dateFormat = "02.01.2006"
		case 2:
			dateFormat = "02.01.06"
		default:
			m = whc.NewMessageByCode(trans.MESSAGE_TEXT_INVALID_YEAR)
			return m, date, nil
		}
		dayStr := match[1]
		monthStr := match[3]
		if len(dayStr) == 1 {
			dayStr = "0" + dayStr
		}
		if len(monthStr) == 1 {
			monthStr = "0" + monthStr
		}
		var day, month int
		if day, err = strconv.Atoi(dayStr); err != nil || day > 31 {
			m = whc.NewMessageByCode(trans.MESSAGE_TEXT_INVALID_DAY)
		}
		if month, err = strconv.Atoi(monthStr); err != nil || month > 31 {
			m = whc.NewMessageByCode(trans.MESSAGE_TEXT_INVALID_MONTH)
		}
		if month > 12 && day < 12 {
			temp := month
			month = day
			day = temp
		}
		dateToParse := fmt.Sprintf("%02d.%02d", day, month) + "." + yearStr
		logus.Debugf(c, "dateToParse: %v", dateFormat)
		if date, err = time.Parse(dateFormat, dateToParse); err != nil {
			m = whc.NewMessageByCode(trans.MESSAGE_TEXT_WRONG_DATE)
		}
	} else {
		logus.Debugf(c, "Regex not matched")
		m = whc.NewMessageByCode(trans.MESSAGE_TEXT_INVALID_DATE)
	}
	return m, date, nil
}
