package facade4splitus

import (
	"context"
	"github.com/sneat-co/sneat-go-backend/src/modules/splitus/const4splitus"
	"github.com/strongo/delaying"
)

func DelayUpdateSpaceWithBill(c context.Context, userID string, billID string) (err error) {
	if err = delayerUpdateUserWithBill.EnqueueWork(c, delaying.With(const4splitus.QueueSplitus, "UpdateUserWithBill", 0), userID, billID); err != nil {
		return
	}
	return
}
