package facade4listus

import (
	"context"

	"github.com/sneat-co/sneat-go-core/facade"
)

func ClearList(ctx context.Context, userCtx facade.UserContext, listID string) {
	_, _, _ = ctx, userCtx, listID
}
