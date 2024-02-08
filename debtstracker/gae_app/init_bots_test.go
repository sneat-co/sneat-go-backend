package gaeapp

import (
	"testing"
)

func TestInitBot(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Should fail")
		}
	}()
	InitBots(nil, nil, nil /*common.DebtusAppContext{}*/)
}
