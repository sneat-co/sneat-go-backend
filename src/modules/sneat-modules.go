package modules

import (
	"github.com/sneat-co/sneat-go-backend/src/auth"
	"github.com/sneat-co/sneat-go-backend/src/modules/assetus"
	"github.com/sneat-co/sneat-go-backend/src/modules/calendarium"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus"
	"github.com/sneat-co/sneat-go-backend/src/modules/generic"
	"github.com/sneat-co/sneat-go-backend/src/modules/invitus"
	"github.com/sneat-co/sneat-go-backend/src/modules/linkage"
	"github.com/sneat-co/sneat-go-backend/src/modules/listus"
	"github.com/sneat-co/sneat-go-backend/src/modules/retrospectus"
	"github.com/sneat-co/sneat-go-backend/src/modules/scrumus"
	"github.com/sneat-co/sneat-go-backend/src/modules/spaceus"
	"github.com/sneat-co/sneat-go-backend/src/modules/sportus"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus"
	"github.com/sneat-co/sneat-go-core/module"
)

func Modules() []module.Module {
	return []module.Module{
		auth.Module(),
		userus.Module(),
		spaceus.Module(),
		linkage.Module(),
		calendarium.Module(),
		contactus.Module(),
		invitus.Module(),
		assetus.Module(),
		listus.Module(),
		scrumus.Module(),
		retrospectus.Module(),
		sportus.Module(),
		generic.Module(),
		debtus.Module(),
	}
}
