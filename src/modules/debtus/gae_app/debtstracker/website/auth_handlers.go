package website

import (
	"github.com/dal-go/dalgo/dal"
	"github.com/julienschmidt/httprouter"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/common4debtus"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/dal4userus"
	"github.com/strongo/logus"
	"net/http"
	"strings"
	"time"
)

func LoginHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	c := r.Context()

	q := r.URL.Query()
	userID := q.Get("user")
	if userID == "" {
		w.WriteHeader(http.StatusBadRequest)
		logus.Infof(c, "Invalid user parameter")
		return
	}
	secret := q.Get("secret")
	secretItems := strings.Split(secret, ":")
	expirySecStr := secretItems[0]
	logus.Infof(c, "expirySeconds: %v; secret: %v", expirySecStr, secret)
	expirySeconds, err := common4debtus.DecodeID(expirySecStr)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logus.Infof(c, "Failed to decode expiry bytes to seconds")
		return
	}

	expiresAt := time.Unix(expirySeconds, 0)

	expectedSecret := common4debtus.SignStrWithExpiry(c, userID, expiresAt)
	if secret != expectedSecret {
		w.WriteHeader(http.StatusUnauthorized)
		logus.Infof(c, "Invalid secret")
		return
	}

	if expiresAt.Before(time.Now()) {
		w.WriteHeader(http.StatusUnauthorized)
		logus.Infof(c, "expiresAt.Before(time.Now())")
		_, _ = w.Write([]byte("<html><body style=font-size:xx-large>Your secret has expired. Please generate a new link</body></html>"))
		return
	}

	if _user, err := dal4userus.GetUserByID(c, nil, userID); err != nil {
		if dal.IsNotFound(err) {
			w.WriteHeader(http.StatusNotFound)
			logus.Infof(c, err.Error())
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			logus.Errorf(c, err.Error())
		}
		return
	} else {
		if _user.Data.Email != "" {
			logus.Infof(c, "_user.Email: %v", _user.Data.Email)
		} else {
			panic("disabled: switch to Firestore authentication") // TODO: switch to Firestore authentication
			//gaeUser := gaeuser.Current(c)
			//if gaeUser == nil {
			//	logus.Infof(c, "appengine.user.Current(): nil")
			//} else {
			//	if gaeUser.Email == "" {
			//		logus.Infof(c, "gaeUser.Email is empty")
			//	} else {
			//		logus.Infof(c, "gaeUser.Email: %v", gaeUser.Email)
			//		var db dal.DB
			//		if db, err = facade4debtus.GetDatabase(c); err != nil {
			//			w.WriteHeader(http.StatusInternalServerError)
			//			logus.Errorf(c, err.Error())
			//			return
			//		}
			//		err = db.RunReadwriteTransaction(c, func(tctx context.Context, tx dal.ReadwriteTransaction) error {
			//			u, err := dal4userus.GetUserByID(tctx, tx, userID)
			//			if err != nil {
			//				return err
			//			}
			//			if u.Data.EmailAddress == "" {
			//				u.Data.SetEmail(gaeUser.Email, true)
			//				if err = facade4debtus.User.SaveUserOBSOLETE(tctx, tx, u); err != nil {
			//					return fmt.Errorf("failed to save user: %w", err)
			//				}
			//			}
			//			return err
			//		}, nil)
			//		if err != nil {
			//			w.WriteHeader(http.StatusInternalServerError)
			//			logus.Errorf(ctx, err.Error())
			//		}
			//	}
			//}
		}
	}

	panic("Not implemented")
	//session, _ := anybot.GetSession(r)
	//session.SetUserID(userID, w)
	//if err = session.Save(r, w); err != nil {
	//	w.WriteHeader(http.StatusInternalServerError)
	//	logus.Errorf(c, err.Error())
	//	return
	//}

	//w.Write([]byte("<html><body style=font-size:xx-large>User signed</body></html>"))
}
