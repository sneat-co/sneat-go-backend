package inspector

import "github.com/julienschmidt/httprouter"

type router interface {
	GET(path string, handle httprouter.Handle)
}

func InitInspector(router router) {
	router.GET("/inspector/user", userPage)
	router.GET("/inspector/contact", contactPage{}.contactPageHandler)
	router.GET("/inspector/api4transfers", transfersPage{}.transfersPageHandler)
}
