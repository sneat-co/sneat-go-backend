package apimapping

import (
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/julienschmidt/httprouter"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/api"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/api/transfers"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/api/unsorted"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/dtdal"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/reminders"
	"github.com/strongo/strongoapp"
)

func InitApi(router *httprouter.Router) {
	router.HandlerFunc("GET", "/api/ping", botsfw.PingHandler)

	HandlerFunc := func(method, path string, handler strongoapp.HttpHandlerWithContext) {
		// TODO: Refactor optionsHandler so it's does not handle GET requests (see AuthOnly() for example)
		router.HandlerFunc(method, path, dtdal.HttpAppHost.HandleWithContext(handler))
		router.HandlerFunc("OPTIONS", path, dtdal.HttpAppHost.HandleWithContext(api.OptionsHandler))
	}

	GET := func(path string, handler strongoapp.HttpHandlerWithContext) {
		HandlerFunc("GET", path, handler)
	}
	POST := func(path string, handler strongoapp.HttpHandlerWithContext) {
		HandlerFunc("POST", path, handler)
	}

	POST("/api/auth/login-id", unsorted.OptionalAuth(unsorted.HandleAuthLoginId))
	POST("/api/auth/sign-in-with-pin", unsorted.OptionalAuth(unsorted.HandleSignInWithPin))
	POST("/api/auth/sign-in-with-fbm", unsorted.OptionalAuth(unsorted.HandleSignInWithFbm))
	POST("/api/auth/sign-in", unsorted.OptionalAuth(unsorted.HandleSignInWithCode))
	POST("/api/auth/fb/signed", unsorted.OptionalAuth(unsorted.HandleSignedWithFacebook))
	POST("/api/auth/google-plus/signed", unsorted.OptionalAuth(unsorted.HandleSignedInWithGooglePlus))
	POST("/api/auth/vk/signed", unsorted.OptionalAuth(unsorted.HandleSignedWithVK))
	POST("/api/auth/email-sign-up", unsorted.HandleSignUpWithEmail)
	POST("/api/auth/email-sign-in", unsorted.HandleSignInWithEmail)
	POST("/api/auth/request-password-reset", unsorted.HandleRequestPasswordReset)
	POST("/api/auth/change-password-and-sign-in", unsorted.HandleChangePasswordAndSignIn)
	POST("/api/auth/confirm-email-and-sign-in", unsorted.HandleConfirmEmailAndSignIn)
	POST("/api/auth/anonymous-sign-up", unsorted.HandleSignUpAnonymously)
	POST("/api/auth/anonymous-sign-in", unsorted.HandleSignInAnonymous)
	POST("/api/auth/disconnect", unsorted.AuthOnly(unsorted.HandleDisconnect))

	GET("/api/receipt-get", unsorted.HandleGetReceipt)
	POST("/api/receipt-create", unsorted.AuthOnly(unsorted.HandleCreateReceipt))
	POST("/api/receipt-send", unsorted.AuthOnlyWithUser(unsorted.HandleSendReceipt))
	POST("/api/receipt-set-channel", unsorted.HandleSetReceiptChannel)
	POST("/api/receipt-ack-accept", unsorted.HandleReceiptAccept)
	POST("/api/receipt-ack-decline", unsorted.HandleReceiptDecline)

	GET("/api/transfer", transfers.HandleGetTransfer)
	POST("/api/create-transfer", unsorted.AuthOnly(transfers.HandleCreateTransfer))

	POST("/api/bill-create", unsorted.AuthOnly(unsorted.HandleCreateBill))
	GET("/api/bill-get", unsorted.AuthOnly(unsorted.HandleGetBill))

	POST("/api/tg-helpers/currency-selected", unsorted.AuthOnly(unsorted.HandleTgHelperCurrencySelected))

	GET("/api/contact-get", unsorted.AuthOnly(unsorted.HandleGetContact))
	POST("/api/contact-create", unsorted.AuthOnly(unsorted.HandleCreateCounterparty))
	POST("/api/contact-update", unsorted.AuthOnly(unsorted.HandleUpdateCounterparty))
	POST("/api/contact-delete", unsorted.AuthOnly(unsorted.HandleDeleteContact))
	POST("/api/contact-archive", unsorted.AuthOnly(unsorted.HandleArchiveCounterparty))
	POST("/api/contact-activate", unsorted.AuthOnly(unsorted.HandleActivateCounterparty))

	POST("/api/group-create", unsorted.AuthOnlyWithUser(unsorted.HandlerCreateGroup))
	POST("/api/group-get", unsorted.AuthOnlyWithUser(unsorted.HandlerGetGroup))
	POST("/api/group-update", unsorted.AuthOnly(unsorted.HandlerUpdateGroup))
	POST("/api/group-delete", unsorted.AuthOnly(unsorted.HandlerDeleteGroup))
	POST("/api/group-set-contacts", unsorted.AuthOnlyWithUser(unsorted.HandlerSetContactsToGroup))
	POST("/api/join-groups", unsorted.AuthOnly(unsorted.HandleJoinGroups))

	GET("/api/user/transfers", unsorted.AuthOnlyWithUser(transfers.HandleUserTransfers))
	GET("/api/user/data/*rest", unsorted.AuthOnly(unsorted.HandleGetUserData))
	GET("/api/user/currencies", unsorted.AuthOnlyWithUser(unsorted.HandleGetUserCurrencies))
	GET("/api/user", unsorted.HandleUserInfo)

	GET("/api/me", unsorted.AuthOnlyWithUser(unsorted.HandleMe))
	POST("/api/user-set-name", unsorted.AuthOnly(unsorted.SetUserName))

	GET("/api/admin/latest/transfers", unsorted.AdminOnly(transfers.HandleAdminLatestTransfers))
	GET("/api/admin/latest/users", unsorted.AdminOnly(unsorted.HandleAdminLatestUsers))
	POST("/api/admin/find-user", unsorted.AdminOnly(unsorted.HandleAdminFindUser))
	GET("/api/admin/merge-user-contacts", unsorted.AdminOnly(unsorted.HandleAdminMergeUserContacts))

	POST("/api/analytics/visitor", unsorted.HandleSaveVisitorData)

	GET("/api/test/email", reminders.TestEmail)
	//POST("/api/invite-friend", inviteFriend)
	POST("/api/send-receipt", reminders.SendReceipt)
	POST("/api/invite/create", unsorted.CreateInvite)
}
