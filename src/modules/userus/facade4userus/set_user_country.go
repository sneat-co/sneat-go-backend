package facade4userus

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/dal4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/spaceus/dto4spaceus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/strongoapp/with"
	"github.com/strongo/validation"
)

type SetUserCountryRequest struct {
	CountryID string `json:"countryID"`
}

func (v SetUserCountryRequest) Validate() error {
	if v.CountryID == "" {
		return validation.NewErrRequestIsMissingRequiredField("countryID")
	}
	if len(v.CountryID) != 2 {
		return validation.NewErrBadRequestFieldValue("countryID", "must be 2 characters long")
	}
	return nil
}

func SetUserCountry(ctx context.Context, userContext facade.User, request SetUserCountryRequest) (err error) {
	return RunUserWorker(ctx, userContext, func(ctx context.Context, tx dal.ReadwriteTransaction, params *UserWorkerParams) error {
		return txSetUserCountry(ctx, userContext, request, tx, params)
	})
}

func txSetUserCountry(ctx context.Context, userContext facade.User, request SetUserCountryRequest, tx dal.ReadwriteTransaction, params *UserWorkerParams) (err error) {
	if params.User.Data.CountryID != request.CountryID {
		params.User.Data.CountryID = request.CountryID
		params.UserUpdates = append(params.UserUpdates,
			dal.Update{Field: "countryID", Value: request.CountryID})
	}
	for teamID, teamBrief := range params.User.Data.Spaces {
		if teamBrief.Type == "family" && IsUnknownCountryID(teamBrief.CountryID) {
			teamBrief.CountryID = request.CountryID
			params.UserUpdates = append(params.UserUpdates, dal.Update{Field: fmt.Sprintf("spaces.%s.countryID", teamID), Value: request.CountryID})
		}
		teamRequest := dto4spaceus.SpaceRequest{SpaceID: teamID}
		err = dal4contactus.RunContactusSpaceWorkerTx(ctx, tx, userContext, teamRequest, func(ctx context.Context, tx dal.ReadwriteTransaction, params *dal4contactus.ContactusSpaceWorkerParams) error {
			if err = params.GetRecords(ctx, tx); err != nil {
				return err
			}
			if params.Space.Data.CountryID == "" || params.Space.Data.CountryID == with.UnknownCountryID {
				params.Space.Data.CountryID = request.CountryID
				params.SpaceUpdates = append(params.SpaceUpdates, dal.Update{Field: "countryID", Value: request.CountryID})
			}
			userContactID, userContactBrief := params.SpaceModuleEntry.Data.GetContactBriefByUserID(userContext.GetID())
			if userContactBrief != nil && IsUnknownCountryID(userContactBrief.CountryID) {
				userContactBrief.CountryID = request.CountryID
				params.SpaceModuleUpdates = append(params.SpaceModuleUpdates, dal.Update{Field: "contacts." + userContactID + ".countryID", Value: request.CountryID})
				userContact := dal4contactus.NewContactEntry(teamID, userContactID)
				if err = tx.Get(ctx, userContact.Record); err != nil {
					return err
				}
				if IsUnknownCountryID(userContact.Data.CountryID) {
					userContact.Data.CountryID = request.CountryID
					if err = tx.Update(ctx, userContact.Key, []dal.Update{{Field: "countryID", Value: request.CountryID}}); err != nil {
						return err
					}
				}
			}
			return nil
		})
	}
	return nil
}

// IsUnknownCountryID checks if countryID is empty or "--" - TODO: move next to dbmodels.UnknownCountryID
func IsUnknownCountryID(countryID string) bool {
	return countryID == "" || countryID == with.UnknownCountryID
}
