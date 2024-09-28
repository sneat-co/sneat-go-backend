package facade4auth

import (
	"fmt"
	"testing"

	"context"
)

func TestNewUserFacebookKey(t *testing.T) {
	const (
		fbAppID  = "12345"
		fbUserID = "975"
	)
	key := NewUserFacebookKey(context.Background(), fbAppID, fbUserID)
	testStringKey(t, fmt.Sprintf("%s:%s", fbAppID, fbUserID), key)
}
