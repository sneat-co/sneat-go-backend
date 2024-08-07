package dbo4listus

import (
	"fmt"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/strongo/validation"
	"strings"
)

// ListsCollection defines collection name
const ListsCollection = "lists"

// ListType list type
type ListType = string

const (
	// ListTypeGeneral = "general"
	ListTypeGeneral ListType = "general"

	// ListTypeToBuy = "to-by"
	ListTypeToBuy ListType = "to-buy"

	// ListTypeToDo = "to-do"
	ListTypeToDo ListType = "to-do"

	// ListTypeToWatch = "to-watch"
	ListTypeToWatch ListType = "to-watch"
)

// IsKnownListType checks if it is a known list type
func IsKnownListType(v string) bool {
	switch v {
	case ListTypeGeneral, ListTypeToBuy, ListTypeToWatch, ListTypeToDo:
		return true
	}
	return false
}

const ListIDSeparator = "!"

// GetFullListID returns full list ContactID
func GetFullListID(listType ListType, listID string) string {
	return listType + ListIDSeparator + listID
}

// ListBase DTO
type ListBase struct {
	Type  ListType `json:"type" firestore:"type"`
	Emoji string   `json:"emoji,omitempty" firestore:"emoji,omitempty"`
	// Title should be unique across owning team/company/group/etc
	Title string `json:"title" firestore:"title"`
}

// Validate returns error if not valid
func (v ListBase) Validate() error {
	if v.Type == "" {
		return validation.NewErrRecordIsMissingRequiredField("type")
	}
	if !IsKnownListType(v.Type) {
		return validation.NewErrBadRecordFieldValue("type", "unknown value: "+v.Type)
	}
	if strings.TrimSpace(v.Title) == "" {
		return validation.NewErrRecordIsMissingRequiredField("title")
	}
	return nil
}

// ListGroup DTO
type ListGroup struct {
	Type  string                `json:"type" firestore:"type"`
	Title string                `json:"title" firestore:"title"`
	Lists map[string]*ListBrief `json:"lists,omitempty"`
}

// Validate returns error if not valid
func (v ListGroup) Validate() error {
	if v.Type == "" {
		return validation.NewErrRecordIsMissingRequiredField("type")
	}
	if v.Title == "" {
		return validation.NewErrRecordIsMissingRequiredField("title")
	}
	if err := validateListBriefs(v.Lists); err != nil {
		return validation.NewErrBadRecordFieldValue("lists", err.Error())
	}
	return nil
}

// ListBrief DTO
type ListBrief struct {
	ListBase
	ItemsCount int `json:"itemsCount" firestore:"itemsCount"`
}

func validateListBriefs(lists map[string]*ListBrief) error {
	for id, list := range lists {
		if err := list.Validate(); err != nil {
			return fmt.Errorf("invalid list brief ID=%s: %w", id, err)
		}
	}
	return nil
}

// Validate returns error if not valid
func (v ListBrief) Validate() error {
	if err := v.ListBase.Validate(); err != nil {
		return err
	}
	if v.ItemsCount < 0 {
		return validation.NewErrBadRecordFieldValue("itemsCount", fmt.Sprintf("should be positive, got: %d", v.ItemsCount))
	}
	return nil
}

// ListDbo DTO
type ListDbo struct {
	ListBase
	dbmodels.WithModified
	dbmodels.WithUserIDs
	dbmodels.WithSpaceIDs

	Items []*ListItemBrief `json:"items,omitempty" firestore:"items,omitempty"`
	Count int              `json:"count" firestore:"count"`
}

// Validate returns error if not valid
func (v ListDbo) Validate() error {
	if err := v.WithSpaceIDs.Validate(); err != nil {
		return err
	}
	if err := v.WithUserIDs.Validate(); err != nil {
		return err
	}
	if err := v.ListBase.Validate(); err != nil {
		return err
	}
	if v.Count < 0 {
		return validation.NewErrBadRecordFieldValue("count", fmt.Sprintf("should be positive, got: %d", v.Count))
	}
	for i, item := range v.Items {
		if err := item.Validate(); err != nil {
			return validation.NewErrBadRecordFieldValue(fmt.Sprintf("items[%d]", i), err.Error())
		}
	}
	if count := len(v.Items); count != v.Count {
		return validation.NewErrBadRecordFieldValue("count", fmt.Sprintf("count != len(items): %d != %d", v.Count, count))
	}
	return nil
}
