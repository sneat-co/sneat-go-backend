package dtb_settings

import (
	"fmt"
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/sneat-co/debtstracker-translations/trans"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/debtusbots/profiles/debtusbot/cmd/dtb_general"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/gae_app/debtstracker/dtdal"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/models4debtus"
	"strings"

	"errors"
	"github.com/sneat-co/debtstracker-translations/emoji"
)

//var reEmail = regexp.MustCompile("^.+@.+\\.\\w+$")

func handleInviteOnStart(whc botsfw.WebhookContext, inviteCode string, invite models4debtus.Invite) (m botsfw.MessageFromBot, err error) {
	claimAndReply := func() {
		if err = dtdal.Invite.ClaimInvite2(whc.Context(), inviteCode, invite, whc.AppUserID(), whc.BotPlatform().ID(), whc.GetBotCode()); err != nil {
			err = fmt.Errorf("failed to ClaimInvite(): %w", err)
			return
		}
		m = whc.NewMessageByCode(trans.MESSAGE_TEXT_WELCOME_ONBOARDING_INVITE_ACCEPTED)
		dtb_general.SetMainMenuKeyboard(whc, &m)
	}
	if invite.Data.Related == INVITE_IS_RELATED_TO_ONBOARDING {
		if invite.Data.CreatedByUserID != whc.AppUserID() {
			return m, errors.New("invite.Related == INVITE_IS_RELATED_TO_ONBOARDING && invite.CreatedByUserID != whc.AppUserID()")
		}
		claimAndReply()
		return
	} else {
		if invite.Data.CreatedByUserID == whc.AppUserID() {
			m = whc.NewMessage(whc.Translate(trans.MESSAGE_TEXT_ATTEMPT_TO_USE_OWN_INVITE))
			dtb_general.SetMainMenuKeyboard(whc, &m)
			return m, nil
		}

		if invite.Data.Related == "" {
			claimAndReply()
			return
		} else {
			//switch invite.RelatedTo() {
			//default:
			m = whc.NewMessage(fmt.Sprintf("Unknown invite.Related: %v", invite.Data.Related))
			//}
		}
	}
	return
}

//func OnboardingTellAboutInviteCodeAction(whc botsfw.WebhookContext) (botsfw.MessageFromBot, error) {
//	logus.Infof(whc.Context(), "onboardingTellAboutInviteCodeAction")
//	m := whc.NewMessageByCode(trans.MESSAGE_TEXT_ONBOARDING_TELL_ABOUT_INVITES, whc.GetSender().GetFirstName())
//	keyboard := tgbotapi.NewReplyKeyboardUsingStrings([][]string{
//		{whc.CommandText(trans.COMMAND_TEXT_I_HAVE_INVITE, emoji.CLOSED_LOCK_WITH_KEY)},
//		{whc.CommandText(trans.COMMAND_TEXT_SEND_ME_NEW_INVITE, emoji.PACKAGE_ICON)},
//		{trans.ChooseLocaleIcon},
//	})
//	keyboard.ResizeKeyboard = true
//	m.Keyboard = keyboard
//	m.Format = botsfw.MessageFormatHTML
//	whc.ChatData().SetAwaitingReplyTo(TELL_ABOUT_INVITE_CODE_COMMAND)
//	return m, nil
//}

//const TELL_ABOUT_INVITE_CODE_COMMAND = "tell-about-invite-code"
//
//var OnboardingTellAboutInviteCodeCommand = botsfw.Command{
//	Code: TELL_ABOUT_INVITE_CODE_COMMAND,
//	Replies: []botsfw.Command{
//		OnboardingAskInviteCodeCommand,
//		OnboardingAskInviteChannelCommand,
//		shared_all.OnboardingAskLocaleCommand,
//	},
//	Action: OnboardingTellAboutInviteCodeAction,
//}
//
//const ASK_INVITE_CHANNEL_COMMAND = "ask-invite-channel"
//
//var OnboardingAskInviteChannelCommand = botsfw.Command{
//	Code:  ASK_INVITE_CHANNEL_COMMAND,
//	Title: trans.COMMAND_TEXT_SEND_ME_NEW_INVITE,
//	Icon:  emoji.PACKAGE_ICON,
//	Replies: []botsfw.Command{
//		OnboardingAskEmailCommand,
//		OnboardingAskPhoneCommand,
//	},
//	Action: func(whc botsfw.WebhookContext) (botsfw.MessageFromBot, error) {
//		input := whc.Input()
//		switch input.(type) {
//		case botsfw.WebhookContactMessage:
//			return onboardingProcessPhoneContact(whc, input.(botsfw.WebhookContactMessage))
//		default: //case botsfw.WebhookTextMessage:
//			m := whc.NewMessageByCode(trans.MESSAGE_TEXT_ASK_INVITE_CHANNEL)
//			telegramKeyboard := &tgbotapi.ReplyKeyboardMarkup{
//				ResizeKeyboard:  true,
//				OneTimeKeyboard: true,
//				Keyboard: [][]tgbotapi.KeyboardButton{
//					{
//						{
//							Text:           SmsChannelCommand.CommandText(whc),
//							RequestContact: true,
//						},
//					},
//					{
//						{
//							Text: EmailChannelCommand.CommandText(whc),
//						},
//					},
//					{
//						{
//							Text: whc.CommandText(trans.COMMAND_TEXT_I_HAVE_INVITE, emoji.CLOSED_LOCK_WITH_KEY),
//						},
//					},
//				},
//			}
//			m.Keyboard = telegramKeyboard
//			chatEntity := whc.ChatData()
//			chatEntity.PushStepToAwaitingReplyTo(ASK_INVITE_CHANNEL_COMMAND)
//			return m, nil
//		}
//	},
//}
//
//const ASK_INVITE_CODE_COMMAND = "ask-invite-code"
//
//var OnboardingAskInviteCodeCommand = botsfw.Command{
//	Code:       ASK_INVITE_CODE_COMMAND,
//	Title:      trans.COMMAND_TEXT_I_HAVE_INVITE,
//	Icon:       emoji.CLOSED_LOCK_WITH_KEY,
//	ExactMatch: trans.COMMAND_TEXT_I_HAVE_INVITE,
//	Replies: []botsfw.Command{
//		OnboardingCheckInviteCommand,
//	},
//	Action: func(whc botsfw.WebhookContext) (botsfw.MessageFromBot, error) {
//		if whc.Input().(botsfw.WebhookTextMessage).Text() == whc.CommandText(trans.COMMAND_TEXT_I_HAVE_INVITE, emoji.CLOSED_LOCK_WITH_KEY) {
//			chatEntity := whc.ChatData()
//			chatEntity.PopStepsFromAwaitingReplyUpToSpecificParent(TELL_ABOUT_INVITE_CODE_COMMAND)
//			chatEntity.PushStepToAwaitingReplyTo(ASK_INVITE_CODE_COMMAND)
//			m := whc.NewMessageByCode(trans.MESSAGE_TEXT_PLEASE_ENTER_INVITE_CODE)
//			m.Keyboard = &tgbotapi.ReplyKeyboardHide{HideKeyboard: true}
//			return m, nil
//		} else {
//			return OnboardingCheckInviteCommand.Action(whc)
//		}
//	},
//}
//
//const CHECK_INVITE_COMMAND = "check-invite"

func NewMistypedCommand(messageToAdd string) botsfw.Command {
	var message []string
	if messageToAdd == "" {
		message = []string{trans.MESSAGE_TEXT_OK_PLEASE_TRY_AGAIN}
	} else {
		message = []string{trans.MESSAGE_TEXT_OK_PLEASE_TRY_AGAIN, messageToAdd}
	}
	return TextCommand(
		trans.COMMAND_TEXT_MISTYPE_WILL_TRY_AGAIN,
		message,
		emoji.SORRY_ICON, emoji.THUMB_UP_ICON, true)
}

//var OnboardingCheckInviteCommand = botsfw.Command{
//	Code: CHECK_INVITE_COMMAND,
//	Replies: []botsfw.Command{
//		OnboardingAskEmailCommand,
//		OnboardingAskPhoneCommand,
//		NewMistypedCommand(trans.MESSAGE_TEXT_PLEASE_ENTER_INVITE_CODE),
//	},
//	Action: func(whc botsfw.WebhookContext) (m botsfw.MessageFromBot, err error) {
//		// Check for code
//
//		c := whc.Context()
//		chatEntity := whc.ChatData()
//		chatEntity.PushStepToAwaitingReplyTo(CHECK_INVITE_COMMAND)
//		inviteCode := strings.ToUpper(whc.Input().(botsfw.WebhookTextMessage).Text())
//		userID := whc.AppUserID()
//
//		if err = dtdal.Invite.ClaimInvite(c, userID, inviteCode, whc.BotPlatform().ContactID(), whc.GetBotCode()); err != nil {
//			if dal.IsNotFound(err) {
//				m = whc.NewMessage(emoji.NO_ENTRY_SIGN_ICON + " " + strings.TrimSpace(fmt.Sprintf(whc.Translate(trans.MESSAGE_TEXT_WRONG_INVITE_CODE), inviteCode)))
//				m.Keyboard = tgbotapi.NewReplyKeyboardUsingStrings([][]string{
//					{NewMistypedCommand("").DefaultTitle(whc)},
//					{EmailChannelCommand.CommandText(whc)},
//					{SmsChannelCommand.CommandText(whc)},
//				})
//				return m, nil
//			}
//			return
//		}
//		if err = botsfw.SetAccessGranted(whc, true); err != nil {
//			err = errors.Wrap(err, "Failed to call botsfw.SetAccessGranted(whc, true)")
//			return
//		}
//
//		m, err = dtb_general.MainMenuCommand.Action(whc)
//
//		return
//	},
//}

func TextCommand(on string, message []string, icon, replyIcon string, hideKeyboard bool) botsfw.Command {
	return botsfw.Command{
		Code:       fmt.Sprintf("TextCommand(%v)", on),
		Title:      on,
		Icon:       icon,
		ExactMatch: on,
		Action: func(whc botsfw.WebhookContext) (botsfw.MessageFromBot, error) {
			for i, untranslated := range message {
				message[i] = whc.Translate(untranslated)
			}
			messageText := strings.Join(message, " ")

			if replyIcon != "" {
				messageText = replyIcon + " " + messageText
			}
			m := whc.NewMessage(messageText)
			m.Keyboard = &tgbotapi.ReplyKeyboardHide{HideKeyboard: true}
			return m, nil
		},
	}
}

//var onboardingCommands map[string]Command = map[string]Command{
//	TELL_ABOUT_INVITE_CODE_COMMAND: OnboardingTellAboutInviteCode,
//}

//func askContactDetailsCommand(cmd ContactsChannelCommand, invalidMessageCode string, altCmd ContactsChannelCommand, telegramKeyboard func(botsfw.WebhookContext) botsfw.Keyboard) botsfw.Command {
//	return botsfw.Command{
//		Code:  cmd.code,
//		Title: cmd.title,
//		Icon:  cmd.icon,
//		//exactMatch: cmd.exactMatch, //TODO: Fix it!
//		Replies: []botsfw.Command{
//			NewMistypedCommand(cmd.message),
//		},
//		Action: func(whc botsfw.WebhookContext) (m botsfw.MessageFromBot, err error) {
//			c := whc.Context()
//
//			logus.Infof(c, "Command(code=%v).Action()", cmd.code)
//			chatEntity := whc.ChatData()
//			awaitingReplyTo := chatEntity.GetAwaitingReplyTo()
//
//			input := whc.Input()
//			switch input.(type) {
//			case botsfw.WebhookContactMessage:
//				return onboardingProcessPhoneContact(whc, input.(botsfw.WebhookContactMessage))
//			default: // case botsfw.WebhookTextMessage:
//				message := whc.Input().(botsfw.WebhookTextMessage)
//				messageText := message.Text()
//				altOption := messageText == altCmd.title
//				//logus.Debugf(c, "code: %v, awaitingReplyTo: %v, messageText: %v, altOption=%v, altCommand: %v", cmd.code, awaitingReplyTo, messageText, altOption, altCmd.title)
//				switch {
//				case messageText == whc.CommandText(trans.COMMAND_TEXT_TELL_ME_MORE_ABOUT_INVITES, emoji.QUESTION_ICON):
//					switch cmd.code {
//					case SmsChannelCommand.code:
//						return OnboardingTellAboutInviteCodeAction(whc)
//					default:
//						return m, fmt.Errorf("Unhandled DebtusSpaceContactEntry message by %v command", cmd.code)
//					}
//
//				case strings.HasSuffix(awaitingReplyTo, cmd.code) && !altOption:
//					switch cmd.code {
//					case SmsChannelCommand.code:
//						if messageText != "" {
//							panic("Not implemented yet")
//						}
//					case EmailChannelCommand.code:
//						return onboardingProcessEmail(messageText, whc, altCmd)
//					}
//				default:
//					if altOption {
//						logus.Debugf(c, "Switching to alt code: %v", altCmd.code)
//						cmd.code = altCmd.code
//						cmd.message = altCmd.message
//					}
//					m = whc.NewMessage(whc.Translate(cmd.message))
//					m.Keyboard = telegramKeyboard(whc)
//					chatEntity.PopStepsFromAwaitingReplyUpToSpecificParent(ASK_INVITE_CHANNEL_COMMAND)
//					chatEntity.PushStepToAwaitingReplyTo(cmd.code)
//				}
//			}
//			return m, nil
//		},
//	}
//}

type ContactsChannelCommand struct {
	//code    string
	icon  string
	title string
	//message string
}

func (c ContactsChannelCommand) CommandText(whc botsfw.WebhookContext) string {
	return whc.CommandText(c.title, c.icon)
}

//var EmailChannelCommand = ContactsChannelCommand{
//	code:    "onboarding-ask-email",
//	icon:    emoji.EMAIL_ICON,
//	title:   trans.COMMAND_TEXT_SEND_ME_NEW_INVITE_BY_EMAIL,
//	message: trans.MESSAGE_TEXT_PLEASE_PROVIDE_YOUR_EMAIL,
//}
//
//var SmsChannelCommand = ContactsChannelCommand{
//	code:    "onboarding-ask-phone",
//	icon:    emoji.PHONE_ICON,
//	title:   trans.COMMAND_TEXT_SEND_ME_NEW_INVITE_BY_SMS,
//	message: trans.MESSAGE_TEXT_PLEASE_PROVIDE_YOUR_PHONE_NUMBER,
//}

//var OnboardingAskPhoneCommand = askContactDetailsCommand(SmsChannelCommand, trans.MESSAGE_TEXT_WRONG_PHONE_NUMBER, EmailChannelCommand,
//	func(whc botsfw.WebhookContext) botsfw.Keyboard {
//		return &tgbotapi.ReplyKeyboardMarkup{
//			ResizeKeyboard:  true,
//			OneTimeKeyboard: true,
//			Keyboard: [][]tgbotapi.KeyboardButton{
//				{
//					{
//						Text:           whc.Translate(trans.COMMAND_TEXT_SEND_MY_PHONE_NUMBER),
//						RequestContact: true,
//					},
//				},
//			},
//		}
//	},
//)
//
//var OnboardingAskEmailCommand = askContactDetailsCommand(
//	EmailChannelCommand,
//	trans.MESSAGE_TEXT_WRONG_EMAIL, SmsChannelCommand,
//	func(whc botsfw.WebhookContext) botsfw.Keyboard {
//		return tgbotapi.NewHideKeyboard(false)
//	},
//)

const ON_USER_CONTACT_RECEIVED_COMMAND = "on-user-contact-received"

var OnboardingOnUserContactReceivedCommand = botsfw.Command{
	Code: ON_USER_CONTACT_RECEIVED_COMMAND,
	Action: func(whc botsfw.WebhookContext) (m botsfw.MessageFromBot, err error) {
		return
	},
}

//func onboardingProcessPhoneContact(whc botsfw.WebhookContext, contact botsfw.WebhookContactMessage) (m botsfw.MessageFromBot, err error) {
//	c := whc.Context()
//	//whc.ChatData().SetAwaitingReplyTo(ON_USER_CONTACT_RECEIVED_COMMAND)
//	invite, err := dtdal.Invite.CreatePersonalInvite(whc.ExecutionContext(), whc.AppUserID(), models.InviteBySms, contact.PhoneNumber(), whc.BotPlatform().ContactID(), whc.GetBotCode(), INVITE_IS_RELATED_TO_ONBOARDING)
//	if err != nil {
//		return m, err
//	}
//
//	utmParams := anybot.UtmParams{
//		Source:   anybot.UtmSourceFromContext(whc),
//		Medium:   "sms",
//		Campaign: anybot.UTM_CAMPAIGN_ONBOARDING_INVITE,
//	}
//	templateParams := invites.InviteTemplateParams{
//		ToName:     contact.FirstName(),
//		FromName:   "DebtsTracker",
//		InviteCode: invite.ContactID,
//		TgBot:      whc.GetBotCode(),
//		Utm:        utmParams.String(),
//	}
//
//	smsText, err := anybot.TextTemplates.RenderTemplate(c, whc, trans.SMS_INVITE_TEXT, templateParams)
//	if err != nil {
//		return m, err
//	}
//	isTestSender, smsResponse, twilioException, err := sms.SendSms(whc.Context(), whc.GetBotSettings().Env == strongoapp.EnvProduction, contact.PhoneNumber(), smsText)
//	if err != nil {
//		return m, err
//	}
//	if twilioException != nil {
//		sms.TwilioExceptionToMessage(whc, twilioException)
//	}
//	if err != nil {
//		return m, err
//	}
//	m = whc.NewMessage(whc.Translate(trans.MESSAGE_TEXT_USER_CONTACT_FOR_INVITE_RECEIVED))
//
//	if isTestSender {
//		m.Text += "\n\n<b>SMS text</b>\n" + smsResponse.Body
//	}
//
//	m.Keyboard = &tgbotapi.InlineKeyboardMarkup{
//		InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{
//			{
//				{
//					Text: "Shares on FB to get instant invite",
//					URL:  "https://apps.facebook.com/debtstracker/",
//				},
//			},
//		},
//	}
//	return m, err
//}

const EMAIL_CONFIRMATION_SENT_COMMAND = "email-confirmation-sent"

const INVITE_IS_RELATED_TO_ONBOARDING = "onboarding=yes"

//func onboardingProcessEmail(messageText string, whc botsfw.WebhookContext, altCmd ContactsChannelCommand) (m botsfw.MessageFromBot, err error) {
//
//	logus.Infof(whc.Context(), "onboardingProcessEmail(messageText=%v)", messageText)
//	email := strings.TrimSpace(messageText)
//
//	if !reEmail.MatchString(email) {
//		m = whc.NewMessageByCode(trans.MESSAGE_TEXT_WRONG_EMAIL)
//		m.Keyboard = tgbotapi.NewReplyKeyboardUsingStrings([][]string{
//			{whc.CommandText(trans.COMMAND_TEXT_MISTYPE_WILL_TRY_AGAIN, emoji.SORRY_ICON)}, //TODO: Reuse command from replies:
//			{altCmd.CommandText(whc)},
//			{whc.CommandText(trans.COMMAND_TEXT_TELL_ME_MORE_ABOUT_INVITES, emoji.QUESTION_ICON)},
//		})
//	} else {
//		//TODO: Try to send email and handle return codes & exceptions
//		invite, err := dtdal.Invite.CreatePersonalInvite(whc.ExecutionContext(), whc.AppUserID(), models.InviteByEmail, email, whc.BotPlatform().ContactID(), whc.GetBotCode(), INVITE_IS_RELATED_TO_ONBOARDING)
//		if err != nil {
//			return m, err
//		}
//		_, err = invites.SendInviteByEmail(whc.ExecutionContext(), "DebtsTracker", email, whc.GetSender().GetFirstName(), invite.ContactID, whc.GetBotCode(), anybot.UtmSourceFromContext(whc))
//		if err != nil {
//			return m, err
//		}
//		whc.ChatData().SetAwaitingReplyTo(EMAIL_CONFIRMATION_SENT_COMMAND)
//		mt := fmt.Sprintf(whc.Translate(trans.MESSAGE_TEXT_USER_EMAIL_FOR_INVITE_RECEIVED), email)
//		if whc.BotPlatform().ContactID() == telegram.PlatformID {
//			mt += "\n\n" + whc.Translate(trans.MESSAGE_TEXT_USER_EMAIL_FOR_INVITE_SENT_TELEGRAM)
//		}
//		m = whc.NewMessage(mt)
//		m.Keyboard = &tgbotapi.ReplyKeyboardMarkup{
//			ResizeKeyboard:  true,
//			OneTimeKeyboard: true,
//			Keyboard: [][]tgbotapi.KeyboardButton{
//				{
//					{Text: whc.Translate(trans.COMMAND_TEXT_I_HAVE_NOT_GOT_EMAIL)},
//				},
//				{
//					{Text: whc.CommandText(trans.COMMAND_TEXT_I_HAVE_INVITE, emoji.CLOSED_LOCK_WITH_KEY)},
//				},
//			},
//		}
//	}
//
//	return m, nil
//}
