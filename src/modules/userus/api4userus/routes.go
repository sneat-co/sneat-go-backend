package api4userus

import (
	"github.com/sneat-co/sneat-go-core/modules"
	"net/http"
)

// RegisterHttpRoutes initiates users module
func RegisterHttpRoutes(handle modules.HTTPHandleFunc) {
	handle(http.MethodPost, "/v0/users/init_user_record", httpInitUserRecord)
	handle(http.MethodPost, "/v0/users/set_user_country", httpSetUserCountry)
	handle(http.MethodPost, "/v0/users/login_with_telegram", httpLoginWithTelegram)
	//handle(http.MethodPost, "/v0/users/create_user", httpPostCreateUser)
}
