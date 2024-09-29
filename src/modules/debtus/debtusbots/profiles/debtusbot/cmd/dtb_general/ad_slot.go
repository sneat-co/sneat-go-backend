package dtb_general

import (
	"fmt"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/sneat-co/debtstracker-translations/trans"
	"github.com/sneat-co/sneat-core-modules/common4all"
	"strings"
)

func AdSlot(whc botsfw.WebhookContext, place string) string {
	utmParams := common4all.FillUtmParams(whc, common4all.UtmParams{Campaign: place})
	link := fmt.Sprintf(`href="https://debtus.app/%v/ads#%v"`, whc.Locale().SiteCode(), utmParams)
	return strings.Replace(whc.Translate(trans.MESSAGE_TEXT_YOUR_AD_COULD_BE_HERE), "href", link, 1)
}
