package delays4userus

import (
	"context"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/const4userus"
	"github.com/strongo/delaying"
	"time"
)

var (
	delaySetUserPreferredLocale delaying.Function
)

func InitDelays4userus(mustRegisterFunc func(key string, i any) delaying.Function) {
	delaySetUserPreferredLocale = mustRegisterFunc("delayedSetUserPreferredLocale", delayedSetUserPreferredLocale)
}

func DelaySetUserPreferredLocale(ctx context.Context, delay time.Duration, userID string, locale string) error {
	return delaySetUserPreferredLocale.EnqueueWork(ctx, delaying.With(const4userus.QueueUsers, "set-user-preferred-locale", delay), userID, locale)
}
