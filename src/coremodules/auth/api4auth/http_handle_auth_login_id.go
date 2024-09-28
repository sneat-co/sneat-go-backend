package api4auth

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/coremodules/auth/models4auth"
	"github.com/sneat-co/sneat-go-backend/src/coremodules/auth/token4auth"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/common4debtus"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/gae_app/debtstracker/api4debtus"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/gae_app/debtstracker/dtdal"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/logus"
	"io"
	"net/http"
	"strings"
)

func HandleAuthLoginId(ctx context.Context, w http.ResponseWriter, r *http.Request, authInfo token4auth.AuthInfo) {
	query := r.URL.Query()
	channel := query.Get("channel")
	var (
		loginID int
		err     error
	)

	loginIdStr := query.Get("id")

	if loginIdStr != "" {
		if loginID, err = common4debtus.DecodeIntID(loginIdStr); err != nil {
			api4debtus.BadRequestError(ctx, w, err)
			return
		}
	}

	returnLoginID := func(loginID int) {
		encoded := common4debtus.EncodeIntID(loginID)
		logus.Infof(ctx, "Login ContactID: %d, Encoded: %s", loginID, encoded)
		if _, err = w.Write([]byte(encoded)); err != nil {
			logus.Criticalf(ctx, "Failed to write login ContactID to response: %v", err)
		}
	}

	if loginID != 0 {
		if loginPin, err := dtdal.LoginPin.GetLoginPinByID(ctx, nil, loginID); err != nil {
			if dal.IsNotFound(err) {
				w.WriteHeader(http.StatusInternalServerError)
				logus.Errorf(ctx, err.Error())
				return
			}
		} else if loginPin.Data.IsActive(channel) {
			returnLoginID(loginID)
			return
		}
	}

	var rBody []byte
	if rBody, err = io.ReadAll(r.Body); err != nil {
		api4debtus.BadRequestError(ctx, w, fmt.Errorf("failed to read request body: %w", err))
		return
	}
	gaClientID := string(rBody)

	if gaClientID != "" {
		if len(gaClientID) > 100 {
			api4debtus.BadRequestMessage(ctx, w, fmt.Sprintf("Google Client ContactID is too long: %d", len(gaClientID)))
			return
		}

		if strings.Count(gaClientID, ".") != 1 {
			api4debtus.BadRequestMessage(ctx, w, "Google Client ContactID has wrong format, a '.' char expected")
			return
		}
	}

	err = facade.RunReadwriteTransaction(ctx, func(ctx context.Context, tx dal.ReadwriteTransaction) (err error) {
		var loginPin models4auth.LoginPin
		if loginPin, err = dtdal.LoginPin.CreateLoginPin(ctx, tx, channel, gaClientID, authInfo.UserID); err != nil {
			api4debtus.ErrorAsJson(ctx, w, http.StatusInternalServerError, err)
			return
		}
		loginID = loginPin.ID
		return err
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logus.Errorf(ctx, err.Error())
		return
	}
	returnLoginID(loginID)
}
