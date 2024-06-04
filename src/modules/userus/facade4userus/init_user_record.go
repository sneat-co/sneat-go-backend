package facade4userus

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/briefs4contactus"
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
func InitUserRecord(ctx context.Context, userContext facade.User, request dto4userus.InitUserRecordRequest) (user dbo4userus.UserContext, err error) {
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
		user.Dbo = nil
		return user, fmt.Errorf("failet to init user record: %w", err)
	}
	if request.Team != nil {
		var hasTeamOfSameType bool
		for _, team := range user.Dbo.Teams {
			if team.Type == request.Team.Type {
				hasTeamOfSameType = true
				break
			}
		}
		if !hasTeamOfSameType && request.Team != nil {
			if _, err = facade4teamus.CreateTeam(ctx, userContext, *request.Team); err != nil {
				err = fmt.Errorf("failed to create team for user: %w", err)
				return
			}
		}
	}

	return
}

func initUserRecordTxWorker(ctx context.Context, tx dal.ReadwriteTransaction, uid string, userInfo *sneatauth.AuthUserInfo, request dto4userus.InitUserRecordRequest) (user dbo4userus.UserContext, err error) {
	var isNewUser bool
	user = dbo4userus.NewUserContext(uid)
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

func createUserRecordTx(ctx context.Context, tx dal.ReadwriteTransaction, request dto4userus.InitUserRecordRequest, user dbo4userus.UserContext, userInfo *sneatauth.AuthUserInfo) error {
	user.Dbo.Status = "active"
	user.Dbo.Type = briefs4contactus.ContactTypePerson
	user.Dbo.CountryID = with.UnknownCountryID
	user.Dbo.AgeGroup = "unknown"
	user.Dbo.Gender = "unknown"

	if request.Names != nil && !request.Names.IsEmpty() {
		user.Dbo.Names = request.Names
	}

	if user.Dbo.Names != nil && user.Dbo.Names.FullName != "" && (user.Dbo.Names.FirstName == "" || user.Dbo.Names.LastName == "") {
		firstName, lastName := person.DeductNamesFromFullName(user.Dbo.Names.FullName)
		if user.Dbo.Names.FirstName == "" || firstName != "" {
			user.Dbo.Names.FirstName = firstName
		}
		if user.Dbo.Names.LastName == "" || lastName != "" {
			user.Dbo.Names.LastName = lastName
		}
	}

	user.Dbo.CreatedAt = time.Now()
	user.Dbo.CreatedBy = request.RemoteClient.HostOrApp
	if i := strings.Index(user.Dbo.CreatedBy, ":"); i > 0 {
		user.Dbo.CreatedBy = user.Dbo.CreatedBy[:i]
	}
	user.Dbo.Created.Client = request.RemoteClient
	if request.Email != "" {
		user.Dbo.Email = request.Email
		user.Dbo.EmailVerified = request.EmailIsVerified
	} else {
		user.Dbo.Email = userInfo.Email
		user.Dbo.EmailVerified = userInfo.EmailVerified
	}
	authProvider := request.AuthProvider
	if authProvider == "" {
		if len(userInfo.ProviderUserInfo) == 1 {
			authProvider = userInfo.ProviderUserInfo[0].ProviderID
		} else {
			authProvider = userInfo.ProviderID
		}
	}
	user.Dbo.Emails = []dbmodels.PersonEmail{
		{
			Type:         "primary",
			Address:      user.Dbo.Email,
			Verified:     user.Dbo.EmailVerified,
			AuthProvider: authProvider,
		},
	}
	if request.IanaTimezone != "" {
		user.Dbo.Timezone = &dbmodels.Timezone{
			Iana: request.IanaTimezone,
		}
	}
	if user.Dbo.Title == "" && user.Dbo.Names.IsEmpty() {
		user.Dbo.Title = user.Dbo.Email
	}
	_ = user.Dbo.UpdateRelatedIDs()
	if err := user.Dbo.Validate(); err != nil {
		return fmt.Errorf("user record prepared for insert is not valid: %w", err)
	}
	if err := tx.Insert(ctx, user.Record); err != nil {
		return fmt.Errorf("failed to insert user record: %w", err)
	}
	return nil
}

func updateUserRecordWithInitData(ctx context.Context, tx dal.ReadwriteTransaction, request dto4userus.InitUserRecordRequest, user dbo4userus.UserContext) error {
	var updates []dal.Update
	if name := request.Names; name != nil {
		if name.FullName == "" && !name.IsEmpty() {
			name.FullName = name.GetFullName()
		}
		if !name.IsEmpty() {
			updates = append(updates, dal.Update{Field: "name", Value: name})
		}
		user.Dbo.Names = name
	}

	if request.IanaTimezone != "" && (user.Dbo.Timezone == nil || user.Dbo.Timezone.Iana == "") {
		if user.Dbo.Timezone == nil {
			user.Dbo.Timezone = &dbmodels.Timezone{}
		}
		user.Dbo.Timezone.Iana = request.IanaTimezone
		updates = append(updates, dal.Update{Field: "timezone.iana", Value: request.IanaTimezone})
	}
	if user.Dbo.Title == user.Dbo.Email && user.Dbo.Names != nil && !user.Dbo.Names.IsEmpty() {
		user.Dbo.Title = ""
		updates = append(updates, dal.Update{Field: "title", Value: dal.DeleteField})
	}
	if len(updates) > 0 {
		if err := user.Dbo.Validate(); err != nil {
			return fmt.Errorf("user record prepared for update is not valid: %w", err)
		}
		if err := tx.Update(ctx, user.Key, updates); err != nil {
			return err
		}
	}
	return nil
}
