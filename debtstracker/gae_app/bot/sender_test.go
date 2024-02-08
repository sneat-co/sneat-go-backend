package bot

import (
	"github.com/bots-go-framework/bots-fw/botsfw"
	"testing"
)

func TestSendRefreshOrNothingChanged(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("should panic")
		}
	}()

	var m botsfw.MessageFromBot
	_, err := SendRefreshOrNothingChanged(nil, m)
	if err != nil {
		t.Fatal(err)
	}
}
