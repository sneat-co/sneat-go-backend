package dto4spaceus

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/spaceus/dbo4spaceus"
)

type spaceRecord struct {
	ID  string               `json:"id"`
	Dto dbo4spaceus.SpaceDbo `json:"dto4debtus"`
}

// SpaceResponse response
type SpaceResponse struct {
	Space spaceRecord `json:"space"`
}
