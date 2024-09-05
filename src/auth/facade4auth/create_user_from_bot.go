package facade4auth

import (
	"context"
	"firebase.google.com/go/v4/auth"
	"fmt"
	"github.com/sneat-co/sneat-go-backend/src/core/facade2firebase"
)

func createFirebaseUser(ctx context.Context, userToCreate DataToCreateUser) (firebaseUserRecord *auth.UserRecord, err error) {
	displayName := userToCreate.Names.FirstName
	if userToCreate.Names.LastName != "" {
		if displayName != "" {
			displayName += " "
		}
		displayName += userToCreate.Names.LastName
	}
	if displayName == "" {
		displayName = userToCreate.Names.UserName
	}
	firebaseUserToCreate := (&auth.UserToCreate{}).
		DisplayName(displayName)
	if userToCreate.PhotoURL != "" {
		firebaseUserToCreate = firebaseUserToCreate.PhotoURL(userToCreate.PhotoURL)
	}

	var fbAuthClient *auth.Client
	if fbAuthClient, err = facade2firebase.GetFirebaseAuthClient(ctx); err != nil {
		return
	}

	if firebaseUserRecord, err = fbAuthClient.CreateUser(ctx, firebaseUserToCreate); err != nil {
		err = fmt.Errorf("failed to create firebase user: %w", err)
		return
	}
	customClaims := map[string]interface{}{
		"authProvider": userToCreate.Account.Provider,
	}
	if err = fbAuthClient.SetCustomUserClaims(ctx, firebaseUserRecord.UID, customClaims); err != nil {
		err = fmt.Errorf("failed to set custom claims: %w", err)
		return
	}
	return
}

func deleteFirebaseUser(ctx context.Context, uid string) (err error) {
	var fbAuthClient *auth.Client
	if fbAuthClient, err = facade2firebase.GetFirebaseAuthClient(ctx); err != nil {
		return
	}
	if err = fbAuthClient.DeleteUser(ctx, uid); err != nil {
		err = fmt.Errorf("failed to delete firebase user: %w", err)
	}
	return
}
