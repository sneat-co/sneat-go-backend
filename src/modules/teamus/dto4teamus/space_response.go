package dto4teamus

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/teamus/dbo4teamus"
)

type spaceRecord struct {
	ID  string              `json:"id"`
	Dto dbo4teamus.SpaceDbo `json:"dto"`
}

// SpaceResponse response
type SpaceResponse struct {
	Space spaceRecord `json:"space"`
}
