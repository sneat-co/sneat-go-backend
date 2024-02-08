package dtb_fbm

//import (
//	"fmt"
//	"net/http"
//
//	"context"
//	"github.com/strongo/strongoapp"
//	"github.com/strongo/bots-api-fbm"
//	"github.com/bots-go-framework/bots-fw/botsfw"
//)
//
//func SetWhitelistedDomains(c context.Context, r *http.Request, bot botsfw.BotSettings, api fbmbotapi.GraphAPI) (err error) {
//	var whitelistedDomainsMessage fbmbotapi.WhitelistedDomainsMessage
//	switch bot.Env {
//	case strongoapp.EnvProduction:
//		whitelistedDomainsMessage = fbmbotapi.WhitelistedDomainsMessage{WhitelistedDomains: []string{
//			"https://debtstracker.io",
//			"https://splitbill.co",
//		}}
//	case strongoapp.EnvLocal:
//		domains := []string{
//			"https://debtstracker.local",
//		}
//		host := r.URL.Query().Get("host")
//		if host != "" {
//			domains = append(domains, fmt.Sprintf("https://%v", host))
//		}
//		whitelistedDomainsMessage = fbmbotapi.WhitelistedDomainsMessage{WhitelistedDomains: domains}
//	case strongoapp.EnvDevTest:
//		whitelistedDomainsMessage = fbmbotapi.WhitelistedDomainsMessage{WhitelistedDomains: []string{
//			"https://debtstracker-dev1.appspot.com",
//		}}
//	default:
//		err = fmt.Errorf("Unknown bot environment: %d=%v", bot.Env, strongoapp.EnvironmentNames[bot.Env])
//		return
//	}
//
//	host := fmt.Sprintf("https://%v", r.Host)
//
//	hasHost := false
//	for _, v := range whitelistedDomainsMessage.WhitelistedDomains {
//		if v == host {
//			hasHost = true
//			break
//		}
//	}
//	if !hasHost {
//		whitelistedDomainsMessage.WhitelistedDomains = append(whitelistedDomainsMessage.WhitelistedDomains, host)
//	}
//
//	if err = api.SetWhitelistedDomains(c, whitelistedDomainsMessage); err != nil {
//		return
//	}
//	return
//}
