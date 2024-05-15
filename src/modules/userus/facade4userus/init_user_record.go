package facade4userus

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/briefs4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/teamus/facade4teamus"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/dto4userus"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/models4userus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/sneat-co/sneat-go-core/sneatauth"
	"github.com/strongo/strongoapp/person"
	"github.com/strongo/strongoapp/with"
	"strings"
	"time"
)

// InitUserRecord sets user title
func InitUserRecord(ctx context.Context, userContext facade.User, request dto4userus.InitUserRecordRequest) (user models4userus.UserContext, err error) {
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
		user.Dto = nil
		return user, fmt.Errorf("failet to init user record: %w", err)
	}
	if request.Team != nil {
		var hasTeamOfSameType bool
		for _, team := range user.Dto.Teams {
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

func initUserRecordTxWorker(ctx context.Context, tx dal.ReadwriteTransaction, uid string, userInfo *sneatauth.AuthUserInfo, request dto4userus.InitUserRecordRequest) (user models4userus.UserContext, err error) {
	var isNewUser bool
	user = models4userus.NewUserContext(uid)
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

func createUserRecordTx(ctx context.Context, tx dal.ReadwriteTransaction, request dto4userus.InitUserRecordRequest, user models4userus.UserContext, userInfo *sneatauth.AuthUserInfo) error {
	user.Dto.Status = "active"
	user.Dto.Type = briefs4contactus.ContactTypePerson
	user.Dto.CountryID = with.UnknownCountryID
	user.Dto.AgeGroup = "unknown"
	user.Dto.Gender = "unknown"

	if request.Names != nil && !request.Names.IsEmpty() {
		user.Dto.Names = request.Names
	}

	if user.Dto.Names != nil && user.Dto.Names.FullName != "" && (user.Dto.Names.FirstName == "" || user.Dto.Names.LastName == "") {
		firstName, lastName := person.DeductNamesFromFullName(user.Dto.Names.FullName)
		if user.Dto.Names.FirstName == "" || firstName != "" {
			user.Dto.Names.FirstName = firstName
		}
		if user.Dto.Names.LastName == "" || lastName != "" {
			user.Dto.Names.LastName = lastName
		}
	}

	user.Dto.CreatedAt = time.Now()
	user.Dto.CreatedBy = request.RemoteClient.HostOrApp
	if i := strings.Index(user.Dto.CreatedBy, ":"); i > 0 {
		user.Dto.CreatedBy = user.Dto.CreatedBy[:i]
	}
	user.Dto.Created.Client = request.RemoteClient
	if request.Email != "" {
		user.Dto.Email = request.Email
		user.Dto.EmailVerified = request.EmailIsVerified
	} else {
		user.Dto.Email = userInfo.Email
		user.Dto.EmailVerified = userInfo.EmailVerified
	}
	authProvider := request.AuthProvider
	if authProvider == "" {
		if len(userInfo.ProviderUserInfo) == 1 {
			authProvider = userInfo.ProviderUserInfo[0].ProviderID
		} else {
			authProvider = userInfo.ProviderID
		}
	}
	user.Dto.Emails = []dbmodels.PersonEmail{
		{
			Type:         "primary",
			Address:      user.Dto.Email,
			Verified:     user.Dto.EmailVerified,
			AuthProvider: authProvider,
		},
	}
	if request.IanaTimezone != "" {
		user.Dto.Timezone = &dbmodels.Timezone{
			Iana: request.IanaTimezone,
		}
	}
	if user.Dto.Title == "" && user.Dto.Names.IsEmpty() {
		user.Dto.Title = user.Dto.Email
	}
	if err := user.Dto.Validate(); err != nil {
		return fmt.Errorf("user record prepared for insert is not valid: %w", err)
	}
	if err := tx.Insert(ctx, user.Record); err != nil {
		return fmt.Errorf("failed to insert user record: %w", err)
	}
	return nil
}

func updateUserRecordWithInitData(ctx context.Context, tx dal.ReadwriteTransaction, request dto4userus.InitUserRecordRequest, user models4userus.UserContext) error {
	var updates []dal.Update
	if name := request.Names; name != nil {
		if name.FullName == "" && !name.IsEmpty() {
			name.FullName = name.GetFullName()
		}
		if !name.IsEmpty() {
			updates = append(updates, dal.Update{Field: "name", Value: name})
		}
		user.Dto.Names = name
	}

	if request.IanaTimezone != "" && (user.Dto.Timezone == nil || user.Dto.Timezone.Iana == "") {
		if user.Dto.Timezone == nil {
			user.Dto.Timezone = &dbmodels.Timezone{}
		}
		user.Dto.Timezone.Iana = request.IanaTimezone
		updates = append(updates, dal.Update{Field: "timezone.iana", Value: request.IanaTimezone})
	}
	if user.Dto.Title == user.Dto.Email && user.Dto.Names != nil && !user.Dto.Names.IsEmpty() {
		user.Dto.Title = ""
		updates = append(updates, dal.Update{Field: "title", Value: dal.DeleteField})
	}
	if len(updates) > 0 {
		if err := user.Dto.Validate(); err != nil {
			return fmt.Errorf("user record prepared for update is not valid: %w", err)
		}
		if err := tx.Update(ctx, user.Key, updates); err != nil {
			return err
		}
	}
	return nil
}
