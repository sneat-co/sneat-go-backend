package api

import "github.com/strongo/delaying"

func InitDelaying(mustRegisterFunc func(key string, i any) delaying.Function) {
	delayChangeTransfersCounterparty = mustRegisterFunc("changeTransfersCounterparty", delayedChangeTransfersCounterparty)
	delayChangeTransferCounterparty = mustRegisterFunc("changeTransferCounterparty", delayedChangeTransferCounterparty)
}

var (
	delayChangeTransfersCounterparty,
	delayChangeTransferCounterparty delaying.Function
)
