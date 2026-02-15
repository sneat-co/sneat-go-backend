package api4meetingus

import (
	"testing"

	"github.com/sneat-co/sneat-go-backend/pkg/extensions/meetingus/facade4meetingus"
	"github.com/sneat-co/sneat-go-core/facade"
)

func TestToggleMeetingTimer(t *testing.T) {
	toggleTimer = func(ctx facade.ContextWithUser, params facade4meetingus.ToggleParams) (response facade4meetingus.ToggleTimerResponse, err error) {
		return
	}
	params := facade4meetingus.Params{
		RecordFactory: nil,
		BeforeSafe:    nil,
	}
	handler := ToggleMeetingTimer(params)
	if handler == nil {
		t.Fatal("handler = nil")
	}
	// TODO: implement
}

func TestToggleMemberTimer(t *testing.T) {
	toggleTimer = func(ctx facade.ContextWithUser, params facade4meetingus.ToggleParams) (response facade4meetingus.ToggleTimerResponse, err error) {
		return
	}
	params := facade4meetingus.Params{
		RecordFactory: nil,
		BeforeSafe:    nil,
	}
	handler := ToggleMemberTimer(params)
	if handler == nil {
		t.Fatal("handler = nil")
	}
	// TODO: implement
}
