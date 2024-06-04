package dto4teamus

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/teamus/dbo4teamus"
)

type teamRecord struct {
	ID  string             `json:"id"`
	Dto dbo4teamus.TeamDbo `json:"dto"`
}

// TeamResponse response
type TeamResponse struct {
	Team teamRecord `json:"team"`
}
