package facade2firebase

import (
	"context"
	firebase "firebase.google.com/go/v4"
	"fmt"
)

func GetFirebaseApp(ctx context.Context) (app *firebase.App, err error) {
	conf := &firebase.Config{
		ServiceAccountID: "my-client-id@my-project-id.iam.gserviceaccount.com",
	}
	if app, err = firebase.NewApp(ctx, conf); err != nil {
		err = fmt.Errorf("faield to initializing Firebase app: %w", err)
		return
	}
	return
}
