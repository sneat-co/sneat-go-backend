package sneatgaeapp

import (
	"cloud.google.com/go/firestore"
	"context"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo2firestore"
	"github.com/sneat-co/sneat-go-core/facade"
	"os"
	"strings"
)

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

func init() {
	facade.GetDatabase = func(ctx context.Context) dal.DB {

		client, err := firestore.NewClient(ctx, projectID)
		if err != nil {
			panic(err)
		}
		return dalgo2firestore.NewDatabase("sneat", client)
	}
}
