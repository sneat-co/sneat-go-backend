package unsorted

import (
	"errors"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/auth/api4auth"
	"github.com/sneat-co/sneat-go-backend/src/auth/token4auth"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/gae_app/debtstracker/api4debtus"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/gae_app/debtstracker/dtdal/gaedal"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/dal4userus"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/dbo4userus"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/facade4userus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/logus"
	"io"
	"net/http"
	"strconv"
	"strings"

	"context"
)

//func getApiUser(ctx context.Context, w http.ResponseWriter, r *http.Request, authInfo auth.AuthInfo) (user models4debtus.AppUser, err error) {
//	if user.ContactID = UserID(ctx, w, r, authInfo); user.ContactID == "" {
//		w.WriteHeader(http.StatusUnauthorized)
//		return
//	}
//
//	if user, err = dal4userus.GetUserByID(c, nil, user.ID); api4debtus.HasError(c, w, err, models4debtus.AppUserKind, user.ContactID, 0) {
//		return
//	} else if user.Data == nil {
//		_, _ = w.Write([]byte(fmt.Sprintf("UserEntry not found by UserID=%v", user.ID)))
//		http.NotFound(w, r) // TODO: Check response output
//		return
//	}
//	return
//}

func HandleUserInfo(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	if userID, err := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write(([]byte)(err.Error()))
	} else {
		if _, err = facade4userus.SaveUserBrowser(ctx, strconv.FormatInt(userID, 10), r.UserAgent()); err != nil {
			logus.Errorf(ctx, err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write(([]byte)(err.Error()))
		}
	}
}

func HandleSaveVisitorData(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		api4debtus.ErrorAsJson(ctx, w, http.StatusBadRequest, err)
		return
	}
	gaClientId := r.FormValue("gaClientId")
	if gaClientId == "" {
		w.WriteHeader(http.StatusBadRequest)
		api4debtus.ErrorAsJson(ctx, w, http.StatusBadRequest, errors.New("missing required parameter gaClientId"))
		return
	}

	userAgent := r.UserAgent()
	ipAddress := strings.SplitN(r.RemoteAddr, ":", 1)[0]

	if _, err := facade4userus.SaveGaClient(ctx, gaClientId, userAgent, ipAddress); err != nil {
		api4debtus.ErrorAsJson(ctx, w, http.StatusInternalServerError, err)
		return
	}
}

func HandleMe(ctx context.Context, w http.ResponseWriter, r *http.Request, authInfo token4auth.AuthInfo, user dbo4userus.UserEntry) {
	api4debtus.ErrorAsJson(ctx, w, http.StatusInternalServerError, errors.New("not implemented"))
	//meDto := dto4debtus.UserMeDto{
	//	UserID:   authInfo.UserID,
	//	FullName: user.Data.GetFullName(),
	//}
	//if ua, err := user.Data.GetAccount("google", ""); err != nil {
	//	api4debtus.ErrorAsJson(ctx, w, http.StatusInternalServerError, err)
	//	return
	//} else if ua != nil {
	//	meDto.GoogleUserID = ua.ContactID
	//}
	//
	//if fbAccounts, err := user.Data.GetAccounts("facebook"); err != nil {
	//	api4debtus.ErrorAsJson(ctx, w, http.StatusInternalServerError, err)
	//	return
	//} else {
	//	for _, ua := range fbAccounts {
	//		meDto.FbUserID = ua.ContactID
	//		break // TODO: change to return map of IDs.
	//	}
	//}
	//
	//if meDto.FullName == models4debtus.NoName {
	//	meDto.FullName = ""
	//}
	//
	//api4debtus.JsonToResponse(ctx, w, meDto)
}

func SetUserName(ctx context.Context, w http.ResponseWriter, r *http.Request, authInfo token4auth.AuthInfo) {

	body, err := io.ReadAll(r.Body)

	if err != nil {
		api4debtus.ErrorAsJson(ctx, w, http.StatusInternalServerError, err)
		return
	}

	if len(body) == 0 {
		api4debtus.ErrorAsJson(ctx, w, http.StatusBadRequest, fmt.Errorf("%w: UserEntry name is required", api4auth.ErrBadRequest))
		return
	}

	userCtx := facade.NewUserContext(authInfo.UserID)
	err = dal4userus.RunUserWorker(ctx, userCtx, true, func(ctx context.Context, tx dal.ReadwriteTransaction, params *dal4userus.UserWorkerParams) error {
		params.User.Data.Names.UserName = string(body)
		params.UserUpdates = append(params.UserUpdates, dal.Update{
			Field: "names.userName",
			Value: params.User.Data.Names.UserName,
		})
		if err = gaedal.DelayUpdateTransfersWithCreatorName(ctx, params.User.ID); err != nil {
			return err
		}
		return err
	})

	if err != nil {
		api4debtus.ErrorAsJson(ctx, w, http.StatusInternalServerError, err)
		return
	}
}
