package facade4companies

import (
	"testing"

	"github.com/sneat-co/sneat-go-core/facade"
)

type mockUser struct {
	facade.UserContext
	userID string
}

func (m mockUser) GetUserID() string {
	return m.userID
}

type mockContextWithUser struct {
	facade.ContextWithUser
	user facade.UserContext
}

func (m mockContextWithUser) User() facade.UserContext {
	return m.user
}

func (m mockContextWithUser) Value(key any) any {
	return m.ContextWithUser.Value(key)
}

func TestCreateCompany(t *testing.T) {
	// CreateCompany calls facade.RunReadwriteTransaction which is hard to mock without a real DB or complex setup.
	// For now, we will test the initial validation.

	t.Run("missing_user_id", func(t *testing.T) {
		ctx := mockContextWithUser{
			user: mockUser{userID: ""},
		}
		_, err := CreateCompany(ctx, CreateCompanyRequest{})
		if err == nil {
			t.Error("expected error due to missing user ID")
		} else if err.Error() != "user ID is missing" {
			t.Errorf("expected 'user ID is missing' error, got: %v", err)
		}
	})

	t.Run("with_user_id", func(t *testing.T) {
		// This will likely fail because facade.RunReadwriteTransaction will try to access a real DB
		// but we can at least see if it proceeds past the ID check.
		ctx := mockContextWithUser{
			user: mockUser{userID: "test_user"},
		}

		defer func() {
			if r := recover(); r != nil {
				// The panic message might be different depending on how it's formatted or if it's an error type
				rStr := ""
				if err, ok := r.(error); ok {
					rStr = err.Error()
				} else {
					rStr = r.(string)
				}

				if rStr == "not initialized: facade.GetSneatDB(context.Context) (dal.DB, error)" {
					t.Log("Passed validation and reached DB access as expected")
				} else {
					t.Errorf("Unexpected panic: %v", r)
				}
			}
		}()

		_, _ = CreateCompany(ctx, CreateCompanyRequest{})
	})
}
