package delays4contactus

import (
	"fmt"
	"golang.org/x/net/context"
)

func delayedUpdateContactusSpaceDboWithContact(_ context.Context, userID string, contactID string) (err error) {
	return fmt.Errorf("not implemented: userID=%s, contactID=%s", userID, contactID)
}
