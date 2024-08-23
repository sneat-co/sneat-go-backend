package api4auth

import "github.com/sneat-co/sneat-go-core/module"

func RegisterHttpRoutes(handle module.HTTPHandleFunc) {
	handle("POST", "/v0/auth/signing-from-telegram-miniapp", httpSignInFromTelegramMiniapp)
}
