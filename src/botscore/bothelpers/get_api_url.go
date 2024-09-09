package bothelpers

import core "github.com/sneat-co/sneat-go-core"

func GetBotWebAppUrl() string {
	if core.IsInProd() {
		return "https://sneat.app/pwa/"
	} else {
		return "https://local-app.sneat.ws/"
	}
}
