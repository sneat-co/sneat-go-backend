package facade4auth

import (
	"context"
	"firebase.google.com/go/v4/auth"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/core/facade2firebase"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/dbo4userus"
)

func createUser(ctx context.Context, tx dal.ReadwriteTransaction, userToCreate *auth.UserToCreate, customClaims map[string]interface{}) (uid string, err error) {

	var fbAuthClient *auth.Client
	if fbAuthClient, err = facade2firebase.GetFirebaseAuthClient(ctx); err != nil {
		return
	}
	var fbUserRecord *auth.UserRecord
	if fbUserRecord, err = fbAuthClient.CreateUser(ctx, userToCreate); err != nil {
		err = fmt.Errorf("failed to create firebase user: %w", err)
		return
	}
	uid = fbUserRecord.UID
	if err = fbAuthClient.SetCustomUserClaims(ctx, uid, customClaims); err != nil {
		err = fmt.Errorf("failed to set custom claims: %w", err)
		return
	}
	user := dbo4userus.NewUserEntry(uid)
	if err = tx.Insert(ctx, user.Record); err != nil {
		err = fmt.Errorf("failed to insert user record: %w", err)
		return
	}
	return
}
