package fbmbots

import (
	"context"
	"errors"
	"github.com/bots-go-framework/bots-fw/botsfw"
)

const (
// fbmProdPageAccessToken  = "TODO_DO_NOT_STORE_IN_GIT"
// fbmTestPageAccessToken  = "TODO_DO_NOT_STORE_IN_GIT"
// fbmLocalPageAccessToken = "TODO_DO_NOT_STORE_IN_GIT"
//
// fbmSplitBillProdPageAccessToken = "TODO_DO_NOT_STORE_IN_GIT"
)

//type fbAppSecrets struct {
//	AppID     string
//	AppSecret string
//	//app       *fb.App
//}

//func (s *fbAppSecrets) App() *fb.App {
//	if s.app == nil {
//		s.app = fb.New(s.AppID, s.AppSecret)
//	}
//	return s.app
//}

var (
//	fbLocal = fbAppSecrets{
//		AppID:     "457648507752783",
//		AppSecret: "23ceb7a7f53516119fd60b19a309cb14",
//	}
//
//	fbDev = fbAppSecrets{
//		AppID:     "579129655604667",
//		AppSecret: "0e3ee2d65e8abae458f121e874950b73",
//	}
//
//	fbProd = fbAppSecrets{
//		AppID:     "454859831364984",
//		AppSecret: "72f6f7382dda3235d48e6a7d60bb4a6a",
//	}
)

var _bots botsfw.SettingsBy

func Bots(_ context.Context) botsfw.SettingsBy {
	//if len(_bots.ByCode) == 0 {
	//	_bots = botsfw.NewBotSettingsBy(nil,
	//		fbm.NewFbmBot(
	//			strongoapp.EnvProduction,
	//			bot.ProfileDebtus,
	//			"debtstracker",
	//			"1587055508253137",
	//			fbmProdPageAccessToken,
	//			"d6087a01-c728-4fdf-983c-1695d76236dc",
	//			common.GA_TRACKING_ID,
	//			trans.SupportedLocalesByCode5[i18n.LocaleCodeEnUS],
	//		),
	//		fbm.NewFbmBot(
	//			strongoapp.EnvProduction,
	//			bot.ProfileSplitus,
	//			"splitbill.co",
	//			"286238251784027",
	//			fbmSplitBillProdPageAccessToken,
	//			"e8535dd1-df3b-4c3f-bd2c-d4a822509bb3",
	//			common.GA_TRACKING_ID,
	//			trans.SupportedLocalesByCode5[i18n.LocaleCodeEnUS],
	//		),
	//		fbm.NewFbmBot(
	//			strongoapp.EnvDevTest,
	//			bot.ProfileDebtus,
	//			"debtstracker.dev",
	//			"942911595837341",
	//			fbmTestPageAccessToken,
	//			"4afb645e-b592-48e6-882c-89f0ec126fbb",
	//			"",
	//			trans.SupportedLocalesByCode5[i18n.LocaleCodeEnUS],
	//		),
	//		fbm.NewFbmBot(
	//			strongoapp.EnvLocal,
	//			bot.ProfileDebtus,
	//			"debtstracker.local",
	//			"300392587037950",
	//			fbmLocalPageAccessToken,
	//			"4afb645e-b592-48e6-882c-89f0ec126fbb",
	//			"",
	//			trans.SupportedLocalesByCode5[i18n.LocaleCodeEnUS],
	//		),
	//	)
	//}
	return _bots
}

var ErrUnknownHost = errors.New("Unknown host")

//func GetFbAppAndHost(r *http.Request) (fbApp *fb.App, host string, err error) {
//	switch r.Host {
//	case "debtstracker.io":
//		return fbProd.App(), r.Host, nil
//	case "debtstracker-io.appspot.com":
//		return fbProd.App(), "debtstracker.io", nil
//	case "debtstracker-dev1.appspot.com":
//		return fbDev.App(), r.Host, nil
//	case "debtstracker.local":
//		return fbLocal.App(), r.Host, nil
//	case "localhost":
//		return fbLocal.App(), "debtstracker.local", nil
//	default:
//		if strings.HasSuffix(r.Host, ".ngrok.io") {
//			return fbLocal.App(), "debtstracker.local", nil
//		}
//	}
//
//	return nil, "", fmt.Errorf(ErrUnknownHost, r.Host)
//}

//func getFbAppAndSession(c context.Context, r *http.Request, getSession func(fbApp *fb.App) (*fb.Session, error)) (
//	fbApp *fb.App, fbSession *fb.Session, err error,
//) {
//	log.Debugf(c, "getFbAppAndSession()")
//	if fbApp, _, err = GetFbAppAndHost(r); err != nil {
//		log.Errorf(c, "getFbAppAndSession() => Failed to get app")
//		return nil, nil, err
//	}
//	if fbSession, err = getSession(fbApp); err != nil {
//		log.Errorf(c, "getFbAppAndSession() => Failed to get session")
//		return nil, nil, err
//	}
//	log.Debugf(c, "getFbAppAndSession() => AppId: %v", fbApp.AppId)
//	return fbApp, fbSession, err
//}

//func FbAppAndSessionFromAccessToken(c context.Context, r *http.Request, accessToken string) (*fb.App, *fb.Session, error) {
//	return getFbAppAndSession(c, r, func(fbApp *fb.App) (fbSession *fb.Session, err error) {
//		fbSession = fbApp.Session(accessToken)
//		fbSession.HttpClient = dtdal.HttpClient(c)
//		return
//	})
//}
//
//func FbAppAndSessionFromSignedRequest(c context.Context, r *http.Request, signedRequest string) (*fb.App, *fb.Session, error) {
//	log.Debugf(c, "FbAppAndSessionFromSignedRequest()")
//	return getFbAppAndSession(c, r, func(fbApp *fb.App) (fbSession *fb.Session, err error) {
//		log.Debugf(c, "FbAppAndSessionFromSignedRequest() => getSession()")
//		//fbSession, err = fbApp.SessionFromSignedRequest(c, signedRequest, dtdal.HttpClient(c))
//		//if err != nil {
//		//	log.Debugf(c, "FbAppAndSessionFromSignedRequest() => getSession(): %v", err.Error())
//		//}
//		panic("not implemented")
//		return
//	})
//}
