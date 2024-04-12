package models4logist

import (
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	"github.com/sneat-co/sneat-go-backend/src/modules/teamus/dal4teamus"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/strongo/validation"
	"strings"
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

// LogistTeamDto is a DTO for LogistTeam
type LogistTeamDto struct {
	dbmodels.WithUserIDs
	Roles             []string
	ContactID         string `json:"contactID,omitempty" firestore:"contactID,omitempty"`
	OrderNumberPrefix string `json:"orderNumberPrefix,omitempty" firestore:"orderNumberPrefix,omitempty"`
	//
	OrderCounters map[string]OrderCounter `json:"orderCounters,omitempty" firestore:"orderCounters,omitempty"`
}

// Validate returns error if invalid
func (v LogistTeamDto) Validate() error {
	if err := v.WithUserIDs.Validate(); err != nil {
		return validation.NewErrBadRecordFieldValue("WithUserIDs", err.Error())
	}
	if len(v.Roles) == 0 {
		return validation.NewErrRecordIsMissingRequiredField("roles")
	}
	for i, role := range v.Roles {
		if !IsKnownLogistCompanyRole(LogistTeamRole(role)) {
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

// LogistTeamContext is a context for LogistTeam
type LogistTeamContext struct {
	record.WithID[string]
	Dto *LogistTeamDto
}

func newLogistTeamKey(teamID string) *dal.Key {
	teamKey := dal4teamus.NewTeamKey(teamID)
	return dal.NewKeyWithParentAndID(teamKey, dal4teamus.TeamModulesCollection, ModuleID)
}

// NewLogistTeamContext creates new LogistTeamContext
func NewLogistTeamContext(teamID string) (logistTeam LogistTeamContext) {
	logistTeam.ID = teamID
	logistTeam.Key = newLogistTeamKey(teamID)
	logistTeam.Dto = new(LogistTeamDto)
	logistTeam.Record = dal.NewRecordWithData(logistTeam.Key, logistTeam.Dto)
	return logistTeam
}