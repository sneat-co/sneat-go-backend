package facade4listus

import (
	"fmt"
	"github.com/sneat-co/sneat-go-backend/src/modules/listus/dbo4listus"
	"github.com/sneat-co/sneat-go-backend/src/modules/teamus/dto4teamus"
	"github.com/sneat-co/sneat-go-core/validate"
	"github.com/strongo/validation"
	"strings"
)

// ListRequest DTO
type ListRequest struct {
	dto4teamus.TeamRequest
	ListID   string `json:"listID"`
	listType dbo4listus.ListType
}

func (v *ListRequest) ListType() dbo4listus.ListType {
	if v.listType != "" {
		return v.listType
	}
	if v.ListID == "" {
		return ""
	}
	i := strings.Index(v.ListID, dbo4listus.ListIDSeparator)
	if i < 0 {
		return ""
	}
	return v.ListID[:i]
}

// Validate returns error if not valid
func (v *ListRequest) Validate() error {
	if err := v.TeamRequest.Validate(); err != nil {
		return err
	}
	if strings.TrimSpace(v.ListID) == "" {
		return validation.NewErrRequestIsMissingRequiredField("list")
	}
	switch v.ListType() {
	case "", "to-buy", "to-do":
	default:
		return fmt.Errorf("")
	}
	return nil
}

// CreateListRequest DTO
type CreateListRequest struct {
	dto4teamus.TeamRequest
	Type  string `json:"type"`
	Title string `json:"title"`
}

// Validate returns error if not valid
func (v CreateListRequest) Validate() error {
	if err := v.TeamRequest.Validate(); err != nil {
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
