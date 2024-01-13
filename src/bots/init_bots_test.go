package bots

import (
	"fmt"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/julienschmidt/httprouter"
	"strings"
	"testing"
)

func TestInitializeBots(t *testing.T) {
	type args struct {
		botHost    botsfw.BotHost
		httpRouter *httprouter.Router
	}
	tests := []struct {
		name         string
		args         args
		expectsPanic string
	}{
		{
			name:         "nil_args",
			args:         args{},
			expectsPanic: "== nil",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.expectsPanic != "" {
				defer func() {
					if r := recover(); r == nil {
						t.Error("Expected panic: " + tt.expectsPanic)
					} else if p := fmt.Sprintf("%v", r); !strings.Contains(p, tt.expectsPanic) {
						t.Errorf("Expected panic '%s', got: %s" + tt.expectsPanic)
					}
				}()
			}
			InitializeBots(tt.args.botHost, tt.args.httpRouter)
		})
	}
}
