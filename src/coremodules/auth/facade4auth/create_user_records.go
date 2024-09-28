package facade4auth

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/coremodules/contactus/briefs4contactus"
	"github.com/sneat-co/sneat-go-backend/src/coremodules/userus/dal4userus"
	"github.com/sneat-co/sneat-go-backend/src/coremodules/userus/dbo4userus"
	"github.com/sneat-co/sneat-go-backend/src/modules/linkage/dbo4linkage"
	"github.com/sneat-co/sneat-go-core/dto4auth"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/sneat-co/sneat-go-core/sneatauth"
	"github.com/strongo/strongoapp/appuser"
	"github.com/strongo/strongoapp/person"
	"github.com/strongo/strongoapp/with"
	"strings"
	"time"
)

// CreateUserRecords sets user title
func CreateUserRecords(ctx context.Context, userCtx facade.UserContext, userToCreate dto4auth.DataToCreateUser) (params CreateUserWorkerParams, err error) {
	if err = userToCreate.Validate(); err != nil {
		err = fmt.Errorf("%w: %v", facade.ErrBadRequest, err)
		return
	}
	userID := userCtx.GetUserID()
	var userInfo *sneatauth.AuthUserInfo
	if userInfo, err = sneatauth.GetUserInfo(ctx, userID); err != nil {
		err = fmt.Errorf("failed to get user info: %w", err)
		return
	}

	err = dal4userus.RunUserWorker(ctx, userCtx, false, func(ctx context.Context, tx dal.ReadwriteTransaction, userWorkerParams *dal4userus.UserWorkerParams) (err error) {
		params = CreateUserWorkerParams{
			UserWorkerParams: userWorkerParams,
		}
		if err = createUserRecordsTxWorker(ctx, tx, userInfo, userToCreate, &params); err != nil {
			return
		}
		if err = params.ApplyChanges(ctx, tx); err != nil {
			err = fmt.Errorf("failed to apply changes returned by createUserRecordsTxWorker(): %w", err)
		}
		return
	})
	if err != nil {
		return params, fmt.Errorf("failed to init user record and to create default user spaces: %w", err)
	}
	return
}

func createUserRecordsTxWorker(
	ctx context.Context,
	tx dal.ReadwriteTransaction,
	userInfo *sneatauth.AuthUserInfo, userToCreate dto4auth.DataToCreateUser, // TODO: Does this 2 duplicate each other?
	params *CreateUserWorkerParams,
) (err error) {
	if params == nil {
		panic("params is nil")
	}
	if userInfo == nil {
		panic("userInfo is nil")
	}
	if err = createOrUpdateUserRecord(userInfo, userToCreate, params); err != nil {
		return
	}

	if !params.User.Record.Exists() {
		if err = createDefaultUserSpacesTx(ctx, tx, params); err != nil {
			return fmt.Errorf("failed to create default user spaces: %w", err)
		}
	}
	return
}

func createOrUpdateUserRecord(userInfo *sneatauth.AuthUserInfo, userToCreate dto4auth.DataToCreateUser, params *CreateUserWorkerParams) (err error) {
	if params == nil {
		panic("params is nil")
	}
	if !params.User.Record.Exists() {
		if err = createUserRecord(userToCreate, params.User, userInfo); err != nil {
			err = fmt.Errorf("failed to populate new user record data: %w", err)
			return
		}
		params.User.Record.MarkAsChanged()
		params.QueueForInsert(params.User.Record)
	} else if err = updateUserRecordWithInitData(userToCreate, params.UserWorkerParams); err != nil {
		err = fmt.Errorf("failed to update user record data: %w", err)
		// It might be too earlier to add updates to RecordsToUpdate?
		//params.RecordsToUpdate = append(params.RecordsToUpdate, record.Updates{Record: params.User.Record, Updates: params.UserUpdates})
		return
	}
	return
}

func createUserRecord(userToCreate dto4auth.DataToCreateUser, user dbo4userus.UserEntry, userInfo *sneatauth.AuthUserInfo) error {
	if userInfo == nil {
		panic("userInfo is nil")
	}
	user.Data.Status = "active"
	user.Data.ContactBrief.Type = briefs4contactus.ContactTypePerson
	user.Data.ContactBrief.AgeGroup = "unknown"
	user.Data.ContactBrief.Gender = "unknown"
	user.Data.OptionalCountryID.CountryID = with.UnknownCountryID

	if !userToCreate.Names.IsEmpty() {
		user.Data.ContactBrief.Names = &userToCreate.Names
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
	user.Data.CreatedBy = userToCreate.RemoteClient.HostOrApp
	if i := strings.Index(user.Data.CreatedBy, ":"); i > 0 {
		user.Data.CreatedBy = user.Data.CreatedBy[:i]
	}
	user.Data.Created.Client = userToCreate.RemoteClient
	if userToCreate.Email != "" {
		user.Data.Email = userToCreate.Email
		user.Data.EmailVerified = userToCreate.EmailIsVerified
	} else {
		user.Data.Email = userInfo.Email
		user.Data.EmailVerified = userInfo.EmailVerified
	}

	if userToCreate.AuthAccount.IsEmpty() {
		if len(userInfo.ProviderUserInfo) == 1 {
			ui := userInfo.ProviderUserInfo[0]
			userToCreate.AuthAccount = appuser.AccountKey{
				Provider: ui.ProviderID,
			}
			if ui.Email != "" {
				userToCreate.AuthAccount.ID = ui.Email
			} else if ui.PhoneNumber != "" {
				userToCreate.AuthAccount.ID = ui.PhoneNumber
			}
		}
	}

	_ = user.Data.AddAccount(userToCreate.AuthAccount)

	if user.Data.Email != "" {
		user.Data.Emails = []dbmodels.PersonEmail{
			{
				Type:         "primary",
				Address:      user.Data.Email,
				Verified:     user.Data.EmailVerified,
				AuthProvider: userToCreate.AuthAccount.Provider,
			},
		}
	}
	if userToCreate.IanaTimezone != "" {
		user.Data.Timezone = &dbmodels.Timezone{
			Iana: userToCreate.IanaTimezone,
		}
	}
	if user.Data.Title == "" && user.Data.Names.IsEmpty() {
		user.Data.Title = user.Data.Email
	}
	_ = dbo4linkage.UpdateRelatedIDs(&user.Data.WithRelated, &user.Data.WithRelatedIDs)
	if err := user.Data.Validate(); err != nil {
		return fmt.Errorf("user record prepared for insert is not valid: %w", err)
	}
	return nil
}

func updateUserRecordWithInitData(userToCreate dto4auth.DataToCreateUser, params *dal4userus.UserWorkerParams) error {
	if name := userToCreate.Names; !name.IsEmpty() {
		if name.FullName == "" {
			name.FullName = name.GetFullName()
		}
		if !name.IsEmpty() {
			params.UserUpdates = append(params.UserUpdates, dal.Update{Field: "name", Value: name})
			params.User.Record.MarkAsChanged()
		}
		params.User.Data.Names = &name
	}

	if userToCreate.IanaTimezone != "" && (params.User.Data.Timezone == nil || params.User.Data.Timezone.Iana == "") {
		if params.User.Data.Timezone == nil {
			params.User.Data.Timezone = &dbmodels.Timezone{}
		}
		params.User.Data.Timezone.Iana = userToCreate.IanaTimezone
		params.UserUpdates = append(params.UserUpdates, dal.Update{Field: "timezone.iana", Value: userToCreate.IanaTimezone})
		params.User.Record.MarkAsChanged()
	}
	if params.User.Data.Title == params.User.Data.Email && params.User.Data.Names != nil && !params.User.Data.Names.IsEmpty() {
		params.User.Data.Title = ""
		params.UserUpdates = append(params.UserUpdates, dal.Update{Field: "title", Value: dal.DeleteField})
		params.User.Record.MarkAsChanged()
	}
	return nil
}
