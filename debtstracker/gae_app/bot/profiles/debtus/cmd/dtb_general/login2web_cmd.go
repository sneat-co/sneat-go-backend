package dtb_general

import (
	"fmt"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/sneat-co/debtstracker-translations/trans"
	"strings"

	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/common"
)

const LOGIN2WEB_COMMAND = "login2web"

var Login2WebCommand = botsfw.Command{
	Code:     LOGIN2WEB_COMMAND,
	Commands: []string{"/login"},
	Action: func(whc botsfw.WebhookContext) (m botsfw.MessageFromBot, err error) {
		mt := whc.Translate(trans.MESSAGE_TEXT_LOGIN_TO_WEB_APP)
		linker := common.NewLinkerFromWhc(whc)
		mt = strings.Replace(mt, "<a>", fmt.Sprintf(`<a href="%v">`, linker.ToMainScreen(whc)), 1)
		m = whc.NewMessage(mt)
		m.Format = botsfw.MessageFormatHTML
		m.DisableWebPagePreview = true
		return
	},
}
