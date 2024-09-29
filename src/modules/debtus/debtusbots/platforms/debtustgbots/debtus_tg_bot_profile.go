package debtustgbots

import (
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/sneat-co/sneat-go-backend/src/coremodules/anybot"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/debtusbots/debtusbotconst"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/debtusbots/profiles/debtusbot"
)

func GetDebtusBotProfile(errFooterText func() string) botsfw.BotProfile {
	_ = errFooterText
	//return botsfw.NewBotProfile(debtusbotconst.DebtusBotProfileID, &debtusbot.Router, newBotChatData, newBotUserData, newAppUserData, getAppUserByID, i18n.LocaleEnUS, nil)
	return anybot.NewProfile(debtusbotconst.DebtusBotProfileID, &debtusbot.Router)
}
