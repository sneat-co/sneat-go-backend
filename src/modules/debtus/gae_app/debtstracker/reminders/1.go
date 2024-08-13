package reminders

import "github.com/strongo/delaying"

func InitDelaying(mustRegisterFunc func(key string, i any) delaying.Function) {
	delaySetChatIsForbidden = mustRegisterFunc("SetChatIsForbidden", SetChatIsForbidden)
}

var (
	delaySetChatIsForbidden delaying.Function
)
