package facade4userus

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/dto4userus"
	"github.com/strongo/strongoapp/person"
	"testing"
)

func TestSetUserTitleRequest_Validate(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		request := dto4userus.InitUserRecordRequest{}
		if err := request.Validate(); err != nil {
			t.Fatal("expected no error for empty request")
		}
	})
	t.Run("valid", func(t *testing.T) {
		request := dto4userus.InitUserRecordRequest{
			AuthProvider: "password",
			Email:        "a@example.com",
			Names: &person.NameFields{
				FullName: "Test User",
			},
		}
		if err := request.Validate(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}
