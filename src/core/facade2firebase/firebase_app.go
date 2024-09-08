package facade2firebase

import (
	"context"
	firebase "firebase.google.com/go/v4"
	core "github.com/sneat-co/sneat-go-core"
)

func GetFirebaseApp(ctx context.Context) (app *firebase.App, err error) {
	var config *firebase.Config
	if !core.IsInProd() {
		config = &firebase.Config{
			ServiceAccountID: "LOCAL-SNEAT-APP@my-project-id.iam.gserviceaccount.com",
		}
	}
	return firebase.NewApp(ctx, config)
}
