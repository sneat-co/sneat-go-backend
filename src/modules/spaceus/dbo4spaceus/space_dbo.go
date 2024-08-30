package dbo4spaceus

import (
	"errors"
	"fmt"
	"github.com/sneat-co/sneat-go-backend/src/modules/spaceus/core4spaceus"
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

// SpaceMetric record
type SpaceMetric struct {
	ID      string      `json:"id" firestore:"id"`
	Title   string      `json:"title" firestore:"title"`
	Mode    string      `json:"mode" firestore:"mode"` // Possible values: personal, team
	Type    string      `json:"type" firestore:"type"` // Possible values: bool, int, str
	Min     *int        `json:"min,omitempty" firestore:"min,omitempty"`
	Max     *int        `json:"max,omitempty" firestore:"max,omitempty"`
	Bool    *BoolMetric `json:"bool,omitempty" firestore:"bool,omitempty"`
	Options []string    `json:"options" firestore:"options"`
}

// Validate validates SpaceMetric record
func (v *SpaceMetric) Validate() error {
	if v.Title == "" {
		return validation.NewErrRecordIsMissingRequiredField("title")
	}
	switch v.Mode {
	case "":
		return validation.NewErrRecordIsMissingRequiredField("mode")
	case "personal", "space":
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

// SpaceMeetingInfo record
type SpaceMeetingInfo struct {
	ID       string     `json:"id,omitempty" firestore:"id,omitempty"`
	Stage    string     `json:"stage" firestore:"stage"`
	Started  *time.Time `json:"started,omitempty" firestore:"started,omitempty"`
	Finished *time.Time `json:"finished,omitempty" firestore:"finished,omitempty"`
}

// EarliestPossibleTime 2020-01-01
var EarliestPossibleTime = time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)

// Validate validates record
func (v *SpaceMeetingInfo) Validate() error {
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

// SpaceMeetings record
type SpaceMeetings struct {
	Retrospective *SpaceMeetingInfo `json:"retrospective,omitempty" firestore:"retrospective,omitempty"`
	Scrum         *SpaceMeetingInfo `json:"scrum,omitempty" firestore:"scrum,omitempty"`
}

// SpaceBrief is a base class for SpaceDbo
type SpaceBrief struct {
	Type   core4spaceus.SpaceType `json:"type" firestore:"type"`
	Title  string                 `json:"title" firestore:"title"`
	Status dbmodels.Status        `json:"status" firestore:"status"`

	Modules []string `json:"modules,omitempty" firestore:"modules,omitempty"`

	with.OptionalCountryID

	// TODO: This should be populated
	ParentSpaceID string `json:"parentSpaceID,omitempty" firestore:"parentSpaceID,omitempty"`
}

func (v SpaceBrief) Validate() error {
	v.Title = strings.TrimSpace(v.Title)
	if v.Type == "" {
		return validation.NewErrRecordIsMissingRequiredField("type")
	}
	if !core4spaceus.IsValidSpaceType(v.Type) {
		if v.Title == "" {
			return validation.NewErrBadRecordFieldValue("type", "unknown value")
		}
		return validation.NewErrBadRecordFieldValue("type", "unknown value for team:"+v.Title)
	}
	if v.Title == "" && v.Type != core4spaceus.SpaceTypeFamily && v.Type != core4spaceus.SpaceTypePrivate {
		return errors.New("non family team is required to have a title")
	}
	if v.Status == "" {
		return validation.NewErrRequestIsMissingRequiredField("status")
	} else if !IsKnownSpaceStatus(v.Status) {
		return validation.NewErrBadRecordFieldValue("status", "unknown value: "+v.Status)
	}
	if err := v.OptionalCountryID.Validate(); err != nil {
		return err
	}
	return nil
}

func IsKnownSpaceStatus(status dbmodels.Status) bool {
	switch status {
	case dbmodels.StatusActive, dbmodels.StatusArchived, dbmodels.StatusDeleted, dbmodels.StatusDraft:
		return true
	}
	return false
}

const NumberOfMembersFieldName = "members"

// SpaceDbo record
type SpaceDbo struct {
	SpaceBrief
	with.CreatedFields
	dbmodels.WithUpdatedAndVersion
	dbmodels.WithUserIDs

	//NumberOf map[string]int `json:"numberOf,omitempty" firestore:"numberOf,omitempty"`
	//
	//Contacts   []*briefs4contactus.ContactBrief `json:"contacts,omitempty" firestore:"contacts,omitempty"`

	//
	Timezone *dbmodels.Timezone `json:"timezone,omitempty" firestore:"timezone,omitempty"`
	//
	dbmodels.WithPreferredLocale
	dbmodels.WithPrimaryCurrency

	//
	Metrics []*SpaceMetric `json:"metrics,omitempty" firestore:"metrics,omitempty"`
}

//func (v *SpaceDbo) SetNumberOf(name string, value int) (update dal.Update) {
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
func (v *SpaceDbo) IncreaseVersion(timestamp time.Time, updatedBy string) int {
	v.Version++
	v.UpdatedAt = timestamp
	v.UpdatedBy = updatedBy
	return v.Version
}

// Validate validates record
func (v *SpaceDbo) Validate() error {
	if err := v.SpaceBrief.Validate(); err != nil {
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

// HasUser checks if team has a user with given ContactID
func (v *SpaceDbo) HasUser(uid string) bool {
	return slice.Index(v.UserIDs, uid) >= 0
}

//func NumberOfUpdateField(name string) string {
//	return "numberOf." + name
//}

//func (v *SpaceDbo) UpdateNumberOf(name string, value int) dal.Update {
//	v.NumberOf[name] = value
//	return dal.Update{
//		Field: NumberOfUpdateField(name),
//		Value: value,
//	}
//}
