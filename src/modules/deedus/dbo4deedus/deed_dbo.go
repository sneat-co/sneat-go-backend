package dbo4deedus

import (
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/strongo/strongoapp/with"
)

type DeedDetails struct {
	Text string `json:"text" firestore:"text"`
}

type DeedDbo struct {
	with.CreatedFields
	Status  dbmodels.Status `json:"status" firestore:"status"`
	Starts  int             `json:"starts" firestore:"starts"`
	Details DeedDetails     `json:"details" firestore:"details"`
}
