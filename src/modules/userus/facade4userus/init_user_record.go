package facade4userus

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/briefs4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/linkage/dbo4linkage"
	"github.com/sneat-co/sneat-go-backend/src/modules/spaceus/facade4spaceus"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/dal4userus"
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
func InitUserRecord(ctx context.Context, userCtx facade.UserContext, request dto4userus.InitUserRecordRequest) (user dbo4userus.UserEntry, err error) {
	if err = request.Validate(); err != nil {
		err = fmt.Errorf("%w: %v", facade.ErrBadRequest, err)
		return
	}
	userID := userCtx.GetUserID()
	var userInfo *sneatauth.AuthUserInfo
	if userInfo, err = sneatauth.GetUserInfo(ctx, userID); err != nil {
		return user, fmt.Errorf("failed to get user info: %w", err)
	}

	err = dal4userus.RunUserWorker(ctx, userCtx, false, func(ctx context.Context, tx dal.ReadwriteTransaction, params *dal4userus.UserWorkerParams) (err error) {
		user = params.User
		if err = initUserRecordTxWorker(userInfo, request, params); err != nil {
			return
		}

		if !user.Record.Exists() {
			if _, _, err = facade4spaceus.CreateDefaultUserSpacesTx(ctx, tx, params); err != nil {
				return fmt.Errorf("failed to create default user spaces: %w", err)
			}
		}

		user.Record.MarkAsChanged()

		return
	})
	if err != nil {
		user.Data = nil
		return user, fmt.Errorf("failed to init user record and to create default user spaces: %w", err)
	}

	return
}

func initUserRecordTxWorker(userInfo *sneatauth.AuthUserInfo, request dto4userus.InitUserRecordRequest, params *dal4userus.UserWorkerParams) (err error) {
	if !params.User.Record.Exists() {
		if err = createUserRecord(request, params.User, userInfo); err != nil {
			err = fmt.Errorf("faield to populate new user record data: %w", err)
			return
		}
	} else if err = updateUserRecordWithInitData(request, params); err != nil {
		err = fmt.Errorf("failed to update user record data: %w", err)
		return
	}
	return
}

func createUserRecord(request dto4userus.InitUserRecordRequest, user dbo4userus.UserEntry, userInfo *sneatauth.AuthUserInfo) error {
	user.Data.Status = "active"
	user.Data.ContactBrief.Type = briefs4contactus.ContactTypePerson
	user.Data.ContactBrief.AgeGroup = "unknown"
	user.Data.ContactBrief.Gender = "unknown"
	user.Data.OptionalCountryID.CountryID = with.UnknownCountryID

	if request.Names != nil && !request.Names.IsEmpty() {
		user.Data.ContactBrief.Names = request.Names
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
	return nil
}

func updateUserRecordWithInitData(request dto4userus.InitUserRecordRequest, params *dal4userus.UserWorkerParams) error {
	if name := request.Names; name != nil {
		if name.FullName == "" && !name.IsEmpty() {
			name.FullName = name.GetFullName()
		}
		if !name.IsEmpty() {
			params.UserUpdates = append(params.UserUpdates, dal.Update{Field: "name", Value: name})
		}
		params.User.Data.Names = name
	}

	if request.IanaTimezone != "" && (params.User.Data.Timezone == nil || params.User.Data.Timezone.Iana == "") {
		if params.User.Data.Timezone == nil {
			params.User.Data.Timezone = &dbmodels.Timezone{}
		}
		params.User.Data.Timezone.Iana = request.IanaTimezone
		params.UserUpdates = append(params.UserUpdates, dal.Update{Field: "timezone.iana", Value: request.IanaTimezone})
	}
	if params.User.Data.Title == params.User.Data.Email && params.User.Data.Names != nil && !params.User.Data.Names.IsEmpty() {
		params.User.Data.Title = ""
		params.UserUpdates = append(params.UserUpdates, dal.Update{Field: "title", Value: dal.DeleteField})
	}
	//if len(params.UserUpdates) > 0 {
	//	if err := params.User.Data.Validate(); err != nil {
	//		return fmt.Errorf("user record prepared for update is not valid: %w", err)
	//	}
	//	if err := tx.Update(ctx, params.User.Key, params.UserUpdates); err != nil {
	//		return err
	//	}
	//}
	return nil
}
