package apimapping

import (
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/julienschmidt/httprouter"
	"github.com/sneat-co/sneat-go-backend/src/coremodules/auth/api4auth"
	"github.com/sneat-co/sneat-go-backend/src/coremodules/common4all"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/gae_app/debtstracker/api4debtus/api4transfers"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/gae_app/debtstracker/api4debtus/unsorted"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/gae_app/debtstracker/dtdal"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/gae_app/debtstracker/reminders"
	"github.com/strongo/strongoapp"
	"net/http"
)

func InitApi(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, "/api4debtus/ping", botsfw.PingHandler)

	HandlerFunc := func(method, path string, handler strongoapp.HttpHandlerWithContext) {
		// TODO: Refactor optionsHandler so it's does not handle GET requests (see AuthOnly() for example)
		router.HandlerFunc(method, path, dtdal.HttpAppHost.HandleWithContext(handler))
		router.HandlerFunc(http.MethodOptions, path, dtdal.HttpAppHost.HandleWithContext(common4all.OptionsHandler))
	}

	GET := func(path string, handler strongoapp.HttpHandlerWithContext) {
		HandlerFunc(http.MethodGet, path, handler)
	}
	POST := func(path string, handler strongoapp.HttpHandlerWithContext) {
		HandlerFunc(http.MethodPost, path, handler)
	}

	POST("/api4debtus/auth/login-id", api4auth.OptionalAuth(api4auth.HandleAuthLoginId))
	POST("/api4debtus/auth/sign-in-with-pin", api4auth.OptionalAuth(api4auth.HandleSignInWithPin))
	POST("/api4debtus/auth/sign-in-with-fbm", api4auth.OptionalAuth(api4auth.HandleSignInWithFbm))
	POST("/api4debtus/auth/sign-in", api4auth.OptionalAuth(api4auth.HandleSignInWithCode))
	POST("/api4debtus/auth/fb/signed", api4auth.OptionalAuth(api4auth.HandleSignedWithFacebook))
	POST("/api4debtus/auth/google-plus/signed", api4auth.OptionalAuth(api4auth.HandleSignedInWithGooglePlus))
	POST("/api4debtus/auth/vk/signed", api4auth.OptionalAuth(api4auth.HandleSignedWithVK))
	//POST("/api4debtus/auth/email-sign-up", api4auth.HandleSignUpWithEmail)
	//POST("/api4debtus/auth/email-sign-in", api4auth.HandleSignInWithEmail)
	POST("/api4debtus/auth/request-password-reset", api4auth.HandleRequestPasswordReset)
	POST("/api4debtus/auth/change-password-and-sign-in", api4auth.HandleChangePasswordAndSignIn)
	POST("/api4debtus/auth/confirm-email-and-sign-in", api4auth.HandleConfirmEmailAndSignIn)
	POST("/api4debtus/auth/anonymous-sign-up", api4auth.HandleSignUpAnonymously)
	POST("/api4debtus/auth/anonymous-sign-in", api4auth.HandleSignInAnonymous)
	POST("/api4debtus/auth/disconnect", api4auth.AuthOnly(api4auth.HandleDisconnect))

	GET("/api4debtus/receipt-get", unsorted.HandleGetReceipt)
	POST("/api4debtus/receipt-create", api4auth.AuthOnly(unsorted.HandleCreateReceipt))
	POST("/api4debtus/receipt-send", api4auth.AuthOnlyWithUser(unsorted.HandleSendReceipt))
	POST("/api4debtus/receipt-set-channel", unsorted.HandleSetReceiptChannel)
	POST("/api4debtus/receipt-ack-accept", unsorted.HandleReceiptAccept)
	POST("/api4debtus/receipt-ack-decline", unsorted.HandleReceiptDecline)

	GET("/api4debtus/transfer", api4transfers.HandleGetTransfer)
	POST("/api4debtus/create-transfer", api4auth.AuthOnly(api4transfers.HandleCreateTransfer))

	POST("/api4debtus/bill-create", api4auth.AuthOnly(unsorted.HandleCreateBill))
	GET("/api4debtus/bill-get", api4auth.AuthOnly(unsorted.HandleGetBill))

	POST("/api4debtus/tg-helpers/currency-selected", api4auth.AuthOnly(unsorted.HandleTgHelperCurrencySelected))

	GET("/api4debtus/contact-get", api4auth.AuthOnly(unsorted.HandleGetContact))
	POST("/api4debtus/contact-create", api4auth.AuthOnly(unsorted.HandleCreateCounterparty))
	POST("/api4debtus/contact-update", api4auth.AuthOnly(unsorted.HandleUpdateCounterparty))
	POST("/api4debtus/contact-delete", api4auth.AuthOnly(unsorted.HandleDeleteContact))
	POST("/api4debtus/contact-archive", api4auth.AuthOnly(unsorted.HandleArchiveCounterparty))
	POST("/api4debtus/contact-activate", api4auth.AuthOnly(unsorted.HandleActivateCounterparty))

	POST("/api4debtus/group-create", api4auth.AuthOnlyWithUser(unsorted.HandlerCreateGroup))
	POST("/api4debtus/group-get", api4auth.AuthOnlyWithUser(unsorted.HandlerGetGroup))
	POST("/api4debtus/group-update", api4auth.AuthOnly(unsorted.HandlerUpdateGroup))
	POST("/api4debtus/group-delete", api4auth.AuthOnly(unsorted.HandlerDeleteGroup))
	POST("/api4debtus/group-set-contacts", api4auth.AuthOnlyWithUser(unsorted.HandlerSetContactsToGroup))
	POST("/api4debtus/join-groups", api4auth.AuthOnly(unsorted.HandleJoinGroups))

	GET("/api4debtus/user/api4transfers", api4auth.AuthOnlyWithUser(api4transfers.HandleUserTransfers))
	GET("/api4debtus/user/data/*rest", api4auth.AuthOnly(unsorted.HandleGetUserData))
	GET("/api4debtus/user/currencies", api4auth.AuthOnlyWithUser(unsorted.HandleGetUserCurrencies))
	GET("/api4debtus/user", unsorted.HandleUserInfo)

	GET("/api4debtus/me", api4auth.AuthOnlyWithUser(unsorted.HandleMe))
	POST("/api4debtus/user-set-name", api4auth.AuthOnly(unsorted.SetUserName))

	GET("/api4debtus/admin/latest/api4transfers", api4auth.AdminOnly(api4transfers.HandleAdminLatestTransfers))
	GET("/api4debtus/admin/latest/users", api4auth.AdminOnly(unsorted.HandleAdminLatestUsers))
	POST("/api4debtus/admin/find-user", api4auth.AdminOnly(unsorted.HandleAdminFindUser))
	GET("/api4debtus/admin/merge-user-contacts", api4auth.AdminOnly(unsorted.HandleAdminMergeUserContacts))

	POST("/api4debtus/analytics/visitor", unsorted.HandleSaveVisitorData)

	GET("/api4debtus/test/email", reminders.TestEmail)
	//POST("/api4debtus/invite-friend", inviteFriend)
	POST("/api4debtus/send-receipt", reminders.SendReceipt)
	POST("/api4debtus/invite/create", unsorted.CreateInvite)
}
