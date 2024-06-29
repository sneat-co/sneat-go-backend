package admin

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/dtdal/gaedal"
	"google.golang.org/appengine/v2"
	"net/http"
)

func LatestPage(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	panic("disabled: switch to Firestore authentication")
	//c := appengine.NewContext(r)
	//
	//if !gaeUser.IsAdmin(c) {
	//	url, _ := gaeUser.LoginURL(c, r.RequestURI)
	//	w.WriteHeader(http.StatusTemporaryRedirect)
	//	w.Header().Add("location", url)
	//	return
	//}
	//
	//var users []models.DebutsAppUserDataOBSOLETE
	//userKeys, err := datastore.NewQuery(models.AppUserKind).Order("-DtCreated").Limit(50).GetAll(c, &users)
	//if err != nil {
	//	logus.Errorf(c, err.Error())
	//	w.WriteHeader(http.StatusInternalServerError)
	//	_, _ = w.Write([]byte(err.Error()))
	//}
	//
	//b := bufio.NewWriter(w)
	//_, _ = b.WriteString("<html><head><style>body{Font-Family:Verdana;font-size:x-small} td{padding: 2px 5px; background-color: #eee} th{padding: 2px 5px; text-align: left; background-color: #ddd} .num{text-align: right} div{float: left}</style></head>")
	//_, _ = b.WriteString("<body><h1>Latest</h1><hr>")
	//_, _ = b.WriteString("<div><h2>Users</h2><table cellspacing=1><thead><tr><th>#</th><th>Full ContactName</th><th>Contacts</th><th>Debts</th><th>Balance</th><th>Invited by</th></tr></thead><tbody>")
	//for i, user := range users {
	//	_, _ = b.WriteString("<tr>")
	//	_, _ = b.WriteString("<td class=num>")
	//	_, _ = b.WriteString(strconv.Itoa(i + 1))
	//	_, _ = b.WriteString("</td><td>")
	//	_, _ = b.WriteString(fmt.Sprintf("<a href='user?id=%v'>%v</a>", userKeys[i].IntID(), html.EscapeString(user.FullName())))
	//	_, _ = b.WriteString("</td><td class=num>")
	//	_, _ = b.WriteString(strconv.Itoa(user.TotalContactsCount()))
	//	_, _ = b.WriteString("</td><td class=num>")
	//	_, _ = b.WriteString(strconv.Itoa(user.CountOfTransfers))
	//	_, _ = b.WriteString("</td><td>")
	//	_, _ = b.WriteString(user.BalanceJson)
	//	_, _ = b.WriteString("</td><td>")
	//	if user.InvitedByUserID != "" {
	//		if invitedByUser, err := facade.User.GetUserByID(c, nil, user.InvitedByUserID); err != nil {
	//			_, _ = b.WriteString(err.Error())
	//		} else {
	//			_, _ = b.WriteString(fmt.Sprintf("<a href='user?id=%v>%v</a>')", user.InvitedByUserID, invitedByUser.Data.FullName()))
	//		}
	//	}
	//	_, _ = b.WriteString("</td></tr>")
	//}
	//_, _ = b.WriteString("</tbody></table></div>")
	//_, _ = b.WriteString("</body></html>")
	//_ = b.Flush()
}

func FixTransfersHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	c := appengine.NewContext(r)
	loadedCount, fixedCount, failedCount, err := gaedal.FixTransfers(c)
	stats := fmt.Sprintf("\nLoaded: %v, Fixed: %v, Failed: %v", loadedCount, fixedCount, failedCount)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
	}
	_, _ = w.Write([]byte(stats))
}
