package facade4userus

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/dal4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/spaceus/core4spaceus"
	"github.com/sneat-co/sneat-go-backend/src/modules/spaceus/dto4spaceus"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/dal4userus"
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

func SetUserCountry(ctx context.Context, userCtx facade.UserContext, request SetUserCountryRequest) (err error) {
	return dal4userus.RunUserWorker(ctx, userCtx, true, func(ctx context.Context, tx dal.ReadwriteTransaction, params *dal4userus.UserWorkerParams) (err error) {
		if err = txSetUserCountry(ctx, tx, userCtx, request, params); err != nil {
			return fmt.Errorf("failed in txSetUserCountry(): %w", err)
		}
		return
	})
}

type RecordToUpdate struct {
	Key     *dal.Key
	Updates []dal.Update
}

func txSetUserCountry(ctx context.Context, tx dal.ReadwriteTransaction, userCtx facade.UserContext, request SetUserCountryRequest, params *dal4userus.UserWorkerParams) (err error) {
	if params.User.Data.CountryID != request.CountryID {
		params.User.Data.CountryID = request.CountryID
		params.User.Record.MarkAsChanged()
		params.UserUpdates = append(params.UserUpdates,
			dal.Update{Field: "countryID", Value: request.CountryID})
	}

	recordsToUpdate := make([]RecordToUpdate, 0, len(params.User.Data.Spaces))

	for spaceID, spaceBrief := range params.User.Data.Spaces {
		if IsUnknownCountryID(spaceBrief.CountryID) && spaceBrief.Type == core4spaceus.SpaceTypeFamily || spaceBrief.Type == core4spaceus.SpaceTypePrivate {
			spaceBrief.CountryID = request.CountryID
			params.UserUpdates = append(params.UserUpdates, dal.Update{Field: fmt.Sprintf("spaces.%s.countryID", spaceID), Value: request.CountryID})
		}
		spaceRequest := dto4spaceus.SpaceRequest{SpaceID: spaceID}
		if err = dal4contactus.RunContactusSpaceWorkerNoUpdate(ctx, tx, userCtx, spaceRequest,
			func(ctx context.Context, tx dal.ReadwriteTransaction, params *dal4contactus.ContactusSpaceWorkerParams) (err error) {
				if err = params.GetRecords(ctx, tx, params.Space.Record); err != nil {
					return
				}
				if IsUnknownCountryID(params.Space.Data.CountryID) {
					params.Space.Data.CountryID = request.CountryID
					params.SpaceUpdates = append(params.SpaceUpdates, dal.Update{Field: "countryID", Value: request.CountryID})
					params.Space.Record.MarkAsChanged()
				}
				userContactID, userContactBrief := params.SpaceModuleEntry.Data.GetContactBriefByUserID(userCtx.GetUserID())
				if userContactBrief != nil && IsUnknownCountryID(userContactBrief.CountryID) {
					userContactBrief.CountryID = request.CountryID
					params.SpaceModuleUpdates = append(params.SpaceModuleUpdates, dal.Update{Field: "contacts." + userContactID + ".countryID", Value: request.CountryID})
					params.SpaceModuleEntry.Record.MarkAsChanged()
					userContact := dal4contactus.NewContactEntry(spaceID, userContactID)
					if err = tx.Get(ctx, userContact.Record); err != nil {
						return
					}
					if IsUnknownCountryID(userContact.Data.CountryID) {
						userContact.Data.CountryID = request.CountryID
						recordsToUpdate = append(recordsToUpdate, RecordToUpdate{Key: userContact.Key, Updates: []dal.Update{{Field: "countryID", Value: request.CountryID}}})
					}
				}
				if params.Space.Record.HasChanged() && len(params.SpaceUpdates) > 0 {
					recordsToUpdate = append(recordsToUpdate, RecordToUpdate{Key: params.Space.Key, Updates: params.SpaceUpdates})
				}
				if params.SpaceModuleEntry.Record.HasChanged() && len(params.SpaceModuleUpdates) > 0 {
					recordsToUpdate = append(recordsToUpdate, RecordToUpdate{Key: params.SpaceModuleEntry.Key, Updates: params.SpaceModuleUpdates})
				}
				return
			}); err != nil {
			return fmt.Errorf("failed to update space %s: %w", spaceID, err)
		}
	}
	if len(recordsToUpdate) > 0 {
		for _, rec := range recordsToUpdate {
			if err = tx.Update(ctx, rec.Key, rec.Updates); err != nil {
				return fmt.Errorf("failed to update record %s: %w", rec.Key, err)
			}
		}
	}
	return
}

// IsUnknownCountryID checks if countryID is empty or "--" - TODO: move next to dbmodels.UnknownCountryID
func IsUnknownCountryID(countryID string) bool {
	return countryID == "" || countryID == with.UnknownCountryID
}
