package emailing

import "github.com/strongo/delaying"

func InitDelaying(mustRegisterFunc func(key string, i any) delaying.Function) {
	delayEmail = mustRegisterFunc(SendEmailTaskCode, delayedSendEmail)
}

var delayEmail delaying.Function
