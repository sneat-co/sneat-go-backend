package cmds4anybot

import (
	"context"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-core/facade"
	"net/url"
)

//func TransactionalCallbackCommand(c botsfw.Command, o db.RunOptions) botsfw.Command {
//	c.CallbackAction = TransactionalCallbackAction(o, c.CallbackAction)
//	return c
//}

func TransactionalCallbackAction(
	f func(whc botsfw.WebhookContext, callbackUrl *url.URL) (m botsfw.MessageFromBot, err error),
) func(whc botsfw.WebhookContext, callbackUrl *url.URL) (m botsfw.MessageFromBot, err error) {
	return func(whc botsfw.WebhookContext, callbackUrl *url.URL) (m botsfw.MessageFromBot, err error) {
		ctx := whc.Context()
		err = facade.RunReadwriteTransaction(ctx, func(tctx context.Context, tx dal.ReadwriteTransaction) error {
			whc.SetContext(tctx)
			m, err = f(whc, callbackUrl)
			whc.SetContext(ctx)
			return err
		})
		return
	}
}
