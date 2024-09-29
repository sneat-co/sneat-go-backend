package facade4debtus

import (
	"github.com/sneat-co/sneat-core-modules/anybot/facade4anybot"
	"github.com/strongo/delaying"
)

func init() {
	delaying.Init(delaying.VoidWithLog)
	facade4anybot.InitDelaying(delaying.MustRegisterFunc)
}
