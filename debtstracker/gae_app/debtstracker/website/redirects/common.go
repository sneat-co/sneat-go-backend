package redirects

import (
	"bytes"
	"fmt"
	"github.com/strongo/logus"
	"google.golang.org/appengine/v2"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/auth"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/common"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/facade"
)

func redirectToWebApp(w http.ResponseWriter, r *http.Request, authRequired bool, path string, p2p map[string]string, optionalParams []string) {
	c := appengine.NewContext(r)
	query := r.URL.Query()

	authInfo, _, err := auth.Authenticate(w, r, authRequired)
	if authRequired && err != nil {
		return
	}

	var redirectTo bytes.Buffer
	redirectTo.WriteString("/app/")

	lang := query.Get("lang")
	if lang == "" {
		if authInfo.UserID != "" {
			user, err := facade.User.GetUserByID(c, nil, authInfo.UserID)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(err.Error()))
				return
			}
			lang = strings.ToLower(user.Data.GetPreferredLocale()[:2])
		} else {
			lang = "en" // TODO: Bad to hard-code. Try to get from receipt?
		}
	}

	redirectTo.WriteString("#" + path)

	if path != "" {
		redirectTo.WriteString("&")
	}
	redirectTo.WriteString("lang=" + lang)

	sep := ""

	for pn, pn2 := range p2p {
		if pv := query.Get(pn); pv == "" {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte("redirectToWebApp: missing required parameter: " + pn))
			return
		} else {
			pv = url.QueryEscape(pv)
			if pn == "id" && pn2 == "receipt" { // TODO: Dirty hack! Please fix!!!
				receiptID, err := common.DecodeID(pv)
				if err != nil {
					logus.Debugf(c, "Failed to decode receipt ID: %v", err)
					w.WriteHeader(http.StatusBadRequest)
					_, _ = w.Write([]byte(fmt.Sprintf("Failed to decod receipt ID: %v", err)))
					return
				}
				pv = strconv.FormatInt(receiptID, 10)
			}
			redirectTo.WriteString(sep + pn2 + "=" + pv)
		}
		sep = "&"
	}

	for _, p := range optionalParams {
		if v := query.Get(p); v != "" {
			redirectTo.WriteString(fmt.Sprintf("&%v=%v", p, url.QueryEscape(v)))
		}
	}

	if utm := query.Get("utm"); utm != "" {
		matches := reUtm.FindAllStringSubmatch(r.URL.RawQuery, -1) // TODO: Looks like a hack. Consider replacing ';' char with something else?
		if len(matches) == 1 {
			utm = matches[0][1]
			utmValues := strings.Split(utm, ";")
			if len(utmValues) == 3 {
				for i, p := range []string{"utm_source", "utm_medium", "utm_campaign"} {
					redirectTo.WriteString(fmt.Sprintf("&%v=%v", p, url.QueryEscape(utmValues[i])))
				}
			} else {
				logus.Warningf(c, "Parameter utm should consist of 3 values seprated by ';' character. Got: [%v]", utm)
			}
		} else {
			logus.Errorf(c, "reUtm: %v", matches)
		}
	} else {
		for _, p := range []string{"utm_source", "utm_medium", "utm_campaign"} {
			if v := query.Get(p); v != "" {
				redirectTo.WriteString(fmt.Sprintf("&%v=%v", p, url.QueryEscape(v)))
			}
		}
	}

	if authInfo.UserID > "" {
		redirectTo.WriteString("&secret=" + query.Get("secret"))
	}
	logus.Debugf(c, "Will redirect to: %v", redirectTo.String())
	http.Redirect(w, r, redirectTo.String(), http.StatusFound)
	//w.WriteHeader(http.StatusFound)
	//w.Header().Set("Location", redirectTo.String())
	//w.Write([]byte(fmt.Sprintf(`<html><head><meta http-equiv="refresh" content="0;URL='%v'" /></head></html>`, redirectTo.String())))
}

var reUtm = regexp.MustCompile(`[&#?]?utm=(.+?)(?:&|#|$)`)
