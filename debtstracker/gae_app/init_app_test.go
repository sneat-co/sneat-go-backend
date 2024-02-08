package gaeapp

import (
	"testing"
)

func TestInit_Nil(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("The code did not panic")
		}
	}()

	Init(nil)
}
