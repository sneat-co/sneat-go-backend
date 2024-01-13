package api4teamus

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	//sneatfb.NewFirestoreContext = func(r *http.Request, authRequired bool) (context *sneatfb.FirestoreContext, err error) {
	//	return
	//}

	os.Exit(m.Run())
}
