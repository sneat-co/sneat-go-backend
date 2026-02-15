package dal4listus

import (
	"github.com/dal-go/dalgo/record"
	"github.com/sneat-co/sneat-go-backend/pkg/extensions/listus/dbo4listus"
)

type ListusSpaceEntry = record.DataWithID[string, *dbo4listus.ListusSpaceDbo]
