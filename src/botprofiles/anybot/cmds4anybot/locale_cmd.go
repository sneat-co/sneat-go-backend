package cmds4anybot

import (
	"context"
	"fmt"
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw/botinput"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/debtstracker-translations/trans"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/debtusbots/profiles/debtusbot/cmd/dtb_general"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/dal4userus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/i18n"
	"github.com/strongo/logus"
	"net/url"
	"strings"
)

const (
	SettingsLocaleListCallbackPath = "settings/locale/list"
	SettingsLocaleSetCallbackPath  = "settings/locale/set"
)

//var localesReplyKeyboard = tgbotapi.NewReplyKeyboard(
//	[]tgbotapi.KeyboardButton{
//		{Text: i18n.LocaleEnUS.TitleWithIcon()},
//		{Text: i18n.LocaleRuRu.TitleWithIcon()},
//	},
//	[]tgbotapi.KeyboardButton{
//		{Text: i18n.LocaleEsEs.TitleWithIcon()},
//		{Text: i18n.LocaleItIt.TitleWithIcon()},
//	},
//	[]tgbotapi.KeyboardButton{
//		{Text: i18n.LocaleDeDe.TitleWithIcon()},
//		{Text: i18n.LocaleFaIr.TitleWithIcon()},
//	},
//)

func getOnboardingLocalesKeyboard(callbackPath string) *tgbotapi.InlineKeyboardMarkup {
	localeInlineKeyboardButton := func(locale i18n.Locale) tgbotapi.InlineKeyboardButton {
		return tgbotapi.InlineKeyboardButton{
			Text:         locale.TitleWithIcon(),
			CallbackData: callbackPath + "?locale=" + locale.Code5,
		}
	}
	return tgbotapi.NewInlineKeyboardMarkup(
		[]tgbotapi.InlineKeyboardButton{
			localeInlineKeyboardButton(i18n.LocaleEnUK),
		},
		[]tgbotapi.InlineKeyboardButton{
			localeInlineKeyboardButton(i18n.LocaleRuRu),
			localeInlineKeyboardButton(i18n.LocaleUaUa),
		},
		[]tgbotapi.InlineKeyboardButton{
			localeInlineKeyboardButton(i18n.LocaleEsEs),
			localeInlineKeyboardButton(i18n.LocalePtPt),
		},
		[]tgbotapi.InlineKeyboardButton{
			localeInlineKeyboardButton(i18n.LocaleFrFr),
			localeInlineKeyboardButton(i18n.LocaleItIt),
		},
		[]tgbotapi.InlineKeyboardButton{
			localeInlineKeyboardButton(i18n.LocaleDeDe),
			localeInlineKeyboardButton(i18n.LocaleFaIr),
		},
		[]tgbotapi.InlineKeyboardButton{
			localeInlineKeyboardButton(i18n.LocaleZhCn),
			localeInlineKeyboardButton(i18n.LocaleJaJp),
		},
		//[]tgbotapi.InlineKeyboardButton{
		//	{Text: "Autodetect", WebApp: &tgbotapi.WebappInfo{Url: "https://sneat.app/telegram-webapp/detect-locale"}},
		//},
	)
}

func onStartAskLocaleCallbackAction(whc botsfw.WebhookContext, localeCode string, setMainMenu SetMainMenuFunc, getWelcomeMessageText WelcomeMessageProvider) (m botsfw.MessageFromBot, err error) {
	if m, err = setPreferredLocaleAction(whc, localeCode, "onboarding", setMainMenu, getWelcomeMessageText); err != nil {
		return m, fmt.Errorf("failed to setPreferredLocaleAction(): %w", err)
	}
	return
}

func onStartAskLocaleAction(whc botsfw.WebhookContext, setMainMenu SetMainMenuFunc, getWelcomeText WelcomeMessageProvider) (m botsfw.MessageFromBot, err error) {
	chatEntity := whc.ChatData()

	if chatEntity.IsAwaitingReplyTo(StartCommandCode) {
		messageText := whc.Input().(botinput.WebhookTextMessage).Text()
		for _, locale := range trans.SupportedLocales {
			if locale.TitleWithIcon() == messageText {
				return setPreferredLocaleAction(whc, locale.Code5, "onboarding", setMainMenu, getWelcomeText)
			}
		}
		m = whc.NewMessageByCode(trans.MESSAGE_TEXT_UNKNOWN_LANGUAGE)
	} else {
		m = whc.NewMessageByCode(trans.MESSAGE_TEXT_ONBOARDING_ASK_TO_CHOOSE_LANGUAGE)
	}
	m.Keyboard = getOnboardingLocalesKeyboard("start")
	m.Format = botsfw.MessageFormatHTML
	return
}

var AskPreferredLocaleFromSettingsCallback = botsfw.Command{
	Code:       SettingsLocaleListCallbackPath,
	InputTypes: []botinput.WebhookInputType{botinput.WebhookInputCallbackQuery},
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

func newSetLocaleCallbackCommand(setMainMenu SetMainMenuFunc, getWelcomeMessageText WelcomeMessageProvider) botsfw.Command {
	return botsfw.Command{
		Code: SettingsLocaleSetCallbackPath,
		CallbackAction: func(whc botsfw.WebhookContext, callbackUrl *url.URL) (m botsfw.MessageFromBot, err error) {
			return setPreferredLocaleAction(whc, callbackUrl.Query().Get("code5"), callbackUrl.Query().Get("mode"), setMainMenu, getWelcomeMessageText)
		},
	}
}

func setPreferredLocaleAction(whc botsfw.WebhookContext, code5, mode string, setMainMenu SetMainMenuFunc, getWelcomeMessageText WelcomeMessageProvider) (m botsfw.MessageFromBot, err error) {
	ctx := whc.Context()
	logus.Debugf(ctx, "setPreferredLocaleAction(code5=%v, mode=%v)", code5, mode)

	var (
		localeChanged  bool
		selectedLocale i18n.Locale
	)

	chatData := whc.ChatData()
	preferredLocale := chatData.GetPreferredLanguage()
	logus.Debugf(ctx, "userEntity.PreferredLanguage: %v, chatData.GetPreferredLanguage(): %v, code5: %v", preferredLocale, chatData.GetPreferredLanguage(), code5)
	if preferredLocale != code5 || chatData.GetPreferredLanguage() != code5 {
		logus.Debugf(ctx, "PreferredLanguage will be updated for userEntity & chat entities.")
		for _, locale := range trans.SupportedLocalesByCode5 {
			if locale.Code5 == code5 {
				_ = whc.SetLocale(locale.Code5)

				if appUserID := whc.AppUserID(); appUserID != "" {
					userCtx := facade.NewUserContext(appUserID)
					if err = dal4userus.RunUserWorker(ctx, userCtx, true, func(ctx context.Context, tx dal.ReadwriteTransaction, params *dal4userus.UserWorkerParams) (err error) {
						if params.UserUpdates, err = params.User.Data.SetPreferredLocale(locale.Code5); err != nil {
							return fmt.Errorf("%w: failed to set preferred locale for user", err)
						}
						params.User.Record.MarkAsChanged()
						chatData.SetPreferredLanguage(locale.Code5)
						chatData.SetAwaitingReplyTo("")
						//chatKey := botsfwmodels.NewChatKey(whc.GetBotCode(), whc.MustBotChatID())
						if err = whc.SaveBotChat(ctx); err != nil {
							return
						}
						return
					}); err != nil {
						return
					}
					localeChanged = true
					selectedLocale = locale
					if whc.GetBotSettings().Env == "prod" {
						ga := whc.GA()
						gaEvent := ga.GaEventWithLabel("settings", "locale-changed", strings.ToLower(locale.Code5))
						if gaErr := ga.Queue(gaEvent); gaErr != nil {
							logus.Warningf(ctx, "Failed to log event: %v", gaErr)
						} else {
							logus.Infof(ctx, "GA event queued: %v", gaEvent)
						}
					}
				}

				break
			}
		}
		if !localeChanged {
			logus.Errorf(ctx, "Unknown locale: %v", code5)
		}
	} else {
		selectedLocale = i18n.GetLocaleByCode5(chatData.GetPreferredLanguage())
	}
	//if localeChanged {

	switch mode {
	case "onboarding":
		logus.Debugf(ctx, "whc.Locale().Code5: %v", whc.Locale().Code5)
		m = whc.NewMessageByCode(trans.MESSAGE_TEXT_YOUR_SELECTED_PREFERRED_LANGUAGE, selectedLocale.NativeTitle)
		if getWelcomeMessageText != nil {
			var welcomeText string
			if welcomeText, err = getWelcomeMessageText(whc); err != nil {
				return
			}
			m.Text = welcomeText + "\n\n" + m.Text
			m.Format = botsfw.MessageFormatHTML
		}
		m.IsEdit = true
		if setMainMenu != nil {
			setMainMenu(whc, &m)
		}
		//if _, err = whc.Responder().SendMessage(ctx, m, botsfw.BotAPISendMessageOverHTTPS); err != nil {
		//	logus.Errorf(ctx, "Failed to notify userEntity about selected language: %v", err)
		//	// Not critical, lets continue
		//}
		return
	case "settings":
		if localeChanged {
			if m, err = dtb_general.MainMenuAction(whc, whc.Translate(trans.MESSAGE_TEXT_LOCALE_CHANGED, selectedLocale.TitleWithIcon()), false); err != nil {
				return
			}
			if _, err = whc.Responder().SendMessage(ctx, m, botsfw.BotAPISendMessageOverHTTPS); err != nil {
				return m, err
			}
			return SettingsMainAction(whc)
			//if _, err = whc.Responder().SendMessage(ctx, m, botsfw.BotAPISendMessageOverHTTPS); err != nil {
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
