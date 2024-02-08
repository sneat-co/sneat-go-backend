package splitus

import "github.com/strongo/delaying"

func InitDelaying(mustRegisterFunc func(key string, i any) delaying.Function) {
	delayUpdateBillCards = mustRegisterFunc("UpdateBillCards", delayedUpdateBillCards)
	delayUpdateBillTgChatCard = mustRegisterFunc("UpdateBillTgChatCard", delayedUpdateBillTgChartCard)
}

var (
	delayUpdateBillCards      delaying.Function
	delayUpdateBillTgChatCard delaying.Function
)
