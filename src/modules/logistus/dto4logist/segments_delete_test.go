package dto4logist

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/logistus/dbo4logist"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDeleteSegmentsRequest_Validate(t *testing.T) {
	type fields struct {
		OrderRequest   OrderRequest
		SegmentsFilter dbo4logist.SegmentsFilter
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "empty",
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.NotNil(t, err)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := DeleteSegmentsRequest{
				OrderRequest:   tt.fields.OrderRequest,
				SegmentsFilter: tt.fields.SegmentsFilter,
			}
			tt.wantErr(t, v.Validate(), "Validate()")
		})
	}
}
