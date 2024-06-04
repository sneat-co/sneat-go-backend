package dal4listus

import (
	"github.com/dal-go/dalgo/record"
	"github.com/sneat-co/sneat-go-backend/src/modules/listus/dbo4listus"
)

type ListusTeamContext = record.DataWithID[string, *dbo4listus.ListusTeamDto]
