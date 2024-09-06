package sneatgaeapp

import (
	"github.com/sneat-co/sneat-go-firebase/sneatfb"
	"os"
	"strings"
)

func init() {
	sneatfb.InitFirebaseForSneat(projectID, "sneat")
}

func getFirebaseProjectID() string {
	if fbProjID := os.Getenv("FIREBASE_PROJECT_ID"); fbProjID != "" {
		return fbProjID
	}
	for _, v := range os.Environ() {
		if strings.HasPrefix(v, "GAE_APPLICATION=") {
			if strings.HasSuffix(v, "sneat-team") {
				return "sneat-team"
			}
			if strings.HasSuffix(v, "sneat-eu") {
				return "sneat-eu"
			}
			if strings.HasSuffix(v, "sneatapp") {
				return "sneatapp"
			}
			if strings.HasSuffix(v, "sneat-eur3-1") {
				return "sneat-eur3-1"
			}
		}
	}
	return "demo-local-sneat-app"
}

var projectID = getFirebaseProjectID()
