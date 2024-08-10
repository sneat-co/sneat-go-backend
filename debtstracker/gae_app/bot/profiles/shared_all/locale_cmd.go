package shared_all

import (
	"fmt"
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw-store/botsfwmodels"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/debtstracker-translations/trans"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/i18n"
	"github.com/strongo/logus"
	"net/url"
	"strings"

	"bytes"
	"context"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/bot/profiles/debtus/cmd/dtb_general"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/facade2debtus"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/models"
)

const (
	SettingsLocaleListCallbackPath = "settings/locale/list"
	SettingsLocaleSetCallbackPath  = "settings/locale/set"
)

const onboardingAskLocaleCommandCode = "onboarding-ask-locale"

var localesReplyKeyboard = tgbotapi.NewReplyKeyboard(
	[]tgbotapi.KeyboardButton{
		{Text: i18n.LocaleEnUS.TitleWithIcon()},
		{Text: i18n.LocaleRuRu.TitleWithIcon()},
	},
	[]tgbotapi.KeyboardButton{
		{Text: i18n.LocaleEsEs.TitleWithIcon()},
		{Text: i18n.LocaleItIt.TitleWithIcon()},
	},
	[]tgbotapi.KeyboardButton{
		{Text: i18n.LocaleDeDe.TitleWithIcon()},
		{Text: i18n.LocaleFaIr.TitleWithIcon()},
	},
)

func createOnboardingAskLocaleCommand(botParams BotParams) botsfw.Command {
	return botsfw.Command{
		Code:       onboardingAskLocaleCommandCode,
		ExactMatch: trans.ChooseLocaleIcon,
		Action: func(whc botsfw.WebhookContext) (m botsfw.MessageFromBot, err error) {
			return onboardingAskLocaleAction(whc, "", botParams)
		},
	}
}

func onboardingAskLocaleAction(whc botsfw.WebhookContext, messagePrefix string, botParams BotParams) (m botsfw.MessageFromBot, err error) {
	chatEntity := whc.ChatData()

	if chatEntity.IsAwaitingReplyTo(onboardingAskLocaleCommandCode) {
		messageText := whc.Input().(botsfw.WebhookTextMessage).Text()
		for _, locale := range trans.SupportedLocales {
			if locale.TitleWithIcon() == messageText {
				return setPreferredLanguageAction(whc, locale.Code5, "onboarding", botParams)
			}
		}
		m = whc.NewMessageByCode(trans.MESSAGE_TEXT_UNKNOWN_LANGUAGE)
		//localesReplyKeyboard.OneTimeKeyboard = true
		m.Keyboard = localesReplyKeyboard
	} else {
		m.Text = messagePrefix + m.Text
		chatEntity.SetAwaitingReplyTo(onboardingAskLocaleCommandCode)
		m = whc.NewMessageByCode(trans.MESSAGE_TEXT_ONBOARDING_ASK_TO_CHOOSE_LANGUAGE, whc.GetSender().GetFirstName())
		//localesReplyKeyboard.OneTimeKeyboard = true
		m.Keyboard = localesReplyKeyboard
	}
	return
}

var askPreferredLocaleFromSettingsCallback = botsfw.Command{
	Code: SettingsLocaleListCallbackPath,
	CallbackAction: func(whc botsfw.WebhookContext, _ *url.URL) (m botsfw.MessageFromBot, err error) {
		callbackData := fmt.Sprintf("%v?mode=settings&code5=", SettingsLocaleSetCallbackPath)
		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			[]tgbotapi.InlineKeyboardButton{
				{Text: i18n.LocaleEnUS.TitleWithIcon(), CallbackData: callbackData + i18n.LocaleEnUS.Code5},
				{Text: i18n.LocaleRuRu.TitleWithIcon(), CallbackData: callbackData + i18n.LocaleRuRu.Code5},
			},
			[]tgbotapi.InlineKeyboardButton{
				{Text: i18n.LocaleEsEs.TitleWithIcon(), CallbackData: callbackData + i18n.LocaleEsEs.Code5},
				{Text: i18n.LocaleItIt.TitleWithIcon(), CallbackData: callbackData + i18n.LocaleItIt.Code5},
			},
			[]tgbotapi.InlineKeyboardButton{
				{Text: i18n.LocaleDeDe.TitleWithIcon(), CallbackData: callbackData + i18n.LocaleDeDe.Code5},
				{Text: i18n.LocaleFaIr.TitleWithIcon(), CallbackData: callbackData + i18n.LocaleFaIr.Code5},
			},
		) //dtb_general.LanguageOptions(whc, false)
		logus.Debugf(whc.Context(), "AskPreferredLanguage(): locale: %v", whc.Locale().Code5)
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, []tgbotapi.InlineKeyboardButton{
			{Text: whc.Translate(trans.COMMAND_TEXT_SETTING), CallbackData: SettingsCommandCode},
		})
		if m, err = whc.NewEditMessage(whc.Translate(trans.MESSAGE_TEXT_CHOOSE_UI_LANGUAGE), botsfw.MessageFormatHTML); err != nil {
			return
		}
		m.Keyboard = keyboard
		return m, err
	},
}

func setLocaleCallbackCommand(botParams BotParams) botsfw.Command {
	return botsfw.Command{
		Code: SettingsLocaleSetCallbackPath,
		CallbackAction: func(whc botsfw.WebhookContext, callbackUrl *url.URL) (m botsfw.MessageFromBot, err error) {
			return setPreferredLanguageAction(whc, callbackUrl.Query().Get("code5"), callbackUrl.Query().Get("mode"), botParams)
		},
	}
}

func setPreferredLanguageAction(whc botsfw.WebhookContext, code5, mode string, botParams BotParams) (m botsfw.MessageFromBot, err error) {
	c := whc.Context()
	logus.Debugf(c, "setPreferredLanguageAction(code5=%v, mode=%v)", code5, mode)

	var appUserData botsfwmodels.AppUserData

	if appUserData, err = whc.AppUserData(); err != nil {
		logus.Errorf(c, ": %v", err)
		return m, fmt.Errorf("%w: failed to load userEntity", err)
	}

	preferredLocale := appUserData.BotsFwAdapter().GetPreferredLocale()

	var (
		localeChanged  bool
		selectedLocale i18n.Locale
	)

	chatData := whc.ChatData()
	logus.Debugf(c, "userEntity.PreferredLanguage: %v, chatData.GetPreferredLanguage(): %v, code5: %v", preferredLocale, chatData.GetPreferredLanguage(), code5)
	if preferredLocale != code5 || chatData.GetPreferredLanguage() != code5 {
		logus.Debugf(c, "PreferredLanguage will be updated for userEntity & chat entities.")
		for _, locale := range trans.SupportedLocalesByCode5 {
			if locale.Code5 == code5 {
				_ = whc.SetLocale(locale.Code5)

				if err = facade.RunReadwriteTransaction(c, func(c context.Context, tx dal.ReadwriteTransaction) (err error) {
					var user models.AppUser
					if user, err = facade2debtus.User.GetUserByID(c, tx, whc.AppUserID()); err != nil {
						return
					}
					if err = user.Data.SetPreferredLocale(locale.Code5); err != nil {
						return fmt.Errorf("%w: failed to set preferred locale for user", err)
					}
					chatData.SetPreferredLanguage(locale.Code5)
					chatData.SetAwaitingReplyTo("")
					//chatKey := botsfwmodels.NewChatKey(whc.GetBotCode(), whc.MustBotChatID())
					if err = whc.SaveBotChat(c); err != nil {
						return
					}
					return facade2debtus.User.SaveUser(c, tx, user)
				}); err != nil {
					return
				}
				localeChanged = true
				selectedLocale = locale
				if whc.GetBotSettings().Env == "prod" {
					ga := whc.GA()
					gaEvent := ga.GaEventWithLabel("settings", "locale-changed", strings.ToLower(locale.Code5))
					if gaErr := ga.Queue(gaEvent); gaErr != nil {
						logus.Warningf(c, "Failed to log event: %v", gaErr)
					} else {
						logus.Infof(c, "GA event queued: %v", gaEvent)
					}
				}
				break
			}
		}
		if !localeChanged {
			logus.Errorf(c, "Unknown locale: %v", code5)
		}
	} else {
		selectedLocale = i18n.GetLocaleByCode5(chatData.GetPreferredLanguage())
	}
	//if localeChanged {

	switch mode {
	case "onboarding":
		logus.Debugf(c, "whc.Locale().Code5: %v", whc.Locale().Code5)
		m = whc.NewMessageByCode(trans.MESSAGE_TEXT_YOUR_SELECTED_PREFERRED_LANGUAGE, selectedLocale.NativeTitle)
		botParams.SetMainMenu(whc, &m)
		if _, err = whc.Responder().SendMessage(c, m, botsfw.BotAPISendMessageOverHTTPS); err != nil {
			logus.Errorf(c, "Failed to notify userEntity about selected language: %v", err)
			// Not critical, lets continue
		}
		return aboutDrawAction(whc, nil)
	case "settings":
		if localeChanged {
			if m, err = dtb_general.MainMenuAction(whc, whc.Translate(trans.MESSAGE_TEXT_LOCALE_CHANGED, selectedLocale.TitleWithIcon()), false); err != nil {
				return
			}
			if _, err = whc.Responder().SendMessage(c, m, botsfw.BotAPISendMessageOverHTTPS); err != nil {
				return m, err
			}
			return SettingsMainAction(whc)
			//if _, err = whc.Responder().SendMessage(c, m, botsfw.BotAPISendMessageOverHTTPS); err != nil {
			//	return m, err
			//}
			//return dtb_general.MainMenuAction(whc, )
		} else {
			return SettingsMainAction(whc)
		}
	default:
		panic(fmt.Sprintf("Unknown mode: %v", mode))
	}
}

const (
	moreAboutDrawCommandCode = "more-about-draw"
	joinDrawCommandCode      = "join-draw"
)

var aboutDrawCommand = botsfw.Command{
	Commands: []string{"/draw"},
	Code:     moreAboutDrawCommandCode,
	Action: func(whc botsfw.WebhookContext) (m botsfw.MessageFromBot, err error) {
		return aboutDrawAction(whc, nil)
	},
	CallbackAction: aboutDrawAction,
}

var joinDrawCommand = botsfw.Command{
	Code:           joinDrawCommandCode,
	CallbackAction: aboutDrawAction,
}

func aboutDrawAction(whc botsfw.WebhookContext, callbackUrl *url.URL) (m botsfw.MessageFromBot, err error) {
	c := whc.Context()
	buf := new(bytes.Buffer)
	sender := whc.GetSender()
	name := sender.GetFirstName()
	if name == "" {
		name = sender.GetUserName()
		if name == "" {
			name = sender.GetLastName()
		}
	}
	buf.WriteString(whc.Translate(trans.MESSAGE_TEXT_ABOUT_DRAW_SHORT, name))
	buf.WriteString("\n\n")
	m.Format = botsfw.MessageFormatHTML
	if callbackUrl == nil {
		buf.WriteString(whc.Translate(trans.MESSAGE_TEXT_ABOUT_DRAW_CALL_TO_ACTION))
		m.Text = buf.String()
		m.Keyboard = tgbotapi.NewInlineKeyboardMarkup(
			[]tgbotapi.InlineKeyboardButton{
				{
					Text:         whc.Translate(trans.COMMAN_TEXT_MORE_ABOUT_DRAW),
					CallbackData: moreAboutDrawCommandCode,
				},
			},
		)
		return
	} else {
		m.IsEdit = true
		buf.WriteString(whc.Translate(trans.MESSAGE_TEXT_ABOUT_DRAW_MORE))
		m.Text = buf.String()
		switch callbackUrl.Path {
		case moreAboutDrawCommandCode:
			m.Keyboard = tgbotapi.NewInlineKeyboardMarkup(
				[]tgbotapi.InlineKeyboardButton{
					{
						Text:         whc.Translate(trans.COMMAN_TEXT_I_AM_IN_DRAW),
						CallbackData: joinDrawCommandCode,
					},
				},
			)
			return
		case joinDrawCommandCode:
			if _, err = whc.Responder().SendMessage(c, m, botsfw.BotAPISendMessageOverHTTPS); err != nil {
				logus.Warningf(c, "Failed to edit message: %v", err)
				err = nil // Not critical
			}
			m.IsEdit = false
			m.Text = whc.Translate(trans.MESSAGE_TEXT_JOINED_DRAW)
			return
		default:
			err = fmt.Errorf("unknown callback command: %v", callbackUrl.String())
			return
		}
	}
}

//func LanguageOptions(whc botsfw.WebhookContext, mainMenu bool) tgbotapi.ReplyKeyboardMarkup {
//
//	buttons := [][]string{}
//	buttons = append(buttons, []string{whc.Locale().TitleWithIcon()})
//	row := []string{"", ""}
//	col := 0
//	whcLocalCode := whc.Locale().Code5
//	for _, locale := range trans.SupportedLocales {
//		logus.Debugf(c, "locale: %v, row: %v", locale, row)
//		if locale.Code5 == whcLocalCode {
//			logus.Debugf(c, "continue")
//			continue
//		}
//		row[col] = locale.TitleWithIcon()
//		logus.Debugf(c, "row: %v", row)
//		if col == 1 {
//			buttons = append(buttons, []string{row[0], row[1]})
//			logus.Debugf(c, "col: %v, keyboard: %v", col, buttons)
//			col = 0
//		} else {
//			col += 1
//		}
//	}
//	if mainMenu {
//		buttons = append(buttons, []string{MainMenuCommand.DefaultTitle(whc)})
//	}
//	logus.Debugf(c, "keyboard: %v", buttons)
//	return tgbotapi.NewReplyKeyboardUsingStrings(buttons)
//}
//
