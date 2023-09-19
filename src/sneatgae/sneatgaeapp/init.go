package sneatgaeapp

import (
	"github.com/sneat-co/sneat-go-firebase/sneatfb"
	"os"
	"strings"
)

func init() {
	sneatfb.InitFirebaseForSneat(projectID, "sneat")
}

func getProjectID() string {
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
		}
	}
	return "demo-local-sneat-app"
}

var projectID = getProjectID()
