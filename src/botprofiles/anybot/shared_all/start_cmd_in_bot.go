package shared_all

import (
	"errors"
)

var ErrUnknownStartParam = errors.New("unknown start parameter")

//func startInBotAction(whc botsfw.WebhookContext, startParams []string, startInBot StartInBotActionFunc, setMainMenu SetMainMenuFunc) (m botsfw.MessageFromBot, err error) {
//	logus.Debugf(whc.Context(), "startInBotAction() => startParams: %v", startParams)
//	if m, err = startInBot(whc, startParams); err != nil {
//		if errors.Is(err, ErrUnknownStartParam) {
//			if whc.ChatData().GetPreferredLanguage() == "" {
//				if m, err = OnboardingAskLocaleAction(whc, whc.Translate(trans.MESSAGE_TEXT_HI)+"\n\n", setMainMenu); err != nil {
//					err = fmt.Errorf("failed in OnboardingAskLocaleAction(): %w", err)
//					return
//				}
//			}
//		}
//		return
//	}
//	if len(startParams) > 0 {
//		switch {
//		case strings.HasPrefix(startParams[0], "how-to"):
//			return howToCommand.Action(whc)
//		}
//	}
//	if m, err = startInBotWelcomeAction(whc); err != nil {
//		err = errors.New("failed in startInBotWelcomeAction(): " + err.Error())
//	}
//	return
//}

//func startInBotWelcomeAction(whc botsfw.WebhookContext) (m botsfw.MessageFromBot, err error) {
//
//	var userName string
//
//	var user dbo4userus.UserEntry
//	if user, err = GetUser(whc); err != nil {
//		if !dal.IsNotFound(err) {
//			return
//		}
//		var botUser record.DataWithID[string, botsfwmodels.PlatformUserData]
//		if botUser, err = whc.BotUser(); err != nil && !dal.IsNotFound(err) {
//			return m, fmt.Errorf("failed to get bot user data: %w", err)
//		}
//		if dal.IsNotFound(err) {
//			userName = "stranger"
//			err = nil
//		} else {
//			botUserBaseData := botUser.Data.BaseData()
//			userName = botUserBaseData.FirstName
//			if userName == "" {
//				userName = botUserBaseData.LastName
//				if userName == "" {
//					userName = botUserBaseData.UserName
//				}
//				if userName == "" {
//					userName = "stranger"
//				}
//			}
//		}
//	} else {
//		userName = user.Data.Names.FirstName
//	}
//
//	buf := new(bytes.Buffer)
//
//	buf.WriteString(whc.Translate(trans.MESSAGE_TEXT_HI_USERNAME, userName))
//	buf.WriteString(" ")
//
//	buf.WriteString(whc.Translate(trans.SPLITUS_TEXT_HI))
//	buf.WriteString("\n\n")
//	buf.WriteString(whc.Translate(trans.SPLITUS_TEXT_ABOUT_ME_AND_CO))
//
//	buf.WriteString("\n\n")
//	buf.WriteString(whc.Translate(trans.MESSAGE_TEXT_ASK_LANG))
//	m.Text = buf.String()
//
//	m.Format = botsfw.MessageFormatHTML
//	m.Keyboard = LangKeyboard
//	return
//}

//func onStartCallbackInBot(whc botsfw.WebhookContext, params BotParams) (m botsfw.MessageFromBot, err error) {
//	ctx := whc.Context()
//	logus.Debugf(ctx, "onStartCallbackInBot()")
//
//	if m, err = params.InBotWelcomeMessage(whc); err != nil {
//		return
//	}
//
//	return
//}