package facade4calendarius

import (
	"strings"
	"time"

	"github.com/sneat-co/sneat-go-backend/pkg/extensions/calendarius/dbo4calendarius"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
)

func CreateCancellation(uid, reason string) dbo4calendarius.Cancellation {
	return dbo4calendarius.Cancellation{
		At:     time.Now(),
		By:     dbmodels.ByUser{UID: uid},
		Reason: strings.TrimSpace(reason),
	}
}
