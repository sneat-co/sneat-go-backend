package dtb_fbm

//import (
//	"bytes"
//	"fmt"
//	"net/http"
//	"time"
//
//	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/bot/platforms/fbmbots"
//	"context"
//	"github.com/julienschmidt/httprouter"
//	"github.com/strongo/bots-api-fbm"
//	"github.com/strongo/logus"
//	"google.golang.org/appengine/v2"
//)
//
//func SetupFbm(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
//	query := r.URL.Query()
//
//	botID := query.Get("bot")
//	c := appengine.NewContext(r)
//	bot, ok := fbmbots.Bots(c).ByCode[botID]
//	if !ok {
//		w.Write([]byte("Unknown bot: " + botID))
//		return
//	}
//
//	c, _ = context.WithDeadline(c, time.Now().Add(20*time.Second))
//	api := fbmbotapi.NewGraphAPI(urlfetch.Client(c), bot.Token)
//
//	var err error
//
//	var buffer bytes.Buffer
//
//	reportError := func(err error) {
//		logus.Errorf(c, err.Error())
//		w.WriteHeader(http.StatusInternalServerError)
//		w.Write(buffer.Bytes())
//		w.Write([]byte(err.Error()))
//	}
//
//	if query.Get("whitelist-domains") == "1" {
//		if err = SetWhitelistedDomains(c, r, bot, api); err != nil {
//			reportError(err)
//			return
//		}
//		buffer.WriteString("Whitelisted domains\n")
//	}
//
//	if query.Get("enable-get-started") == "1" {
//		getStartedMessage := fbmbotapi.GetStartedMessage{}
//
//		getStartedMessage.GetStarted.Payload = "fbm-get-started"
//
//		if err = api.SetGetStarted(c, getStartedMessage); err != nil {
//			reportError(err)
//			return
//		}
//		buffer.WriteString("Enabled 'Get started'\n")
//	}
//
//	if query.Get("set-persistent-menu") == "1" {
//		if err = SetPersistentMenu(c, r, bot, api); err != nil {
//			reportError(err)
//			return
//		}
//		buffer.WriteString("Enabled 'Persistent menu'\n")
//	}
//	logus.Debugf(c, buffer.String())
//	w.Header().Set("Content-ExtraType", "text/plain")
//	w.Write(buffer.Bytes())
//	w.Write([]byte(fmt.Sprintf("OK! %v", time.Now())))
//}
