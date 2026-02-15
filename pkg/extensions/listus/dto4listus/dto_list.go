package dto4listus

import (
	"strings"

	"github.com/sneat-co/sneat-core-modules/spaceus/dto4spaceus"
	"github.com/sneat-co/sneat-go-backend/pkg/extensions/listus/dbo4listus"
	"github.com/sneat-co/sneat-go-core/validate"
	"github.com/strongo/validation"
)

// ListRequest DTO
type ListRequest struct {
	dto4spaceus.SpaceRequest
	ListID dbo4listus.ListKey `json:"listID"`
}

// Validate returns error if not valid
func (v *ListRequest) Validate() error {
	if err := v.SpaceRequest.Validate(); err != nil {
		return err
	}
	if err := v.ListID.Validate(); err != nil {
		return validation.NewErrBadRecordFieldValue("listID", err.Error())
	}
	return nil
}

// CreateListRequest DTO
type CreateListRequest struct {
	dto4spaceus.SpaceRequest
	Type  dbo4listus.ListType `json:"type"`
	Title string              `json:"title"`
}

// Validate returns error if not valid
func (v CreateListRequest) Validate() error {
	if err := v.SpaceRequest.Validate(); err != nil {
		return err
	}
	if strings.TrimSpace(v.Type) == "" {
		return validation.NewErrRequestIsMissingRequiredField("type")
	}
	if err := validate.RequestTitle(v.Title, "title"); err != nil {
		return err
	}
	return nil
}

// ListItemRequest DTO
type ListItemRequest struct {
	ListRequest
	ItemID string `json:"item"`
}

// Validate returns error if not valid
func (v ListItemRequest) Validate() error {
	if err := v.ListRequest.Validate(); err != nil {
		return err
	}
	if strings.TrimSpace(v.ItemID) == "" {
		return validation.NewErrRequestIsMissingRequiredField("item")
	}
	return nil
}

// CreateListResponse DTO
type CreateListResponse struct {
	ID string `json:"id"`
}
