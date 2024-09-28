package dtb_general

import (
	"fmt"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/sneat-co/sneat-go-backend/src/modules/auth/token4auth"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/common4debtus"
)

const BetaCommandCode = "beta"

var BetaCommand = botsfw.Command{
	Code:     BetaCommandCode,
	Commands: []string{"/beta"},
	Action: func(whc botsfw.WebhookContext) (botsfw.MessageFromBot, error) {
		ctx := whc.Context()
		bot := whc.GetBotSettings()
		userID := whc.AppUserID()
		botPlatformID := whc.BotPlatform().ID()
		token, err := token4auth.IssueBotToken(ctx, userID, botPlatformID, bot.Code)
		if err != nil {
			return botsfw.MessageFromBot{}, err
		}
		host := common4debtus.GetWebsiteHost(bot.Code)
		betaUrl := fmt.Sprintf(
			"https://%v/app/#lang=%v&secret=%v",
			host, whc.Locale().SiteCode(), token,
		)
		return whc.NewMessage(betaUrl), nil
	},
}
