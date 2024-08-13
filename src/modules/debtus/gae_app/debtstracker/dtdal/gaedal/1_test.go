package gaedal

import (
	"github.com/strongo/delaying"
)

func init() {
	delaying.Init(delaying.VoidWithLog)
	InitDelaying(delaying.MustRegisterFunc)
}
