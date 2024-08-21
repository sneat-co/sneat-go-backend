package models4auth

import (
	"github.com/dal-go/dalgo/record"
	"time"
)

const GaClientKind = "UserGaClient"

type GaClientEntity struct {
	Created   time.Time
	UserAgent string `firestore:",omitempty"`
	IpAddress string `firestore:",omitempty"`
}

type GaClient struct {
	record.WithID[string]
	*GaClientEntity
}

func (GaClient) Kind() string {
	return GaClientKind
}

func (gaClient GaClient) Entity() interface{} {
	return gaClient.GaClientEntity
}

func (GaClient) NewEntity() interface{} {
	return new(GaClientEntity)
}

func (gaClient *GaClient) SetEntity(entity interface{}) {
	if entity == nil {
		gaClient.GaClientEntity = nil

	} else {
		gaClient.GaClientEntity = entity.(*GaClientEntity)

	}
}

//var _ db.EntityHolder = (*GaClient)(nil)
