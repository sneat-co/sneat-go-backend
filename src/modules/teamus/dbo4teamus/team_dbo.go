package dbo4teamus

import (
	"errors"
	"fmt"
	"github.com/sneat-co/sneat-go-backend/src/modules/teamus/core4teamus"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/strongo/slice"
	"github.com/strongo/strongoapp/with"
	"github.com/strongo/validation"
	"strings"
	"time"
)

// BoolMetricVal record
type BoolMetricVal struct {
	Label string `json:"label,omitempty" firestore:"label,omitempty"`
	Color string `json:"color,omitempty" firestore:"color,omitempty"`
}

// BoolMetric record
type BoolMetric struct {
	True  *BoolMetricVal `json:"true,omitempty" firestore:"true,omitempty"`
	False *BoolMetricVal `json:"false,omitempty" firestore:"false,omitempty"`
}

// TeamMetric record
type TeamMetric struct {
	ID      string      `json:"id" firestore:"id"`
	Title   string      `json:"title" firestore:"title"`
	Mode    string      `json:"mode" firestore:"mode"` // Possible values: personal, team
	Type    string      `json:"type" firestore:"type"` // Possible values: bool, int, str
	Min     *int        `json:"min,omitempty" firestore:"min,omitempty"`
	Max     *int        `json:"max,omitempty" firestore:"max,omitempty"`
	Bool    *BoolMetric `json:"bool,omitempty" firestore:"bool,omitempty"`
	Options []string    `json:"options" firestore:"options"`
}

// Validate validates TeamMetric record
func (v *TeamMetric) Validate() error {
	if v.Title == "" {
		return validation.NewErrRecordIsMissingRequiredField("title")
	}
	switch v.Mode {
	case "":
		return validation.NewErrRecordIsMissingRequiredField("mode")
	case "personal", "team":
		break
	default:
		return validation.NewErrBadRecordFieldValue("mode", fmt.Sprintf("unknwon value: %s", v.Mode))
	}
	switch v.Type {
	case "":
		return validation.NewErrRecordIsMissingRequiredField("type")
	case "bool", "int", "options":
		break
	default:
		return validation.NewErrBadRecordFieldValue("type", fmt.Sprintf("unknwon value: %s", v.Type))
	}
	return nil
}

// TeamMeetingInfo record
type TeamMeetingInfo struct {
	ID       string     `json:"id,omitempty" firestore:"id,omitempty"`
	Stage    string     `json:"stage" firestore:"stage"`
	Started  *time.Time `json:"started,omitempty" firestore:"started,omitempty"`
	Finished *time.Time `json:"finished,omitempty" firestore:"finished,omitempty"`
}

// EarliestPossibleTime 2020-01-01
var EarliestPossibleTime = time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)

// Validate validates record
func (v *TeamMeetingInfo) Validate() error {
	if v.ID == "" {
		return validation.NewErrRecordIsMissingRequiredField("id")
	}
	if v.Stage == "" {
		return validation.NewErrRecordIsMissingRequiredField("stage")
	}
	if v.Finished != nil {
		if v.Finished.Before(EarliestPossibleTime) {
			return validation.NewErrBadRecordFieldValue("finished", fmt.Sprintf("Too early end date: %v", v.Finished))
		}
		if v.Started == nil {
			return validation.NewErrBadRecordFieldValue("started", "started time should be specified when finished time is presented")
		}
	}
	if v.Started != nil && v.Started.Before(EarliestPossibleTime) {
		return validation.NewErrBadRecordFieldValue("started", fmt.Sprintf("Too early start date: %v", v.Started))
	}
	return nil
}

// TeamMeetings record
type TeamMeetings struct {
	Retrospective *TeamMeetingInfo `json:"retrospective,omitempty" firestore:"retrospective,omitempty"`
	Scrum         *TeamMeetingInfo `json:"scrum,omitempty" firestore:"scrum,omitempty"`
}

// TeamBrief is a base class for TeamDbo
type TeamBrief struct {
	Type   core4teamus.TeamType `json:"type" firestore:"type"`
	Title  string               `json:"title" firestore:"title"`
	Status dbmodels.Status      `json:"status" firestore:"status"`

	Modules []string `json:"modules,omitempty" firestore:"modules,omitempty"`

	with.RequiredCountryID

	// TODO: This should be populated
	ParentTeamID string `json:"parentTeamID,omitempty" firestore:"parentTeamID,omitempty"`
}

func (v TeamBrief) Validate() error {
	v.Title = strings.TrimSpace(v.Title)
	if v.Type == "" {
		return validation.NewErrRecordIsMissingRequiredField("type")
	}
	if v.Type != core4teamus.TeamTypeFamily && v.Title == "" {
		return errors.New("non family team is required to have a title")
	}
	if !core4teamus.IsValidTeamType(v.Type) {
		if v.Title == "" {
			return validation.NewErrBadRecordFieldValue("type", "unknown value")
		}
		return validation.NewErrBadRecordFieldValue("type", "unknown value for team:"+v.Title)
	}
	if v.Status == "" {
		return validation.NewErrRequestIsMissingRequiredField("status")
	} else if !IsKnownTeamStatus(v.Status) {
		return validation.NewErrBadRecordFieldValue("status", "unknown value: "+v.Status)
	}
	if err := v.RequiredCountryID.Validate(); err != nil {
		return err
	}
	return nil
}

func IsKnownTeamStatus(status dbmodels.Status) bool {
	switch status {
	case dbmodels.StatusActive, dbmodels.StatusArchived, dbmodels.StatusDeleted, dbmodels.StatusDraft:
		return true
	}
	return false
}

const NumberOfMembersFieldName = "members"

// TeamDbo record
type TeamDbo struct {
	TeamBrief
	with.CreatedFields
	dbmodels.WithUpdatedAndVersion
	dbmodels.WithUserIDs

	//NumberOf map[string]int `json:"numberOf,omitempty" firestore:"numberOf,omitempty"`
	//
	//Contacts   []*briefs4contactus.ContactBrief `json:"contacts,omitempty" firestore:"contacts,omitempty"`

	//
	Timezone *dbmodels.Timezone `json:"timezone,omitempty" firestore:"timezone,omitempty"`
	//
	Metrics []*TeamMetric `json:"metrics,omitempty" firestore:"metrics,omitempty"`
}

//func (v *TeamDbo) SetNumberOf(name string, value int) (update dal.Update) {
//	if v.NumberOf == nil {
//		v.NumberOf = make(map[string]int)
//	}
//	v.NumberOf[name] = value
//	return dal.Update{
//		Field: "numberOf." + name,
//		Value: value,
//	}
//}

// IncreaseVersion increases record version and sets timestamp
func (v *TeamDbo) IncreaseVersion(timestamp time.Time, updatedBy string) int {
	v.Version++
	v.UpdatedAt = timestamp
	v.UpdatedBy = updatedBy
	return v.Version
}

// Validate validates record
func (v *TeamDbo) Validate() error {
	if err := v.TeamBrief.Validate(); err != nil {
		return err
	}
	if err := v.CreatedFields.Validate(); err != nil {
		return err
	}
	if err := v.WithUpdatedAndVersion.Validate(); err != nil {
		return err
	}
	for i, userID := range v.UserIDs {
		if strings.TrimSpace(userID) == "" {
			return validation.NewErrBadRecordFieldValue("userIDs", fmt.Errorf("empty value at index %v", i).Error())
		}
	}
	if err := v.Timezone.Validate(); err != nil {
		return fmt.Errorf("invalid 'timezone' field: %w", err)
	}
	{ // Validate "numberOf" field
		//validateCounter := func(counterName, briefsFieldName string, briefsCount int) error {
		//	if briefsCount == 0 {
		//		return nil
		//	}
		//	counter := v.NumberOf[counterName]
		//	if counter < 0 {
		//		return validation.NewErrBadRecordFieldValue("numberOf."+counterName,
		//			fmt.Sprintf("should be positive, got: %d", counter))
		//	}
		//	if briefsCount != counter {
		//		return validation.NewErrBadRecordFieldValue("numberOf."+counterName,
		//			fmt.Sprintf("%v does not match number of items in '%v' field: %v", counter, briefsFieldName, briefsCount))
		//	}
		//	return nil
		//}

		//for name, value := range v.NumberOf {
		//	if value < 0 {
		//		return validation.NewErrBadRecordFieldValue("numberOf."+name,
		//			fmt.Sprintf("should be positive, got: %d", value))
		//	}
		//}
	}

	return nil
}

// HasUser checks if team has a user with given ID
func (v *TeamDbo) HasUser(uid string) bool {
	return slice.Index(v.UserIDs, uid) >= 0
}

//func NumberOfUpdateField(name string) string {
//	return "numberOf." + name
//}

//func (v *TeamDbo) UpdateNumberOf(name string, value int) dal.Update {
//	v.NumberOf[name] = value
//	return dal.Update{
//		Field: NumberOfUpdateField(name),
//		Value: value,
//	}
//}
