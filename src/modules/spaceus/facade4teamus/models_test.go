package facade4teamus

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/spaceus/dto4spaceus"
	"testing"
)

func TestCreateSpaceRequest_Validate(t *testing.T) {
	request := dto4spaceus.CreateSpaceRequest{Title: ""}
	if request.Validate() == nil {
		t.Error("request.Validate() == nil")
	}
}
