package admin

import (
	"context"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/models4debtus"
	"testing"
)

func TestSendFeedbackToAdmins(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("should panic")
		}
	}()
	_ = SendFeedbackToAdmins(context.Background(), "", models4debtus.Feedback{})
}
