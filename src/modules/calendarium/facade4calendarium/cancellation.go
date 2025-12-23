package facade4calendarium

import (
	"strings"
	"time"

	"github.com/sneat-co/sneat-go-backend/src/modules/calendarium/dbo4calendarium"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
)

func CreateCancellation(uid, reason string) dbo4calendarium.Cancellation {
	return dbo4calendarium.Cancellation{
		At:     time.Now(),
		By:     dbmodels.ByUser{UID: uid},
		Reason: strings.TrimSpace(reason),
	}
}
