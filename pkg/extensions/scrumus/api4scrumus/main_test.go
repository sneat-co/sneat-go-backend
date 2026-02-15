package api4scrumus

import (
	"context"
	"net/http"
	"os"
	"testing"

	"github.com/sneat-co/sneat-go-core/apicore"
	"github.com/sneat-co/sneat-go-core/apicore/verify"
	"github.com/sneat-co/sneat-go-core/facade"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

type mockUser struct {
	facade.UserContext
	userID string
}

func (m mockUser) GetUserID() string {
	return m.userID
}

type mockContextWithUser struct {
	facade.ContextWithUser
	ctx  context.Context
	user facade.UserContext
}

func (m mockContextWithUser) User() facade.UserContext {
	return m.user
}

func (m mockContextWithUser) Value(key any) any {
	return m.ctx.Value(key)
}

func setupMockVerify(t *testing.T) {
	oldVerify := apicore.VerifyRequestAndCreateUserContext
	apicore.VerifyRequestAndCreateUserContext = func(w http.ResponseWriter, r *http.Request, options verify.RequestOptions) (facade.ContextWithUser, error) {
		return mockContextWithUser{
			ctx:  t.Context(),
			user: mockUser{userID: "u1"},
		}, nil
	}
	t.Cleanup(func() {
		apicore.VerifyRequestAndCreateUserContext = oldVerify
	})
}
