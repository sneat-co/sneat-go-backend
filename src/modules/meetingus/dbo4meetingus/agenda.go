package dbo4meetingus

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/spaceus/dbo4spaceus"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
)

// Agenda record
type Agenda struct {
	dbmodels.WithUserIDs
	Status string                          `json:"status" firestore:"status"` // active, archived
	Title  string                          `json:"title" firestore:"title"`
	Spaces []*dbo4spaceus.SpaceMeetingInfo `json:"api4meetingus" firestore:"api4meetingus"`
}

// AgendaTopic record
type AgendaTopic struct {
	ID    string `json:"id" firestore:"id"`
	Title string `json:"title" firestore:"title"`
}
