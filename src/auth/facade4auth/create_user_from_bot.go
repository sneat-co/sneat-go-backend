package facade4auth

import (
	"context"
	"firebase.google.com/go/v4/auth"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/core/facade2firebase"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/briefs4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/dbo4userus"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"time"
)

func createUserFromBot(ctx context.Context, tx dal.ReadwriteTransaction, userToCreate DataToCreateUser) (user dbo4userus.UserEntry, err error) {

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
	fbUserToCreate := (&auth.UserToCreate{}).
		DisplayName(displayName)
	if userToCreate.PhotoURL != "" {
		fbUserToCreate = fbUserToCreate.PhotoURL(userToCreate.PhotoURL)
	}

	var fbAuthClient *auth.Client
	if fbAuthClient, err = facade2firebase.GetFirebaseAuthClient(ctx); err != nil {
		return
	}

	var fbUserRecord *auth.UserRecord
	if fbUserRecord, err = fbAuthClient.CreateUser(ctx, fbUserToCreate); err != nil {
		err = fmt.Errorf("failed to create firebase user: %w", err)
		return
	}
	user.ID = fbUserRecord.UID
	customClaims := map[string]interface{}{
		"authProvider": userToCreate.Account.Provider,
	}
	if err = fbAuthClient.SetCustomUserClaims(ctx, user.ID, customClaims); err != nil {
		err = fmt.Errorf("failed to set custom claims: %w", err)
		return
	}
	now := time.Now()
	user = dbo4userus.NewUserEntry(user.ID)
	user.Data.SetCreatedAt(now)
	user.Data.SetLastLoginAt(now)
	user.Data.Type = briefs4contactus.ContactTypePerson
	user.Data.Status = briefs4contactus.ContactStatusActive
	user.Data.AgeGroup = dbmodels.AgeGroupUnknown
	user.Data.Gender = dbmodels.GenderUnknown
	user.Data.Created = dbmodels.CreatedInfo{
		Client: userToCreate.RemoteClient,
	}
	if !userToCreate.Names.IsEmpty() {
		user.Data.Names = &userToCreate.Names
	}
	user.Data.PreferredLocale = userToCreate.LanguageCode
	user.Data.AddAccount(userToCreate.Account)

	if err = user.Data.Validate(); err != nil {
		err = fmt.Errorf("invalid new user record data: %w", err)
		return
	}

	if tx != nil {
		if err = tx.Insert(ctx, user.Record); err != nil {
			err = fmt.Errorf("failed to insert user record: %w", err)
			return
		}
	}
	return
}
