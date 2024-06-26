package website

import (
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/strongo/i18n"
	"github.com/strongo/logus"
	"github.com/strongo/strongoapp"
	"google.golang.org/appengine/v2"
	"net/http"
	"strconv"

	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/dtdal"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/website/admin"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/website/pages"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/website/redirects"
	//"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/common"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/auth"
	//"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/api"
	"github.com/julienschmidt/httprouter"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/website/pages/inspector"
)

type router interface {
	GET(path string, handle httprouter.Handle)
}

func InitWebsite(router router) {
	router.GET("/", pages.IndexRootPage)

	redirects.InitRedirects(router)

	for _, locale := range i18n.LocalesByCode5 {
		localeSiteCode := locale.SiteCode()
		router.GET(fmt.Sprintf("/%v/ads", localeSiteCode), pages.AdsPage)
		router.GET(fmt.Sprintf("/%v/help-us", localeSiteCode), pages.HelpUsPage)
		router.GET(fmt.Sprintf("/%v/login", localeSiteCode), LoginHandler)
		router.GET(fmt.Sprintf("/%v/counterparty", localeSiteCode), pages.CounterpartyPage)
		router.GET(fmt.Sprintf("/%v/", localeSiteCode), pages.IndexPage)
		//strongoapp.AddHTTPHandler(fmt.Sprintf("/%v/create-mass-invite", localeSiteCode), api.AuthOnly(CreateInvitePage))

	}
	router.GET("/en/songs/annie-iou-a-dance", pages.AnnieIOUaDancePage)
	router.GET("/en/songs/iou-by-dappy", pages.IOWDappyPage)

	admin.InitAdmin(router)
	inspector.InitInspector(router)
}

func CreateInvitePage(w http.ResponseWriter, r *http.Request, authInfo auth.AuthInfo) {
	c := appengine.NewContext(r)
	logus.Infof(c, "CreateInvitePage()")
	//panic("Not implemented")
	userID := authInfo.UserID
	//session, _ := common.GetSession(r)
	//userID := session.GetUserID()
	//if userID == 0 {
	//	w.WriteHeader(http.StatusUnauthorized)
	//	return
	//}
	switch r.Method {
	case http.MethodGet:
		http.ServeFile(w, r, "templates/create-mass-invite.html")
		return
	case http.MethodPost:
		if err := r.ParseForm(); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(err.Error()))
		}
		inviteCode := r.Form.Get("Code")
		if !dtdal.InviteCodeRegex.Match([]byte(inviteCode)) {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(fmt.Sprintf("Invate code [%v] does not match pattern: %v", inviteCode, dtdal.InviteCodeRegex.String())))
			return
		}
		maxClaimsCount, err := strconv.ParseInt(r.Form.Get("MaxClaimsCount"), 10, 32)
		if err != nil || inviteCode == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if _, err = dtdal.Invite.GetInvite(c, nil, inviteCode); err != nil {
			if dal.IsNotFound(err) {
				w.WriteHeader(http.StatusBadRequest)
				_, _ = w.Write([]byte(fmt.Sprintf("Invate code [%v] already exists", inviteCode)))
			} else {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(err.Error()))
			}
			return
		}
		//translator := i18n.NewSingleMapTranslator(i18n.GetLocaleByCode5(i18n.LocaleCodeEnUS), i18n.NewMapTranslator(c, trans.TRANS))
		ec := strongoapp.NewExecutionContext(c)
		if _, err = dtdal.Invite.CreateMassInvite(ec, userID, inviteCode, int32(maxClaimsCount), "web"); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}
		_, _ = w.Write([]byte(fmt.Sprintf("Invite created, code: %v, MaxClaimsCount: %v", inviteCode, maxClaimsCount)))
	default:
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}
