package facade4teamus

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/teamus/dto4teamus"
	"testing"
)

func TestCreateSpaceRequest_Validate(t *testing.T) {
	request := dto4teamus.CreateSpaceRequest{Title: ""}
	if request.Validate() == nil {
		t.Error("request.Validate() == nil")
	}
}
