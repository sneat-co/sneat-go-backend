package debtustgbots

import (
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/debtusbots/debtusbotconst"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/debtusbots/profiles/debtusbot"
	"github.com/strongo/i18n"
)

func GetDebtusBotProfile(errFooterText func() string) botsfw.BotProfile {
	_ = errFooterText
	return botsfw.NewBotProfile(debtusbotconst.DebtusBotProfileID, &debtusbot.Router, newBotChatData, newBotUserData, newAppUserData, getAppUserByID, i18n.LocaleEnUS, nil)
	//commandsByType := map[botsfw.WebhookInputType][]botsfw.Command{
	//	//botsfw.WebhookInputText: {startCommand},
	//}
	//shared.AddSharedCommands(commandsByType)
	//router := botsfw.NewWebhookRouter(commandsByType, errFooterText)
	//return shared.NewProfile(debtusbotconst.DebtusBotProfileID, &router)
}
