package dtb_invite

import (
	"fmt"
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/sneat-co/debtstracker-translations/trans"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/common4debtus"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/debtusbots/profiles/debtusbot/cmd/dtb_general"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/gae_app/debtstracker/dtdal"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/gae_app/invites"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/models4debtus"
	"github.com/strongo/logus"
	"net/url"
	"strings"

	"github.com/sneat-co/debtstracker-translations/emoji"
)

var AskInviteAddressTelegramCommand = AskInviteAddress(string(models4debtus.InviteByTelegram), emoji.ROCKET_ICON, trans.COMMAND_TEXT_INVITE_BY_TELEGRAM, trans.MESSAGE_TEXT_INVITE_BY_TELEGRAM, trans.MESSAGE_TEXT_NO_CONTACT_RECEIVED)
var AskInviteAddressEmailCommand = AskInviteAddress(string(models4debtus.InviteByEmail), emoji.EMAIL_ICON, trans.COMMAND_TEXT_SEND_BY_EMAIL, trans.MESSAGE_TEXT_INVITE_BY_EMAIL, trans.MESSAGE_TEXT_INVALID_EMAIL)
var AskInviteAddressSmsCommand = AskInviteAddress(string(models4debtus.InviteBySms), emoji.PHONE_ICON, trans.COMMAND_TEXT_SEND_BY_SMS, trans.MESSAGE_TEXT_INVITE_BY_SMS, trans.MESSAGE_TEXT_INVALID_PHONE_NUMBER)

func AskInviteAddress(channel, icon, commandText, messageCode, invalidMessageCode string) botsfw.Command {
	code := fmt.Sprintf("ask-%v-address-for-invite", channel)
	return botsfw.Command{
		Code:  code,
		Icon:  icon,
		Title: commandText,
		Action: func(whc botsfw.WebhookContext) (m botsfw.MessageFromBot, err error) {
			chatEntity := whc.ChatData()

			if chatEntity.IsAwaitingReplyTo(code) {
				email := strings.TrimSpace(whc.Input().(botsfw.WebhookTextMessage).Text())
				isValid := channel == string(models4debtus.InviteByEmail) && strings.Contains(email, "@") && strings.Contains(email, ".")
				if isValid {
					invite, err := dtdal.Invite.CreatePersonalInvite(whc, whc.AppUserID(), models4debtus.InviteByEmail, email, whc.BotPlatform().ID(), whc.GetBotCode(), "counterparty=?")
					if err != nil {
						logus.Errorf(whc.Context(), "Failed to call invites.CreateInvite()")
						return m, err
					}
					var emailID string
					emailID, err = invites.SendInviteByEmail(
						whc.ExecutionContext(),
						whc,
						whc.GetSender().GetFirstName(),
						"alex@debtusbot.io",
						"Stranger",
						invite.ID,
						whc.GetBotCode(),
						common4debtus.UtmSourceFromContext(whc),
					)
					if err != nil {
						return m, err
					}
					m = whc.NewMessageByCode(trans.MESSAGE_TEXT_INVITE_CREATED, emailID)
				} else {
					m = whc.NewMessageByCode(invalidMessageCode)
					m.Keyboard = tgbotapi.NewReplyKeyboardUsingStrings([][]string{
						{whc.Translate(trans.COMMAND_TEXT_MISTYPE_WILL_TRY_AGAIN)},
						{whc.Translate(trans.COMMAND_TEXT_OTHER_WAYS_TO_SEND_INVITE)},
						{dtb_general.MainMenuCommand.DefaultTitle(whc)},
					})
				}
			} else {
				m = whc.NewMessageByCode(messageCode)
				chatEntity.PushStepToAwaitingReplyTo(code)
			}
			return m, nil
		},
	}
}

var AskInviteAddressCallbackCommand = botsfw.Command{
	Code: "invite",
	CallbackAction: func(whc botsfw.WebhookContext, callbackUrl *url.URL) (m botsfw.MessageFromBot, err error) {
		q := callbackUrl.Query()
		echoSelection := func(mt string) error {
			if m, err = whc.NewEditMessage(whc.Translate(trans.MESSAGE_TEXT_ABOUT_INVITES)+"\n\n"+mt, botsfw.MessageFormatHTML); err != nil {
				return err
			}
			_, err := whc.Responder().SendMessage(whc.Context(), m, botsfw.BotAPISendMessageOverHTTPS)
			return fmt.Errorf("failed to edit callback message: %w", err)
		}
		_ = whc.ChatData() // To switch locale
		switch q.Get("by") {
		case string(models4debtus.InviteByEmail):
			if err = echoSelection(whc.Translate(trans.MESSAGE_TEXT_YOU_SELECTED_INVITE_BY_EMAIL)); err != nil {
				return
			}
			return AskInviteAddressEmailCommand.Action(whc)
		case string(models4debtus.InviteBySms):
			if err = echoSelection(whc.Translate(trans.MESSAGE_TEXT_YOU_SELECTED_INVITE_BY_SMS)); err != nil {
				return
			}
			return AskInviteAddressSmsCommand.Action(whc)
		case "":
			logus.Warningf(whc.Context(), "AskInviteAddressCallbackCommand: got request to create invite without specifying a channel - not implemented yet. Need to ask a channel first. Check how it works if message forwarded to secret chat.")
			m.Text = whc.Translate(trans.MESSAGE_TEXT_NOT_IMPLEMENTED_YET)
			return
		default:
			err = fmt.Errorf("unknown invite channel: %v", q.Get("by"))
			return
		}
	},
}
