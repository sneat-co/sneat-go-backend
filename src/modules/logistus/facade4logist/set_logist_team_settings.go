package facade4logist

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/briefs4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/dal4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/dbo4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/dto4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/facade4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/logistus/const4logistus"
	"github.com/sneat-co/sneat-go-backend/src/modules/logistus/dbo4logist"
	"github.com/sneat-co/sneat-go-backend/src/modules/logistus/dto4logist"
	"github.com/sneat-co/sneat-go-backend/src/modules/teamus/dal4teamus"
	"github.com/sneat-co/sneat-go-backend/src/modules/teamus/dbo4teamus"
	"github.com/sneat-co/sneat-go-backend/src/modules/teamus/dto4teamus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/slice"
)

// SetLogistSpaceSettings sets team settings for logistus module
func SetLogistSpaceSettings(
	ctx context.Context,
	userContext facade.User,
	request dto4logist.SetLogistSpaceSettingsRequest,
) error {
	if err := request.Validate(); err != nil {
		return err
	}
	return dal4teamus.RunModuleSpaceWorker(ctx, userContext, request.SpaceRequest,
		const4logistus.ModuleID,
		new(dbo4logist.LogistSpaceDbo),
		func(ctx context.Context, tx dal.ReadwriteTransaction, teamWorkerParams *dal4teamus.ModuleSpaceWorkerParams[*dbo4logist.LogistSpaceDbo]) (err error) {
			return setLogistSpaceSettingsTx(ctx /*userContext,*/, request, tx, teamWorkerParams)
		},
	)
}

func setLogistSpaceSettingsTx(
	ctx context.Context,
	//userContext facade.User,
	request dto4logist.SetLogistSpaceSettingsRequest,
	tx dal.ReadwriteTransaction,
	workerParams *dal4teamus.ModuleSpaceWorkerParams[*dbo4logist.LogistSpaceDbo],
) (err error) {
	if workerParams.Space.Data.CountryID != request.Address.CountryID {
		workerParams.Space.Data.CountryID = request.Address.CountryID
		workerParams.SpaceUpdates = append(workerParams.SpaceUpdates, dal.Update{
			Field: "countryID",
			Value: request.Address.CountryID,
		})
	}

	logistSpace := dbo4logist.NewLogistSpaceEntry(request.SpaceID)
	if err = tx.Get(ctx, logistSpace.Record); err != nil {
		if !dal.IsNotFound(err) {
			return err
		}
	} else if err = logistSpace.Data.Validate(); err != nil {
		return fmt.Errorf("loaded logistus team recod is not valid: %w", err)
	}
	var teamContact dal4contactus.ContactEntry
	if teamContact, err = facade4contactus.GetContactByID(ctx, tx, logistSpace.ID, request.SpaceID); err != nil {
		if !dal.IsNotFound(err) {
			return fmt.Errorf("failed to get contact record: %w", err)
		}
	}
	if dal.IsNotFound(err) {
		createContactRequest := dto4contactus.CreateContactRequest{
			Status:       "active",
			ContactID:    request.SpaceID,
			Type:         briefs4contactus.ContactTypeCompany,
			SpaceRequest: request.SpaceRequest,
			Company: &dto4contactus.CreateCompanyRequest{
				Title:     workerParams.Space.Data.Title,
				VATNumber: request.VATNumber,
				Address:   &request.Address,
			},
		}
		for _, role := range request.Roles {
			createContactRequest.Roles = append(createContactRequest.Roles, string(role))
		}

		contactusWorkerParams := &dal4teamus.ModuleSpaceWorkerParams[*models4contactus.ContactusSpaceDbo]{
			SpaceWorkerParams: workerParams.SpaceWorkerParams,
			SpaceModuleEntry:  dal4contactus.NewContactusSpaceModuleEntry(request.SpaceID),
		}

		if teamContact, err = facade4contactus.CreateContactTx(ctx, tx, false, createContactRequest, contactusWorkerParams); err != nil {
			// Intentionally do not use original error to prevent wrongly returner HTTP status BadRequest=400
			return fmt.Errorf("failed to create team contact record: %v", err.Error())
		}
	} else if contactUpdates := updateContact(teamContact.Data, request); len(contactUpdates) > 0 {
		request := dto4contactus.UpdateContactRequest{
			ContactRequest: dto4contactus.ContactRequest{
				ContactID:    teamContact.ID,
				SpaceRequest: dto4teamus.SpaceRequest{SpaceID: request.SpaceID},
			},
			VatNumber: &request.VATNumber,
		}
		var contactWorkerParams *dal4contactus.ContactWorkerParams
		if err = facade4contactus.UpdateContactTx(ctx, tx, request, contactWorkerParams); err != nil {
			return fmt.Errorf("failed to update team contact record: %w", err)
		}
	}

	updates := updateLogistSpace(logistSpace.Data, workerParams.Space.Data, teamContact, request)

	if len(updates) > 0 {
		if err = logistSpace.Data.Validate(); err != nil {
			return fmt.Errorf("logistus team recod is not valid before saving: %w", err)
		}
		if logistSpace.Record.Exists() {
			if err = tx.Update(ctx, logistSpace.Key, updates); err != nil {
				return fmt.Errorf("failed to update logistus team record: %w", err)
			}
		} else if err = tx.Insert(ctx, logistSpace.Record); err != nil {
			return fmt.Errorf("failed to insert logistus team record: %w", err)
		}
	}
	return nil
}

func updateLogistSpace(logistSpaceDbo *dbo4logist.LogistSpaceDbo, spaceDbo *dbo4teamus.SpaceDbo, teamContact dal4contactus.ContactEntry, request dto4logist.SetLogistSpaceSettingsRequest) (updates []dal.Update) {
	if logistSpaceDbo.ContactID != teamContact.ID {
		logistSpaceDbo.ContactID = teamContact.ID
		updates = append(updates, dal.Update{Field: "contactID", Value: teamContact.ID})
	}
	if request.OrderNumberPrefix != "" {
		logistSpaceDbo.OrderNumberPrefix = request.OrderNumberPrefix
		updates = append(updates, dal.Update{Field: "orderNumberPrefix", Value: request.OrderNumberPrefix})
	}
	if dbo4logist.RolesChanged(logistSpaceDbo.Roles, request.Roles) {
		logistSpaceDbo.Roles = dbo4logist.ConvertLogistSpaceRolesToStringSlice(request.Roles)
		updates = append(updates, dal.Update{Field: "roles", Value: request.Roles})
	}
	if !slice.SameUniqueValues(logistSpaceDbo.UserIDs, spaceDbo.UserIDs) {
		logistSpaceDbo.UserIDs = spaceDbo.UserIDs
		updates = append(updates, dal.Update{Field: "userIDs", Value: spaceDbo.UserIDs})
	}
	return updates
}

func updateContact(contactDto *models4contactus.ContactDbo, request dto4logist.SetLogistSpaceSettingsRequest) (updates []dal.Update) {
	if contactDto.VATNumber != request.VATNumber {
		contactDto.VATNumber = request.VATNumber
		updates = append(updates, dal.Update{Field: "vatNumber", Value: request.VATNumber})
	}
	if dbo4logist.RolesChanged(contactDto.Roles, request.Roles) {
		contactDto.Roles = dbo4logist.ConvertLogistSpaceRolesToStringSlice(request.Roles)
		updates = append(updates, dal.Update{Field: "roles", Value: request.Roles})
	}
	return
}
