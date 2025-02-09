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
	// ListTypeGeneral = "general" - not used at the moment
	ListTypeGeneral ListType = "general" // According to ChatGPT, it's more correct then "Generic"

	// ListTypeToBuy = "to-by"
	ListTypeToBuy ListType = "buy"

	// ListTypeToDo = "to-do"
	ListTypeToDo ListType = "do"

	// ListTypeToWatch = "to-watch"
	ListTypeToWatch ListType = "watch"

	ListTypeToRead ListType = "read"
)

// IsKnownListType checks if it is a known list type
func IsKnownListType(v string) bool {
	switch v {
	case ListTypeGeneral, ListTypeToBuy, ListTypeToDo, ListTypeToWatch, ListTypeToRead:
		return true
	}
	return false
}

const ListIDSeparator = "!"

// NewListKey returns ist key in format of "{listType}!{listSubID}"
func NewListKey(listType ListType, listID string) ListKey {
	return ListKey(listType + ListIDSeparator + listID)
}

type ListKey string

const (
	BuyGroceriesListID = "buy!groceries"
	DoTasksListID      = "do!tasks"
	ReadBooksListID    = "read!books"
	WatchMoviesListID  = "watch!movies"
)

func IsStandardList(listKey ListKey) bool {
	return listKey == BuyGroceriesListID || listKey == DoTasksListID || listKey == ReadBooksListID || listKey == WatchMoviesListID
}

func (v ListKey) Validate() error {
	if s := strings.TrimSpace(string(v)); s == "" {
		return validation.NewValidationError("list key is empty string")
	} else if s != string(v) {
		return validation.NewValidationError("list key has leading or trailing spaces")
	} else if i := strings.Index(s, ListIDSeparator); i < 0 {
		return validation.NewValidationError(fmt.Sprintf("list key does not contain required separator '%s'", ListIDSeparator))
	} else if separatorsCount := strings.Count(s, ListIDSeparator); separatorsCount > 1 {
		return validation.NewValidationError(fmt.Sprintf("list key expected to have only 1 '%s' spearator, contains %d separators: %s", ListIDSeparator, separatorsCount, v))
	} else if strings.TrimSpace(s[:i]) == "" {
		return validation.NewValidationError("list type is empty")
	} else if strings.TrimSpace(s[i+1:]) == "" {
		return validation.NewValidationError("list sub ID is empty")
	}
	if listType := v.ListType(); !IsKnownListType(listType) {
		return validation.NewValidationError("unknown list type: " + listType)
	}
	return nil
}

func (v ListKey) ListType() ListType {
	if i := strings.Index(string(v), ListIDSeparator); i > 0 {
		return ListType(v[:i])
	}
	return ""
}

func (v ListKey) ListSubID() string {
	if i := strings.Index(string(v), ListIDSeparator); i > 0 {
		return string(v[i+1:])
	}
	return ""
}

// ListBase DTO
type ListBase struct {
	Type ListType `json:"type" firestore:"type"`

	// Title should be unique across owning team/company/group/etc
	Title string `json:"title" firestore:"title"`

	// Emoji is optional, by default uses emoji of the list type
	Emoji string `json:"emoji,omitempty" firestore:"emoji,omitempty"`
}

// Validate returns error if not valid
func (v ListBase) Validate() error {
	if v.Type == "" {
		return validation.NewErrRecordIsMissingRequiredField("type")
	}
	if !IsKnownListType(v.Type) {
		return validation.NewErrBadRecordFieldValue("type", "unknown value: "+v.Type)
	}
	//if strings.TrimSpace(v.Title) == "" {
	//	return validation.NewErrRecordIsMissingRequiredField("title")
	//}
	return nil
}

// ListGroup DTO
type ListGroup struct {
	Type  string     `json:"type" firestore:"type"`
	Title string     `json:"title" firestore:"title"`
	Lists ListBriefs `json:"lists,omitempty"`
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

type ListBriefs map[string]*ListBrief

func validateListBriefs(lists ListBriefs) error {
	for id, list := range lists {
		if err := list.Validate(); err != nil {
			return fmt.Errorf("invalid list brief ContactID=%s: %w", id, err)
		}
	}
	return nil
}

// Validate returns error if not valid
func (v *ListBrief) Validate() error {
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

	Items       []*ListItemBrief `json:"items,omitempty" firestore:"items,omitempty"`
	RecentItems []*ListItemBrief `json:"recentItems,omitempty" firestore:"recentItems,omitempty"`

	Count int `json:"count" firestore:"count"`
}

func (v *ListDbo) AddListItem(item *ListItemBrief) (addedItem *ListItemBrief) {
	for _, existingItem := range v.Items {
		if existingItem.Title == item.Title && existingItem.Emoji == item.Emoji {
			addedItem = existingItem
			if existingItem.IsDone() {
				existingItem.Status = item.Status
			}
			return
		}
	}
	v.Items = append(v.Items, item)
	v.Count = len(v.Items)
	addedItem = item
	return
}

// Validate returns error if not valid
func (v *ListDbo) Validate() error {
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
