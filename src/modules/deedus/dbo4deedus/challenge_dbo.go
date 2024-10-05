package dbo4deedus

import "github.com/sneat-co/sneat-go-core/models/dbmodels"

type ChallengeDbo struct {
	Title  string          `json:"title" firestore:"title"`
	Status dbmodels.Status `json:"status"  firestore:"status"`
	Stars  int             `json:"stars"  firestore:"stars"`
}
