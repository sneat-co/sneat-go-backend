package dtb_general

import (
	"fmt"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/sneat-co/sneat-go-backend/src/auth/token4auth"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/common4debtus"
)

const BetaCommandCode = "beta"

var BetaCommand = botsfw.Command{
	Code:     BetaCommandCode,
	Commands: []string{"/beta"},
	Action: func(whc botsfw.WebhookContext) (botsfw.MessageFromBot, error) {
		bot := whc.GetBotSettings()
		userID := whc.AppUserID()
		botPlatformID := whc.BotPlatform().ID()
		token := token4auth.IssueBotToken(userID, botPlatformID, bot.Code)
		host := common4debtus.GetWebsiteHost(bot.Code)
		betaUrl := fmt.Sprintf(
			"https://%v/app/#lang=%v&secret=%v",
			host, whc.Locale().SiteCode(), token,
		)
		return whc.NewMessage(betaUrl), nil
	},
}
