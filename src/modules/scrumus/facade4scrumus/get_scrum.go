package facade4scrumus

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/scrumus/dbo4scrumus"
	"github.com/sneat-co/sneat-go-core/facade"
)

// GetScrum returns scrum data
func GetScrum(_ facade.ContextWithUser, _ facade.IDRequest) (scrum dbo4scrumus.Scrum, err error) {
	panic("implement me")
}
