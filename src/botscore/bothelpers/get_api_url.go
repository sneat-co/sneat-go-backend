package bothelpers

import core "github.com/sneat-co/sneat-go-core"

func GetBotWebAppUrl() string {
	if core.IsInProd() {
		return "https://sneat-eur3-1.appspot.com/pwa/"
	} else {
		return "https://local-app.sneat.ws/"
	}
}
