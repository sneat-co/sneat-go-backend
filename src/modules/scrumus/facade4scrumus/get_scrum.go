package facade4scrumus

import (
	"context"
	"github.com/sneat-co/sneat-go-backend/src/modules/scrumus/dbo4scrumus"
	"github.com/sneat-co/sneat-go-core/facade"
)

// GetScrum returns scrum data
func GetScrum(_ context.Context, user facade.User, _ facade.IDRequest) (scrum dbo4scrumus.Scrum, err error) {
	return
}
