package api4sportus

import (
	"github.com/sneat-co/sneat-go-core/modules"
	"net/http"
)

const spotBuddiesPathPrefix = "/v0/spot_buddies/"

func registerSpotHandlers(handle modules.HTTPHandleFunc) {
	//handle(http.MethodGet, spotBuddiesPathPrefix+"spots/my_spots", mySpots)
	handle(http.MethodPost, spotBuddiesPathPrefix+"spots/join_spot", joinSpot)
	handle(http.MethodPost, spotBuddiesPathPrefix+"spots/leave_spot", leaveSpot)
	handle(http.MethodPost, spotBuddiesPathPrefix+"spots/rsvp_to_spot", rsvpToSpot)
	handle(http.MethodPost, spotBuddiesPathPrefix+"spots/checkin_to_spot", checkinToSpot)
	handle(http.MethodPost, spotBuddiesPathPrefix+"spots/checkout_to_spot", checkoutFromSpot)
}

func joinSpot(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func leaveSpot(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

//func mySpots(w http.ResponseWriter, r *http.Request) {
//	w.WriteHeader(http.StatusNotImplemented)
//}

func rsvpToSpot(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func checkinToSpot(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func checkoutFromSpot(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}
