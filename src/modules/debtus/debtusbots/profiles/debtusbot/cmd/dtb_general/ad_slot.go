package dtb_general

import (
	"fmt"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/sneat-co/debtstracker-translations/trans"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/common4debtus"
	"strings"
)

func AdSlot(whc botsfw.WebhookContext, place string) string {
	utmParams := common4debtus.FillUtmParams(whc, common4debtus.UtmParams{Campaign: place})
	link := fmt.Sprintf(`href="https://debtusbot.io/%v/ads#%v"`, whc.Locale().SiteCode(), utmParams)
	return strings.Replace(whc.Translate(trans.MESSAGE_TEXT_YOUR_AD_COULD_BE_HERE), "href", link, 1)
}