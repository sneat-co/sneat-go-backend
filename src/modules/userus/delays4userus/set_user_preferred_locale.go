package delays4userus

import (
	"context"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/facade4userus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/logus"
)

func delayedSetUserPreferredLocale(c context.Context, userID string, localeCode5 string) (err error) {
	logus.Debugf(c, "delayedSetUserPreferredLocale(userID=%v, localeCode5=%v)", userID, localeCode5)
	userContext := facade.NewUserContext(userID)
	return facade4userus.SetUserPreferredLocale(c, userContext, localeCode5)
}