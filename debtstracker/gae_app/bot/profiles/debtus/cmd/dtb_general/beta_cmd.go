package dtb_general

import (
	"fmt"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/auth"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/common"
)

const BetaCommandCode = "beta"

var BetaCommand = botsfw.Command{
	Code:     BetaCommandCode,
	Commands: []string{"/beta"},
	Action: func(whc botsfw.WebhookContext) (botsfw.MessageFromBot, error) {
		bot := whc.GetBotSettings()
		token := auth.IssueToken(whc.AppUserID(), whc.BotPlatform().ID()+":"+bot.Code, false)
		host := common.GetWebsiteHost(bot.Code)
		betaUrl := fmt.Sprintf(
			"https://%v/app/#lang=%v&secret=%v",
			host, whc.Locale().SiteCode(), token,
		)
		return whc.NewMessage(betaUrl), nil
	},
}
