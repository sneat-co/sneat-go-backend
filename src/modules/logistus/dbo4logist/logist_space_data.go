package dbo4logist

import (
	"fmt"
	"strings"

	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	dbo4spaceus2 "github.com/sneat-co/sneat-core-modules/spaceus/dbo4spaceus"
	"github.com/sneat-co/sneat-go-backend/src/modules/logistus/const4logistus"
	"github.com/sneat-co/sneat-go-core/coretypes"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/strongo/validation"
)

// OrderCounter is a counter for orders
type OrderCounter struct {
	Prefix     string `json:"prefix" firestore:"prefix"`
	LastNumber int    `json:"lastNumber" firestore:"lastNumber"`
}

// Validate returns nil if valid, or error if not
func (v OrderCounter) Validate() error {
	//if strings.TrimSpace(v.Prefix) == "" {
	//	return validation.NewErrRecordIsMissingRequiredField("prefix")
	//}
	if v.LastNumber < 0 {
		return validation.NewErrBadRecordFieldValue("lastNumber", "should be positive integer")
	}
	return nil
}

// LogistSpaceDbo is a DTO for LogistTeam
type LogistSpaceDbo struct {
	dbmodels.WithUserIDs
	Roles             []string
	ContactID         string `json:"contactID,omitempty" firestore:"contactID,omitempty"`
	OrderNumberPrefix string `json:"orderNumberPrefix,omitempty" firestore:"orderNumberPrefix,omitempty"`
	//
	OrderCounters map[string]OrderCounter `json:"orderCounters,omitempty" firestore:"orderCounters,omitempty"`
}

// Validate returns error if invalid
func (v LogistSpaceDbo) Validate() error {
	if err := v.WithUserIDs.Validate(); err != nil {
		return validation.NewErrBadRecordFieldValue("WithUserIDs", err.Error())
	}
	if len(v.Roles) == 0 {
		return validation.NewErrRecordIsMissingRequiredField("roles")
	}
	for i, role := range v.Roles {
		if !IsKnownLogistCompanyRole(LogistSpaceRole(role)) {
			return validation.NewErrBadRequestFieldValue(fmt.Sprintf("roles[%d]", i),
				fmt.Sprintf("should be one of: %+v", KnownLogistCompanyRoles))
		}
	}
	if strings.TrimSpace(v.ContactID) != v.ContactID {
		return validation.NewErrBadRecordFieldValue("contactID", "should be trimmed")
	}
	if strings.TrimSpace(v.OrderNumberPrefix) != v.OrderNumberPrefix {
		return validation.NewErrBadRecordFieldValue("orderNumberPrefix", "should be trimmed")
	}
	if len(v.OrderNumberPrefix) > 5 {
		return validation.NewErrBadRecordFieldValue("vatNumber", "should not be longer than 5 characters")
	}
	for name, counter := range v.OrderCounters {
		if err := counter.Validate(); err != nil {
			return validation.NewErrBadRecordFieldValue("orderCounters."+name, err.Error())
		}
	}
	return nil
}

// LogistSpaceEntry is a context for LogistTeam
type LogistSpaceEntry = record.DataWithID[coretypes.ExtID, *LogistSpaceDbo]

func newLogistSpaceKey(spaceID coretypes.SpaceID) *dal.Key {
	key := dbo4spaceus2.NewSpaceKey(spaceID)
	return dal.NewKeyWithParentAndID(key, dbo4spaceus2.SpaceModulesCollection, const4logistus.ExtensionID)
}

// NewLogistSpaceEntry creates new LogistSpaceEntry
func NewLogistSpaceEntry(spaceID coretypes.SpaceID) (logistSpace LogistSpaceEntry) {
	logistSpace.ID = const4logistus.ExtensionID
	logistSpace.Key = newLogistSpaceKey(spaceID)
	logistSpace.Data = new(LogistSpaceDbo)
	logistSpace.Record = dal.NewRecordWithData(logistSpace.Key, logistSpace.Data)
	return logistSpace
}
