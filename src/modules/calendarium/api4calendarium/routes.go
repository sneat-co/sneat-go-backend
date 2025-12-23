package api4calendarium

import (
	"net/http"

	"github.com/sneat-co/sneat-go-core/extension"
)

// RegisterHttpRoutes register calendarium routes
func RegisterHttpRoutes(handle extension.HTTPHandleFunc) {
	handle(http.MethodPost, "/v0/happenings/create_happening", httpPostCreateHappening)
	handle(http.MethodPost, "/v0/happenings/update_happening_texts", httpUpdateHappeningTexts)
	handle(http.MethodDelete, "/v0/happenings/delete_happening", httpDeleteHappening)
	handle(http.MethodDelete, "/v0/happenings/delete_slot", httpDeleteSlot)
	handle(http.MethodPost, "/v0/happenings/cancel_happening", httpCancelHappening)
	handle(http.MethodPost, "/v0/happenings/revoke_happening_cancellation", httpRevokeHappeningCancellation)
	handle(http.MethodPost, "/v0/happenings/add_participants", httpAddParticipantsToHappening)
	handle(http.MethodPost, "/v0/happenings/remove_participants", httpRemoveParticipantsFromHappening)

	handle(http.MethodPost, "/v0/happenings/add_slot", httpAddSlot)
	handle(http.MethodPost, "/v0/happenings/update_slot", httpUpdateSlot)

	//  temporary changes slot (for example, time changed for a specific date, or first class has been canceled)
	handle(http.MethodPost, "/v0/happenings/adjust_slot", httpAdjustSlot)
	handle(http.MethodPost, "/v0/happenings/cancel_adjustment", httpCancelAdjustment)

	handle(http.MethodPost, "/v0/happenings/set_prices", httpSetHappeningPrices)
	handle(http.MethodPost, "/v0/happenings/delete_prices", httpDeleteHappeningPrices)
}
