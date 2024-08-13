package unsorted

import (
	"testing"

	"context"
)

func TestApiUserInfo(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("The code did not panic")
		}
	}()

	c := context.Background()
	HandleUserInfo(c, nil, nil)
}
