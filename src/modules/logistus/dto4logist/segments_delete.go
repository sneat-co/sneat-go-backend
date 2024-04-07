package dto4logist

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/logistus/models4logist"
)

// DeleteSegmentsRequest represents a request to delete segments
type DeleteSegmentsRequest struct {
	OrderRequest
	models4logist.SegmentsFilter
}

// Validate returns an error if the request is invalid.
func (v DeleteSegmentsRequest) Validate() error {
	if err := v.OrderRequest.Validate(); err != nil {
		return err
	}
	if err := v.SegmentsFilter.Validate(); err != nil {
		return err
	}
	return nil
}
