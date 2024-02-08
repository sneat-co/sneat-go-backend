package unsorted

import (
	"github.com/strongo/delaying"
)

func InitDelaying(mustRegisterFunc func(key string, i any) delaying.Function) {
	delayChangeTransfersCounterparty = mustRegisterFunc("changeTransfersCounterparty", DelayedChangeTransfersCounterparty)
	delayChangeTransferCounterparty = mustRegisterFunc("changeTransferCounterparty", DelayedChangeTransferCounterparty)
}

var (
	delayChangeTransfersCounterparty,
	delayChangeTransferCounterparty delaying.Function
)
