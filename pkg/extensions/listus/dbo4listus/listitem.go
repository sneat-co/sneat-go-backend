package dbo4listus

import (
	"strings"

	"github.com/sneat-co/sneat-go-backend/pkg/extensions/listus/const4listus"
	"github.com/strongo/strongoapp/with"
	"github.com/strongo/validation"
)

// ListItemBase DTO
type ListItemBase struct {
	Title string `json:"title" firestore:"title"`
	Emoji string `json:"emoji,omitempty" firestore:"emoji,omitempty"`

	Status const4listus.ListItemStatus `json:"status,omitempty" firestore:"status,omitempty"`
}

func (v ListItemBase) IsDone() bool {
	return v.Status == const4listus.ListItemStatusDone
}

// Validate returns error if not valid
func (v ListItemBase) Validate() error {
	if strings.TrimSpace(v.Title) == "" {
		return validation.NewErrRecordIsMissingRequiredField("title")
	}
	return nil
}

// ListItemBrief DTO
type ListItemBrief struct {
	ID string `json:"id" firestore:"id"`
	ListItemBase
	with.CreatedFields
}

// Validate returns error if not valid
func (v ListItemBrief) Validate() error {
	if strings.TrimSpace(v.ID) == "" {
		return validation.NewErrRecordIsMissingRequiredField("id")
	}
	if err := v.ListItemBase.Validate(); err != nil {
		return err
	}
	if err := v.CreatedFields.Validate(); err != nil {
		return err
	}
	return nil
}

// ListItemDbo DTO
type ListItemDbo struct {
	ListItemBase
}
