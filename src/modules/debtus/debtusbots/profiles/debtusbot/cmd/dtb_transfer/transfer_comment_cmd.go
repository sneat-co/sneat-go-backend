package dtb_transfer

import (
	//"github.com/sneat-co/debtusbot-translations/emoji"
	//"fmt"
	"fmt"
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw/botinput"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/sneat-co/debtstracker-translations/trans"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/facade4debtus"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/dal4userus"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/dbo4userus"
	"github.com/strongo/logus"
	"net/url"
	"strings"

	"errors"
	"golang.org/x/net/html"
)

const (
	//TRANSFER_WIZARD_PARAM_NOTE    = "note"
	TRANSFER_WIZARD_PARAM_COMMENT = "comment"
)

//const (
//	ADD_NOTE_COMMAND    = "add-note"
//	ADD_COMMENT_COMMAND = "add-comment"
//)
//
//func createTransferAddNoteOrCommentCommand(code string, anotherCommand *botsfw.Command, nextCommand botsfw.Command) botsfw.Command {
//	var icon, title string
//	switch code {
//	case ADD_NOTE_COMMAND:
//		icon = emoji.MEMO_ICON
//		title = trans.COMMAND_TEXT_ADD_NOTE_TO_TRANSFER
//	case ADD_COMMENT_COMMAND:
//		icon = emoji.NEWSPAPER_ICON
//		title = trans.COMMAND_TEXT_ADD_COMMENT_TO_TRANSFER
//	}
//
//	return botsfw.Command{
//		Code:  code,
//		Icon:  icon,
//		Title: title,
//		Action: func(whc botsfw.WebhookContext) (m botsfw.MessageFromBot, err error) {
//
//			logus.Debugf(c, "createTransferAddNoteOrCommentCommand().Action(), code=%v", code)
//			if code != ADD_NOTE_COMMAND && code != ADD_COMMENT_COMMAND {
//				panic(fmt.Sprintf("Unknown code: %v", code))
//			}
//			chatEntity := whc.ChatData()
//			if chatEntity.IsAwaitingReplyTo(code) {
//				switch code {
//				case ADD_NOTE_COMMAND:
//					chatEntity.AddWizardParam(TRANSFER_WIZARD_PARAM_NOTE, whc.Input().(botsfw.WebhookTextMessage).Text())
//					if chatEntity.GetWizardParam(TRANSFER_WIZARD_PARAM_COMMENT) != "" {
//						return nextCommand.Action(whc)
//					}
//					m = whc.NewMessage(whc.Translate(trans.MESSAGE_TEXT_TRANSFER_NOTE_ADDED_ASK_FOR_COMMENT))
//					m.Keyboard = tgbotapi.NewReplyKeyboard(
//						[]tgbotapi.KeyboardButton{tgbotapi.NewKeyboardButton(anotherCommand.DefaultTitle(whc))},
//						[]tgbotapi.KeyboardButton{tgbotapi.NewKeyboardButton(whc.Translate(trans.COMMAND_TEXT_NO_COMMENT_FOR_TRANSFER))},
//					)
//				case ADD_COMMENT_COMMAND:
//					chatEntity.AddWizardParam(TRANSFER_WIZARD_PARAM_COMMENT, whc.Input().(botsfw.WebhookTextMessage).Text())
//					if chatEntity.GetWizardParam(TRANSFER_WIZARD_PARAM_NOTE) != "" {
//						return nextCommand.Action(whc)
//					}
//					m = whc.NewMessage(whc.Translate(trans.MESSAGE_TEXT_TRANSFER_COMMENT_ADDED_ASK_FOR_NOTE))
//					m.Keyboard = tgbotapi.NewReplyKeyboard(
//						[]tgbotapi.KeyboardButton{tgbotapi.NewKeyboardButton(anotherCommand.DefaultTitle(whc))},
//						[]tgbotapi.KeyboardButton{tgbotapi.NewKeyboardButton(whc.Translate(trans.COMMAND_TEXT_NO_NOTE_FOR_TRANSFER))},
//					)
//				default:
//					panic(fmt.Sprintf("Unknown code: %v", code))
//				}
//				chatEntity.PopStepsFromAwaitingReplyUpToSpecificParent(ASK_NOTE_OR_COMMENT_FOR_TRANSFER_COMMAND)
//			} else {
//				chatEntity.PushStepToAwaitingReplyTo(code)
//				switch code {
//				case ADD_NOTE_COMMAND:
//					m = whc.NewMessage(whc.Translate(trans.MESSAGE_TEXT_TRANSFER_ASK_FOR_NOTE))
//				case ADD_COMMENT_COMMAND:
//					m = whc.NewMessage(whc.Translate(trans.MESSAGE_TEXT_TRANSFER_ASK_FOR_COMMENT))
//				default:
//					panic(fmt.Sprintf("Unknown code: %v", code))
//				}
//				m.Keyboard = tgbotapi.NewHideKeyboard(true)
//			}
//			m.Format = botsfw.MessageFormatHTML
//			return m, err
//		},
//	}
//}

func createTransferAskNoteOrCommentCommand(code string, nextCommand botsfw.Command) botsfw.Command {
	var addNoteCommand botsfw.Command
	var addCommentCommand botsfw.Command

	//addNoteCommand = createTransferAddNoteOrCommentCommand(ADD_NOTE_COMMAND, &addCommentCommand, nextCommand)
	//addCommentCommand = createTransferAddNoteOrCommentCommand(ADD_COMMENT_COMMAND, &addNoteCommand, nextCommand)

	return botsfw.Command{
		Code: code,
		Replies: []botsfw.Command{
			addNoteCommand,
			addCommentCommand,
		},
		Action: func(whc botsfw.WebhookContext) (m botsfw.MessageFromBot, err error) {
			ctx := whc.Context()
			logus.Infof(ctx, "createTransferAskNoteOrCommentCommand().Action()")
			chatEntity := whc.ChatData()
			//noOptionSelected := false
			if chatEntity.IsAwaitingReplyTo(code) {
				if m, err = interestAction(whc, nextCommand.Action); err != nil || m.Text != "" {
					return
				}
				mt := whc.Input().(botinput.WebhookTextMessage).Text()
				switch mt {
				//case whc.Translate(trans.COMMAND_TEXT_ADD_NOTE_TO_TRANSFER):
				//	return addNoteCommand.Action(whc)
				//case whc.Translate(trans.COMMAND_TEXT_ADD_COMMENT_TO_TRANSFER):
				//	return addCommentCommand.Action(whc)
				//case whc.Translate(trans.COMMAND_TEXT_NO_COMMENT_OR_NOTE_FOR_TRANSFER):
				//	return nextCommand.Action(whc)
				case whc.Translate(trans.COMMAND_TEXT_NO_COMMENT_FOR_TRANSFER):
					return nextCommand.Action(whc)
					//case whc.Translate(trans.COMMAND_TEXT_NO_NOTE_FOR_TRANSFER):
					//	return nextCommand.Action(whc)
				default:
					chatEntity.AddWizardParam(TRANSFER_WIZARD_PARAM_COMMENT, mt)
					return nextCommand.Action(whc)
					//noOptionSelected = true
				}
			} else {
				chatEntity.PushStepToAwaitingReplyTo(code)
			}

			m = whc.NewMessage(whc.Translate(trans.MESSAGE_TEXT_TRANSFER_ASK_FOR_INTEREST_SHORT))
			m.Format = botsfw.MessageFormatHTML
			m.Keyboard = tgbotapi.NewInlineKeyboardMarkup(
				[]tgbotapi.InlineKeyboardButton{
					tgbotapi.NewInlineKeyboardButtonData(whc.Translate(trans.COMMAND_TEXT_MORE_ABOUT_INTEREST_COMMAND), ASK_FOR_INTEREST_AND_COMMENT_COMMAND),
				},
			)
			if _, err = whc.Responder().SendMessage(ctx, m, botsfw.BotAPISendMessageOverHTTPS); err != nil {
				return
			}

			user := dbo4userus.NewUserEntry(whc.AppUserID())
			if err = dal4userus.GetUser(ctx, nil, user); err != nil {
				return
			}
			spaceID := user.Data.GetFamilySpaceID()
			var transferWizard TransferWizard
			if transferWizard, err = NewTransferWizard(whc); err != nil {
				return
			}
			counterpartyID := transferWizard.CounterpartyID(ctx)
			if counterpartyID == "" {
				return m, errors.New("transferWizard.CounterpartyID() == 0")
			}
			counterparty, err := facade4debtus.GetDebtusSpaceContactByID(whc.Context(), nil, spaceID, counterpartyID)
			m.Text = strings.TrimLeft(fmt.Sprintf("%v\n(<i>%v</i>)",
				whc.Translate(trans.MESSAGE_TEXT_TRANSFER_ASK_FOR_COMMENT_ONLY),
				whc.Translate(trans.MESSAGE_TEXT_VISIBLE_TO_YOU_AND_COUNTERPARTY, html.EscapeString(counterparty.Data.FullName()))),
				"\n ",
			)

			replyKeyboard := tgbotapi.NewReplyKeyboard([]tgbotapi.KeyboardButton{
				{Text: whc.Translate(trans.COMMAND_TEXT_NO_COMMENT_FOR_TRANSFER)},
			})
			replyKeyboard.OneTimeKeyboard = true
			m.Keyboard = replyKeyboard
			m.Format = botsfw.MessageFormatHTML
			return
		},
	}
}

const ASK_FOR_INTEREST_AND_COMMENT_COMMAND = "ask-for-interest-and-comment-long"

var AskForInterestAndCommentCallbackCommand = botsfw.Command{
	Code: ASK_FOR_INTEREST_AND_COMMENT_COMMAND,
	CallbackAction: func(whc botsfw.WebhookContext, callbackUrl *url.URL) (m botsfw.MessageFromBot, err error) {
		m.Text = whc.Translate(trans.MESSAGE_TEXT_TRANSFER_ASK_FOR_INTEREST_LONG)
		m.Format = botsfw.MessageFormatHTML
		m.IsEdit = true
		return
	},
}

const ASK_NOTE_OR_COMMENT_FOR_TRANSFER_COMMAND = "ask-note-or-comment"

var TransferFromUserAskNoteOrCommentCommand = createTransferAskNoteOrCommentCommand(
	ASK_NOTE_OR_COMMENT_FOR_TRANSFER_COMMAND,
	BorrowingWizardCompletedCommand,
)

var TransferToUserAskNoteOrCommentCommand = createTransferAskNoteOrCommentCommand(
	ASK_NOTE_OR_COMMENT_FOR_TRANSFER_COMMAND,
	LendingWizardCompletedCommand,
)
