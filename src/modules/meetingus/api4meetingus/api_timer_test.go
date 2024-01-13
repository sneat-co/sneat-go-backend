package api4meetingus

import (
	"context"
	"github.com/sneat-co/sneat-go-backend/src/modules/meetingus/facade4meetingus"
	"github.com/sneat-co/sneat-go-core/facade"
	"testing"
)

func TestToggleMeetingTimer(t *testing.T) {
	toggleTimer = func(ctx context.Context, userContext facade.User, params facade4meetingus.ToggleParams) (response facade4meetingus.ToggleTimerResponse, err error) {
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
	toggleTimer = func(ctx context.Context, userContext facade.User, params facade4meetingus.ToggleParams) (response facade4meetingus.ToggleTimerResponse, err error) {
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
