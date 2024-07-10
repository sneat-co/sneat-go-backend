package facade4userus

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/briefs4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/linkage/dbo4linkage"
	"github.com/sneat-co/sneat-go-backend/src/modules/teamus/facade4teamus"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/dbo4userus"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/dto4userus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/sneat-co/sneat-go-core/sneatauth"
	"github.com/strongo/strongoapp/person"
	"github.com/strongo/strongoapp/with"
	"strings"
	"time"
)

// InitUserRecord sets user title
func InitUserRecord(ctx context.Context, userContext facade.User, request dto4userus.InitUserRecordRequest) (user dbo4userus.UserEntry, err error) {
	if err = request.Validate(); err != nil {
		err = fmt.Errorf("%w: %v", facade.ErrBadRequest, err)
		return
	}
	userID := userContext.GetID()
	var userInfo *sneatauth.AuthUserInfo
	if userInfo, err = sneatauth.GetUserInfo(ctx, userID); err != nil {
		return user, fmt.Errorf("failed to get user info: %w", err)
	}
	err = runReadwriteTransaction(ctx, func(ctx context.Context, tx dal.ReadwriteTransaction) error {
		user, err = initUserRecordTxWorker(ctx, tx, userID, userInfo, request)
		return err
	})
	if err != nil {
		user.Data = nil
		return user, fmt.Errorf("failet to init user record: %w", err)
	}
	if request.Space != nil {
		var hasSpaceOfSameType bool
		for _, team := range user.Data.Spaces {
			if team.Type == request.Space.Type {
				hasSpaceOfSameType = true
				break
			}
		}
		if !hasSpaceOfSameType && request.Space != nil {
			if _, err = facade4teamus.CreateSpace(ctx, userContext, *request.Space); err != nil {
				err = fmt.Errorf("failed to create team for user: %w", err)
				return
			}
		}
	}

	return
}

func initUserRecordTxWorker(ctx context.Context, tx dal.ReadwriteTransaction, uid string, userInfo *sneatauth.AuthUserInfo, request dto4userus.InitUserRecordRequest) (user dbo4userus.UserEntry, err error) {
	var isNewUser bool
	user = dbo4userus.NewUserEntry(uid)
	if err = TxGetUserByID(ctx, tx, user.Record); err != nil {
		if dal.IsNotFound(err) {
			isNewUser = true
		} else {
			return
		}
	}
	if isNewUser {
		if err = createUserRecordTx(ctx, tx, request, user, userInfo); err != nil {
			err = fmt.Errorf("faield to create user record: %w", err)
			return
		}
	} else if err = updateUserRecordWithInitData(ctx, tx, request, user); err != nil {
		err = fmt.Errorf("faield to update user record: %w", err)
		return
	}
	return
}

func createUserRecordTx(ctx context.Context, tx dal.ReadwriteTransaction, request dto4userus.InitUserRecordRequest, user dbo4userus.UserEntry, userInfo *sneatauth.AuthUserInfo) error {
	user.Data.Status = "active"
	user.Data.Type = briefs4contactus.ContactTypePerson
	user.Data.CountryID = with.UnknownCountryID
	user.Data.AgeGroup = "unknown"
	user.Data.Gender = "unknown"

	if request.Names != nil && !request.Names.IsEmpty() {
		user.Data.Names = request.Names
	}

	if user.Data.Names != nil && user.Data.Names.FullName != "" && (user.Data.Names.FirstName == "" || user.Data.Names.LastName == "") {
		firstName, lastName := person.DeductNamesFromFullName(user.Data.Names.FullName)
		if user.Data.Names.FirstName == "" || firstName != "" {
			user.Data.Names.FirstName = firstName
		}
		if user.Data.Names.LastName == "" || lastName != "" {
			user.Data.Names.LastName = lastName
		}
	}

	user.Data.CreatedAt = time.Now()
	user.Data.CreatedBy = request.RemoteClient.HostOrApp
	if i := strings.Index(user.Data.CreatedBy, ":"); i > 0 {
		user.Data.CreatedBy = user.Data.CreatedBy[:i]
	}
	user.Data.Created.Client = request.RemoteClient
	if request.Email != "" {
		user.Data.Email = request.Email
		user.Data.EmailVerified = request.EmailIsVerified
	} else {
		user.Data.Email = userInfo.Email
		user.Data.EmailVerified = userInfo.EmailVerified
	}
	authProvider := request.AuthProvider
	if authProvider == "" {
		if len(userInfo.ProviderUserInfo) == 1 {
			authProvider = userInfo.ProviderUserInfo[0].ProviderID
		} else {
			authProvider = userInfo.ProviderID
		}
	}
	user.Data.Emails = []dbmodels.PersonEmail{
		{
			Type:         "primary",
			Address:      user.Data.Email,
			Verified:     user.Data.EmailVerified,
			AuthProvider: authProvider,
		},
	}
	if request.IanaTimezone != "" {
		user.Data.Timezone = &dbmodels.Timezone{
			Iana: request.IanaTimezone,
		}
	}
	if user.Data.Title == "" && user.Data.Names.IsEmpty() {
		user.Data.Title = user.Data.Email
	}
	_ = dbo4linkage.UpdateRelatedIDs(&user.Data.WithRelated, &user.Data.WithRelatedIDs)
	if err := user.Data.Validate(); err != nil {
		return fmt.Errorf("user record prepared for insert is not valid: %w", err)
	}
	if err := tx.Insert(ctx, user.Record); err != nil {
		return fmt.Errorf("failed to insert user record: %w", err)
	}
	return nil
}

func updateUserRecordWithInitData(ctx context.Context, tx dal.ReadwriteTransaction, request dto4userus.InitUserRecordRequest, user dbo4userus.UserEntry) error {
	var updates []dal.Update
	if name := request.Names; name != nil {
		if name.FullName == "" && !name.IsEmpty() {
			name.FullName = name.GetFullName()
		}
		if !name.IsEmpty() {
			updates = append(updates, dal.Update{Field: "name", Value: name})
		}
		user.Data.Names = name
	}

	if request.IanaTimezone != "" && (user.Data.Timezone == nil || user.Data.Timezone.Iana == "") {
		if user.Data.Timezone == nil {
			user.Data.Timezone = &dbmodels.Timezone{}
		}
		user.Data.Timezone.Iana = request.IanaTimezone
		updates = append(updates, dal.Update{Field: "timezone.iana", Value: request.IanaTimezone})
	}
	if user.Data.Title == user.Data.Email && user.Data.Names != nil && !user.Data.Names.IsEmpty() {
		user.Data.Title = ""
		updates = append(updates, dal.Update{Field: "title", Value: dal.DeleteField})
	}
	if len(updates) > 0 {
		if err := user.Data.Validate(); err != nil {
			return fmt.Errorf("user record prepared for update is not valid: %w", err)
		}
		if err := tx.Update(ctx, user.Key, updates); err != nil {
			return err
		}
	}
	return nil
}
