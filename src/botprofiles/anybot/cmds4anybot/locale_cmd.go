package cmds4anybot

import (
	"bytes"
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

func createOnboardingAskLocaleCommand(setMainMenu SetMainMenuFunc) botsfw.Command {
	return botsfw.Command{
		Code:       onboardingAskLocaleCommandCode,
		InputTypes: []botinput.WebhookInputType{botinput.WebhookInputText, botinput.WebhookInputCallbackQuery},
		ExactMatch: trans.ChooseLocaleIcon,
		Action: func(whc botsfw.WebhookContext) (m botsfw.MessageFromBot, err error) {
			return OnboardingAskLocaleAction(whc, "", setMainMenu)
		},
	}
}

func OnboardingAskLocaleAction(whc botsfw.WebhookContext, messagePrefix string, setMainMenu SetMainMenuFunc) (m botsfw.MessageFromBot, err error) {
	chatEntity := whc.ChatData()

	if chatEntity.IsAwaitingReplyTo(onboardingAskLocaleCommandCode) {
		messageText := whc.Input().(botinput.WebhookTextMessage).Text()
		for _, locale := range trans.SupportedLocales {
			if locale.TitleWithIcon() == messageText {
				return setPreferredLanguageAction(whc, locale.Code5, "onboarding", setMainMenu)
			}
		}
		m = whc.NewMessageByCode(trans.MESSAGE_TEXT_UNKNOWN_LANGUAGE)
		//localesReplyKeyboard.OneTimeKeyboard = true
		m.Keyboard = localesReplyKeyboard
	} else {
		m.Text = messagePrefix + m.Text
		chatEntity.SetAwaitingReplyTo(onboardingAskLocaleCommandCode)
		m = whc.NewMessageByCode(trans.MESSAGE_TEXT_ONBOARDING_ASK_TO_CHOOSE_LANGUAGE, whc.Input().GetSender().GetFirstName())
		m.Format = botsfw.MessageFormatHTML
		//localesReplyKeyboard.OneTimeKeyboard = true
		m.Keyboard = localesReplyKeyboard
	}
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

func newSetLocaleCallbackCommand(setMainMenu SetMainMenuFunc) botsfw.Command {
	return botsfw.Command{
		Code: SettingsLocaleSetCallbackPath,
		CallbackAction: func(whc botsfw.WebhookContext, callbackUrl *url.URL) (m botsfw.MessageFromBot, err error) {
			return setPreferredLanguageAction(whc, callbackUrl.Query().Get("code5"), callbackUrl.Query().Get("mode"), setMainMenu)
		},
	}
}

func setPreferredLanguageAction(whc botsfw.WebhookContext, code5, mode string, setMainMenu SetMainMenuFunc) (m botsfw.MessageFromBot, err error) {
	ctx := whc.Context()
	logus.Debugf(ctx, "setPreferredLanguageAction(code5=%v, mode=%v)", code5, mode)

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
		setMainMenu(whc, &m)
		if _, err = whc.Responder().SendMessage(ctx, m, botsfw.BotAPISendMessageOverHTTPS); err != nil {
			logus.Errorf(ctx, "Failed to notify userEntity about selected language: %v", err)
			// Not critical, lets continue
		}
		return aboutDrawAction(whc, nil)
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

const (
	aboutDrawCommandCode = "about-draw"
	joinDrawCommandCode  = "join-draw"
)

var aboutDrawCommand = botsfw.Command{
	Commands:   []string{"/draw"},
	Code:       aboutDrawCommandCode,
	InputTypes: []botinput.WebhookInputType{botinput.WebhookInputText, botinput.WebhookInputCallbackQuery},
	Action: func(whc botsfw.WebhookContext) (m botsfw.MessageFromBot, err error) {
		return aboutDrawAction(whc, nil)
	},
	CallbackAction: aboutDrawAction,
}

var joinDrawCommand = botsfw.Command{
	Code:           joinDrawCommandCode,
	InputTypes:     []botinput.WebhookInputType{botinput.WebhookInputCallbackQuery},
	CallbackAction: aboutDrawAction,
}

func aboutDrawAction(whc botsfw.WebhookContext, callbackUrl *url.URL) (m botsfw.MessageFromBot, err error) {
	ctx := whc.Context()
	buf := new(bytes.Buffer)
	sender := whc.Input().GetSender()
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
					CallbackData: aboutDrawCommandCode,
				},
			},
		)
		return
	} else {
		m.IsEdit = true
		buf.WriteString(whc.Translate(trans.MESSAGE_TEXT_ABOUT_DRAW_MORE))
		m.Text = buf.String()
		switch callbackUrl.Path {
		case aboutDrawCommandCode:
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
			if _, err = whc.Responder().SendMessage(ctx, m, botsfw.BotAPISendMessageOverHTTPS); err != nil {
				logus.Warningf(ctx, "Failed to edit message: %v", err)
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
