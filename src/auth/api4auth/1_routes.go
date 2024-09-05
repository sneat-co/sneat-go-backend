package api4auth

import (
	"github.com/sneat-co/sneat-go-core/module"
	"net/http"
)

func RegisterHttpRoutes(handle module.HTTPHandleFunc) {
	handle(http.MethodPost, "/v0/auth/login-from-telegram-miniapp", httpLoginFromTelegramMiniapp)
	handle(http.MethodPost, "/v0/auth/login-from-telegram-widget", httpLoginFromTelegramWidget)
	handle(http.MethodDelete, "/v0/auth/disconnect", disconnect)
}
