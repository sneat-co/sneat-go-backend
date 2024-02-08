package main

import (
	"testing"

	"github.com/strongo/log"
)

func TestInit(t *testing.T) {
	t.Skip("TODO: fix")
	if log.NumberOfLoggers() == 0 {
		t.Error("At least 1 logger should be added")
	}
}
