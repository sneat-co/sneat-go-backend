package api4meetingus

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/meetingus/facade4meetingus"
	"github.com/sneat-co/sneat-go-core/apicore"
	"github.com/sneat-co/sneat-go-core/apicore/verify"
	"github.com/strongo/validation"
	"net/http"
	"strings"
)

var toggleTimer = facade4meetingus.ToggleTimer

// ToggleMeetingTimer switches api4meetingus timer
func ToggleMeetingTimer(params facade4meetingus.Params) func(w http.ResponseWriter, r *http.Request) {
	return toggleTimerEndpoint(params, nil)
}

// ToggleMemberTimer switches person timer
func ToggleMemberTimer(params facade4meetingus.Params) func(w http.ResponseWriter, r *http.Request) {
	return toggleTimerEndpoint(params, func(request facade4meetingus.ToggleTimerRequest) error {
		if strings.TrimSpace(request.Member) == "" {
			return validation.NewErrRecordIsMissingRequiredField("members")
		}
		return nil
	})
}

func toggleTimerEndpoint(params facade4meetingus.Params, requestValidator func(request facade4meetingus.ToggleTimerRequest) error) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, err := apicore.VerifyRequestAndCreateUserContext(w, r, verify.DefaultJsonWithAuthRequired)
		if err != nil {
			return
		}
		var request facade4meetingus.ToggleTimerRequest
		if err = apicore.DecodeRequestBody(w, r, &request); err != nil {
			return
		}
		if requestValidator != nil {
			if err = requestValidator(request); err != nil {
				apicore.ReturnError(ctx, w, r, err)
			}
		}
		response, err := toggleTimer(ctx, facade4meetingus.ToggleParams{Params: params, Request: request})
		if err == nil {
			if err = response.Validate(); err != nil {
				apicore.ReturnError(ctx, w, r, err)
			}
		}
		apicore.ReturnJSON(ctx, w, r, http.StatusOK, err, response)
	}
}
