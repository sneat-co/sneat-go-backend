package api

import (
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/julienschmidt/httprouter"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/dtdal"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/reminders"
	"github.com/strongo/strongoapp"
)

func InitApi(router *httprouter.Router) {
	router.HandlerFunc("GET", "/api/ping", botsfw.PingHandler)

	HandlerFunc := func(method, path string, handler strongoapp.HttpHandlerWithContext) {
		// TODO: Refactor optionsHandler so it's does not handle GET requests (see AuthOnly() for example)
		router.HandlerFunc(method, path, dtdal.HttpAppHost.HandleWithContext(handler))
		router.HandlerFunc("OPTIONS", path, dtdal.HttpAppHost.HandleWithContext(optionsHandler))
	}

	GET := func(path string, handler strongoapp.HttpHandlerWithContext) {
		HandlerFunc("GET", path, handler)
	}
	POST := func(path string, handler strongoapp.HttpHandlerWithContext) {
		HandlerFunc("POST", path, handler)
	}

	POST("/api/auth/login-id", OptionalAuth(handleAuthLoginId))
	POST("/api/auth/sign-in-with-pin", OptionalAuth(handleSignInWithPin))
	POST("/api/auth/sign-in-with-fbm", OptionalAuth(handleSignInWithFbm))
	POST("/api/auth/sign-in", OptionalAuth(handleSignInWithCode))
	POST("/api/auth/fb/signed", OptionalAuth(handleSignedWithFacebook))
	POST("/api/auth/google-plus/signed", OptionalAuth(handleSignedInWithGooglePlus))
	POST("/api/auth/vk/signed", OptionalAuth(handleSignedWithVK))
	POST("/api/auth/email-sign-up", handleSignUpWithEmail)
	POST("/api/auth/email-sign-in", handleSignInWithEmail)
	POST("/api/auth/request-password-reset", handleRequestPasswordReset)
	POST("/api/auth/change-password-and-sign-in", handleChangePasswordAndSignIn)
	POST("/api/auth/confirm-email-and-sign-in", handleConfirmEmailAndSignIn)
	POST("/api/auth/anonymous-sign-up", handleSignUpAnonymously)
	POST("/api/auth/anonymous-sign-in", handleSignInAnonymous)
	POST("/api/auth/disconnect", AuthOnly(handleDisconnect))

	GET("/api/receipt-get", handleGetReceipt)
	POST("/api/receipt-create", AuthOnly(handleCreateReceipt))
	POST("/api/receipt-send", AuthOnlyWithUser(handleSendReceipt))
	POST("/api/receipt-set-channel", handleSetReceiptChannel)
	POST("/api/receipt-ack-accept", handleReceiptAccept)
	POST("/api/receipt-ack-decline", handleReceiptDecline)

	GET("/api/transfer", handleGetTransfer)
	POST("/api/create-transfer", AuthOnly(handleCreateTransfer))

	POST("/api/bill-create", AuthOnly(handleCreateBill))
	GET("/api/bill-get", AuthOnly(handleGetBill))

	POST("/api/tg-helpers/currency-selected", AuthOnly(handleTgHelperCurrencySelected))

	GET("/api/contact-get", AuthOnly(handleGetContact))
	POST("/api/contact-create", AuthOnly(handleCreateCounterparty))
	POST("/api/contact-update", AuthOnly(handleUpdateCounterparty))
	POST("/api/contact-delete", AuthOnly(handleDeleteContact))
	POST("/api/contact-archive", AuthOnly(handleArchiveCounterparty))
	POST("/api/contact-activate", AuthOnly(handleActivateCounterparty))

	POST("/api/group-create", AuthOnlyWithUser(handlerCreateGroup))
	POST("/api/group-get", AuthOnlyWithUser(handlerGetGroup))
	POST("/api/group-update", AuthOnly(handlerUpdateGroup))
	POST("/api/group-delete", AuthOnly(handlerDeleteGroup))
	POST("/api/group-set-contacts", AuthOnlyWithUser(handlerSetContactsToGroup))
	POST("/api/join-groups", AuthOnly(handleJoinGroups))

	GET("/api/user/transfers", AuthOnlyWithUser(handleUserTransfers))
	GET("/api/user/data/*rest", AuthOnly(handleGetUserData))
	GET("/api/user/currencies", AuthOnlyWithUser(handleGetUserCurrencies))
	GET("/api/user", handleUserInfo)

	GET("/api/me", AuthOnlyWithUser(handleMe))
	POST("/api/user-set-name", AuthOnly(setUserName))

	GET("/api/admin/latest/transfers", adminOnly(handleAdminLatestTransfers))
	GET("/api/admin/latest/users", adminOnly(handleAdminLatestUsers))
	POST("/api/admin/find-user", adminOnly(handleAdminFindUser))
	GET("/api/admin/merge-user-contacts", adminOnly(handleAdminMergeUserContacts))

	POST("/api/analytics/visitor", handleSaveVisitorData)

	GET("/api/test/email", reminders.TestEmail)
	//POST("/api/invite-friend", inviteFriend)
	POST("/api/send-receipt", reminders.SendReceipt)
	POST("/api/invite/create", CreateInvite)
}
