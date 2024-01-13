package facade4listus

import (
	"fmt"
	"github.com/sneat-co/sneat-go-backend/src/modules/listus/models4listus"
	"github.com/sneat-co/sneat-go-core/validate"
	"github.com/strongo/validation"
	"strings"
)

// CreateListItemsRequest DTO
type CreateListItemsRequest struct {
	ListRequest
	Items []CreateListItemRequest `json:"items"`
}

// Validate returns error if not valid
func (v CreateListItemsRequest) Validate() error {
	if err := v.TeamRequest.Validate(); err != nil {
		return err
	}
	for i, item := range v.Items {
		if err := item.Validate(); err != nil {
			return validation.NewErrBadRequestFieldValue(fmt.Sprintf("items[%v]", i), err.Error())
		}
	}
	return nil
}

type CreateListItemRequest struct {
	ID string `json:"id"`
	models4listus.ListItemBase
}

func (v CreateListItemRequest) Validate() error {
	if strings.TrimSpace(v.ID) != "" {
		if err := validate.RecordID(v.ID); err != nil {
			return validation.NewErrBadRequestFieldValue("id", err.Error())
		}
	}
	if err := v.ListItemBase.Validate(); err != nil {
		return validation.NewBadRequestError(err)
	}
	return nil
}

// CreateListItemResponse DTO
type CreateListItemResponse struct {
	ID string `json:"id"`
}

// ListItemIDsRequest DTO
type ListItemIDsRequest struct {
	ListRequest
	ItemIDs []string `json:"itemIDs"`
}

// Validate returns error if not valid
func (v ListItemIDsRequest) Validate() error {
	if err := v.ListRequest.Validate(); err != nil {
		return err
	}
	// Intentionally do not check for empty ItemIDs
	//if len(v.ItemIDs) == 0 {
	//	return validation.NewErrRequestIsMissingRequiredField("itemIDs")
	//}
	for i, id := range v.ItemIDs {
		if strings.TrimSpace(id) == "" {
			return validation.NewErrRecordIsMissingRequiredField(fmt.Sprintf("itemIDs[%v]", i))
		}
	}
	return nil
}

// ReorderListItemsRequest DTO
type ReorderListItemsRequest struct {
	ListItemIDsRequest
	ToIndex int `json:"toIndex"`
}

// Validate returns error if not valid
func (v ReorderListItemsRequest) Validate() error {
	if err := v.ListItemIDsRequest.Validate(); err != nil {
		return err
	}
	if v.ToIndex < 0 {
		return validation.NewErrBadRecordFieldValue("toIndex", "must be >= 0")
	}
	return nil
}

// ListItemsSetIsDoneRequest DTO
type ListItemsSetIsDoneRequest struct {
	ListItemIDsRequest
	IsDone bool `json:"isDone"`
}

// Validate returns error if not valid
func (v ListItemsSetIsDoneRequest) Validate() error {
	if err := v.ListItemIDsRequest.Validate(); err != nil {
		return err
	}
	return nil
}
