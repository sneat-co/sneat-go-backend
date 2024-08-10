package unsorted

import (
	"errors"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/api"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/facade2debtus/dto"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/logus"
	"io"
	"net/http"
	"strconv"
	"strings"

	"context"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/auth"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/dtdal"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/dtdal/gaedal"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/facade2debtus"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/models"
)

func getApiUser(c context.Context, w http.ResponseWriter, r *http.Request, authInfo auth.AuthInfo) (user models.AppUser, err error) {
	if user.ID = getUserID(c, w, r, authInfo); user.ID == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if user, err = facade2debtus.User.GetUserByID(c, nil, user.ID); api.HasError(c, w, err, models.AppUserKind, user.ID, 0) {
		return
	} else if user.Data == nil {
		_, _ = w.Write([]byte(fmt.Sprintf("User not found by ID=%v", user.ID)))
		http.NotFound(w, r) // TODO: Check response output
		return
	}
	return
}

func HandleUserInfo(c context.Context, w http.ResponseWriter, r *http.Request) {
	if userID, err := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write(([]byte)(err.Error()))
	} else {
		if err := SaveUserAgent(c, strconv.FormatInt(userID, 10), r.UserAgent()); err != nil {
			logus.Errorf(c, err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write(([]byte)(err.Error()))
		}
	}
}

func SaveUserAgent(c context.Context, userID string, userAgent string) error {
	userAgent = strings.TrimSpace(userAgent)
	if userAgent == "" {
		return nil
	}
	_, err := dtdal.UserBrowser.SaveUserBrowser(c, userID, userAgent)
	return err
}

func HandleSaveVisitorData(c context.Context, w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		api.ErrorAsJson(c, w, http.StatusBadRequest, err)
		return
	}
	gaClientId := r.FormValue("gaClientId")
	if gaClientId == "" {
		w.WriteHeader(http.StatusBadRequest)
		api.ErrorAsJson(c, w, http.StatusBadRequest, errors.New("missing required parameter gaClientId"))
		return
	}

	userAgent := r.UserAgent()
	ipAddress := strings.SplitN(r.RemoteAddr, ":", 1)[0]

	if _, err := dtdal.UserGaClient.SaveGaClient(c, gaClientId, userAgent, ipAddress); err != nil {
		api.ErrorAsJson(c, w, http.StatusInternalServerError, err)
		return
	}
}

func HandleMe(c context.Context, w http.ResponseWriter, r *http.Request, authInfo auth.AuthInfo, user models.AppUser) {
	meDto := dto.UserMeDto{
		UserID:   authInfo.UserID,
		FullName: user.Data.FullName(),
	}
	if ua, err := user.Data.GetAccount("google", ""); err != nil {
		api.ErrorAsJson(c, w, http.StatusInternalServerError, err)
		return
	} else if ua != nil {
		meDto.GoogleUserID = ua.ID
	}

	if fbAccounts, err := user.Data.GetAccounts("facebook"); err != nil {
		api.ErrorAsJson(c, w, http.StatusInternalServerError, err)
		return
	} else {
		for _, ua := range fbAccounts {
			meDto.FbUserID = ua.ID
			break // TODO: change to return map of IDs.
		}
	}

	if meDto.FullName == models.NoName {
		meDto.FullName = ""
	}

	api.JsonToResponse(c, w, meDto)
}

func SetUserName(c context.Context, w http.ResponseWriter, r *http.Request, authInfo auth.AuthInfo) {

	body, err := io.ReadAll(r.Body)

	if err != nil {
		api.ErrorAsJson(c, w, http.StatusInternalServerError, err)
		return
	}

	if len(body) == 0 {
		api.ErrorAsJson(c, w, http.StatusBadRequest, fmt.Errorf("%w: User name is required", ErrBadRequest))
		return
	}

	err = facade.RunReadwriteTransaction(c, func(c context.Context, tx dal.ReadwriteTransaction) error {
		user, err := facade2debtus.User.GetUserByID(c, tx, authInfo.UserID)
		if err != nil {
			return err
		}
		user.Data.Username = string(body)
		if err = facade2debtus.User.SaveUser(c, tx, user); err != nil {
			return err
		}
		if err = gaedal.DelayUpdateTransfersWithCreatorName(c, user.ID); err != nil {
			return err
		}
		return err
	})

	if err != nil {
		api.ErrorAsJson(c, w, http.StatusInternalServerError, err)
		return
	}
}
