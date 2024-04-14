package dto4teamus

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/teamus/models4teamus"
)

type teamRecord struct {
	ID  string                `json:"id"`
	Dto models4teamus.TeamDbo `json:"dto"`
}

// TeamResponse response
type TeamResponse struct {
	Team teamRecord `json:"team"`
}
