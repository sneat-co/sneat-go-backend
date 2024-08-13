package shared_all

import (
	"bytes"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/sneat-co/debtstracker-translations/trans"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/dbo4userus"
	"github.com/strongo/logus"
	"strings"

	"errors"
)

var ErrUnknownStartParam = errors.New("unknown start parameter")

func startInBotAction(whc botsfw.WebhookContext, startParams []string, botParams BotParams) (m botsfw.MessageFromBot, err error) {
	logus.Debugf(whc.Context(), "startInBotAction() => startParams: %v", startParams)
	if m, err = botParams.StartInBotAction(whc, startParams); err != nil {
		if err == ErrUnknownStartParam {
			if whc.ChatData().GetPreferredLanguage() == "" {
				return onboardingAskLocaleAction(whc, whc.Translate(trans.MESSAGE_TEXT_HI)+"\n\n", botParams)
			}
		}
		return
	}
	if len(startParams) > 0 {
		switch {
		case strings.HasPrefix(startParams[0], "how-to"):
			return howToCommand.Action(whc)
		}
	}
	return startInBotWelcomeAction(whc, botParams)
}

func startInBotWelcomeAction(whc botsfw.WebhookContext, botParams BotParams) (m botsfw.MessageFromBot, err error) {
	var user dbo4userus.UserEntry
	if user, err = GetUser(whc); err != nil {
		return
	}

	buf := new(bytes.Buffer)

	buf.WriteString(whc.Translate(trans.MESSAGE_TEXT_HI_USERNAME, user.Data.Names.FirstName))
	buf.WriteString(" ")

	buf.WriteString(whc.Translate(trans.SPLITUS_TEXT_HI))
	buf.WriteString("\n\n")
	buf.WriteString(whc.Translate(trans.SPLITUS_TEXT_ABOUT_ME_AND_CO))

	buf.WriteString("\n\n")
	buf.WriteString(whc.Translate(trans.MESSAGE_TEXT_ASK_LANG))
	m.Text = buf.String()

	m.Format = botsfw.MessageFormatHTML
	m.Keyboard = LangKeyboard
	return
}

//func onStartCallbackInBot(whc botsfw.WebhookContext, params BotParams) (m botsfw.MessageFromBot, err error) {
//	c := whc.Context()
//	logus.Debugf(c, "onStartCallbackInBot()")
//
//	if m, err = params.InBotWelcomeMessage(whc); err != nil {
//		return
//	}
//
//	return
//}
