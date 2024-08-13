package delays4contactus

import (
	"context"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/const4contactus"
	"github.com/strongo/delaying"
	"time"
)

func InitDelays4contactus(mustRegisterFunc func(key string, i any) delaying.Function) {
	delayerUpdateContactusSpaceDboWithContact = mustRegisterFunc("delayedUpdateContactusSpaceDboWithContact", delayedUpdateContactusSpaceDboWithContact)
}

var (
	delayerUpdateContactusSpaceDboWithContact delaying.Function
)

func DelayUpdateContactusSpaceDboWithContact(ctx context.Context, delay time.Duration, userID string, contactID string) error {
	return delayerUpdateContactusSpaceDboWithContact.EnqueueWork(ctx, delaying.With(const4contactus.QueueContacts, "UpdateContactusSpaceDboWithContact", delay), userID, contactID)
}
