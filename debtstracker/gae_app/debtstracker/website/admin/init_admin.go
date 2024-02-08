package admin

import (
	"github.com/julienschmidt/httprouter"
)

type router interface {
	GET(path string, handle httprouter.Handle)
}

func InitAdmin(router router) {
	router.GET("/admin/latest", LatestPage)
	router.GET("/admin/clean", CleanupPage)
	//strongoapp.AddHttpHandler("/admin/mass-invites", LatestPage)
	router.GET("/admin/fix/transfers", FixTransfersHandler)
	//router.GET("/admin/fbm/set", dtb_fbm.SetupFbm)
}
