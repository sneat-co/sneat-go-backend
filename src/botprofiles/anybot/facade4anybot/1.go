package facade4anybot

import (
	"github.com/strongo/delaying"
)

func InitDelaying(mustRegisterFunc func(key string, i any) delaying.Function) {
	//delayerSetUserReferrer = mustRegisterFunc("delayedSetUserReferrer", delayedSetUserReferrer)
}

//var delayerSetUserReferrer delaying.Function
