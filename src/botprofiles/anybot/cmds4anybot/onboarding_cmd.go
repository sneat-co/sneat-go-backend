package cmds4anybot

import (
	"errors"
	"github.com/bots-go-framework/bots-fw/botinput"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/dbo4userus"
)

var onboardingCommand = botsfw.Command{
	Code:       "onboarding",
	Commands:   []string{"/onboarding"},
	InputTypes: []botinput.WebhookInputType{botinput.WebhookInputText, botinput.WebhookInputCallbackQuery},
	Action:     onboardingAction,
	//CallbackAction: onboardingCallbackAction,
}

var ErrOnboardingCompleted = errors.New("onboarding completed")

func onboardingAction(whc botsfw.WebhookContext) (m botsfw.MessageFromBot, err error) {
	var user dbo4userus.UserEntry
	if user, err = GetUser(whc); err != nil {
		return
	}
	if user.Data.PreferredLocale == "" {
		m, err = onboardingAskLocaleAction(whc, "", nil)
		return
	}
	if user.Data.PrimaryCurrency == "" {
		m.Text = "Please select your primary currency"
		return
	}
	return m, ErrOnboardingCompleted
}
