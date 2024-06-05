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

// SetLogistTeamSettings sets team settings for logistus module
func SetLogistTeamSettings(
	ctx context.Context,
	userContext facade.User,
	request dto4logist.SetLogistTeamSettingsRequest,
) error {
	if err := request.Validate(); err != nil {
		return err
	}
	return dal4teamus.RunModuleTeamWorker(ctx, userContext, request.TeamRequest,
		const4logistus.ModuleID,
		new(dbo4logist.LogistTeamDbo),
		func(ctx context.Context, tx dal.ReadwriteTransaction, teamWorkerParams *dal4teamus.ModuleTeamWorkerParams[*dbo4logist.LogistTeamDbo]) (err error) {
			return setLogistTeamSettingsTx(ctx /*userContext,*/, request, tx, teamWorkerParams)
		},
	)
}

func setLogistTeamSettingsTx(
	ctx context.Context,
	//userContext facade.User,
	request dto4logist.SetLogistTeamSettingsRequest,
	tx dal.ReadwriteTransaction,
	workerParams *dal4teamus.ModuleTeamWorkerParams[*dbo4logist.LogistTeamDbo],
) (err error) {
	if workerParams.Team.Data.CountryID != request.Address.CountryID {
		workerParams.Team.Data.CountryID = request.Address.CountryID
		workerParams.TeamUpdates = append(workerParams.TeamUpdates, dal.Update{
			Field: "countryID",
			Value: request.Address.CountryID,
		})
	}

	logistTeam := dbo4logist.NewLogistTeamEntry(request.TeamID)
	if err = tx.Get(ctx, logistTeam.Record); err != nil {
		if !dal.IsNotFound(err) {
			return err
		}
	} else if err = logistTeam.Data.Validate(); err != nil {
		return fmt.Errorf("loaded logistus team recod is not valid: %w", err)
	}
	var teamContact dal4contactus.ContactEntry
	if teamContact, err = facade4contactus.GetContactByID(ctx, tx, logistTeam.ID, request.TeamID); err != nil {
		if !dal.IsNotFound(err) {
			return fmt.Errorf("failed to get contact record: %w", err)
		}
	}
	if dal.IsNotFound(err) {
		createContactRequest := dto4contactus.CreateContactRequest{
			Status:      "active",
			ContactID:   request.TeamID,
			Type:        briefs4contactus.ContactTypeCompany,
			TeamRequest: request.TeamRequest,
			Company: &dto4contactus.CreateCompanyRequest{
				Title:     workerParams.Team.Data.Title,
				VATNumber: request.VATNumber,
				Address:   &request.Address,
			},
		}
		for _, role := range request.Roles {
			createContactRequest.Roles = append(createContactRequest.Roles, string(role))
		}

		contactusWorkerParams := &dal4teamus.ModuleTeamWorkerParams[*models4contactus.ContactusTeamDbo]{
			TeamWorkerParams: workerParams.TeamWorkerParams,
			TeamModuleEntry:  dal4contactus.NewContactusTeamModuleEntry(request.TeamID),
		}

		if teamContact, err = facade4contactus.CreateContactTx(ctx, tx, false, createContactRequest, contactusWorkerParams); err != nil {
			// Intentionally do not use original error to prevent wrongly returner HTTP status BadRequest=400
			return fmt.Errorf("failed to create team contact record: %v", err.Error())
		}
	} else if contactUpdates := updateContact(teamContact.Data, request); len(contactUpdates) > 0 {
		request := dto4contactus.UpdateContactRequest{
			ContactRequest: dto4contactus.ContactRequest{
				ContactID:   teamContact.ID,
				TeamRequest: dto4teamus.TeamRequest{TeamID: request.TeamID},
			},
			VatNumber: &request.VATNumber,
		}
		var contactWorkerParams *dal4contactus.ContactWorkerParams
		if err = facade4contactus.UpdateContactTx(ctx, tx, request, contactWorkerParams); err != nil {
			return fmt.Errorf("failed to update team contact record: %w", err)
		}
	}

	updates := updateLogistTeam(logistTeam.Data, workerParams.Team.Data, teamContact, request)

	if len(updates) > 0 {
		if err = logistTeam.Data.Validate(); err != nil {
			return fmt.Errorf("logistus team recod is not valid before saving: %w", err)
		}
		if logistTeam.Record.Exists() {
			if err = tx.Update(ctx, logistTeam.Key, updates); err != nil {
				return fmt.Errorf("failed to update logistus team record: %w", err)
			}
		} else if err = tx.Insert(ctx, logistTeam.Record); err != nil {
			return fmt.Errorf("failed to insert logistus team record: %w", err)
		}
	}
	return nil
}

func updateLogistTeam(logistTeamDto *dbo4logist.LogistTeamDbo, teamDto *dbo4teamus.TeamDbo, teamContact dal4contactus.ContactEntry, request dto4logist.SetLogistTeamSettingsRequest) (updates []dal.Update) {
	if logistTeamDto.ContactID != teamContact.ID {
		logistTeamDto.ContactID = teamContact.ID
		updates = append(updates, dal.Update{Field: "contactID", Value: teamContact.ID})
	}
	if request.OrderNumberPrefix != "" {
		logistTeamDto.OrderNumberPrefix = request.OrderNumberPrefix
		updates = append(updates, dal.Update{Field: "orderNumberPrefix", Value: request.OrderNumberPrefix})
	}
	if dbo4logist.RolesChanged(logistTeamDto.Roles, request.Roles) {
		logistTeamDto.Roles = dbo4logist.ConvertLogistTeamRolesToStringSlice(request.Roles)
		updates = append(updates, dal.Update{Field: "roles", Value: request.Roles})
	}
	if !slice.SameUniqueValues(logistTeamDto.UserIDs, teamDto.UserIDs) {
		logistTeamDto.UserIDs = teamDto.UserIDs
		updates = append(updates, dal.Update{Field: "userIDs", Value: teamDto.UserIDs})
	}
	return updates
}

func updateContact(contactDto *models4contactus.ContactDbo, request dto4logist.SetLogistTeamSettingsRequest) (updates []dal.Update) {
	if contactDto.VATNumber != request.VATNumber {
		contactDto.VATNumber = request.VATNumber
		updates = append(updates, dal.Update{Field: "vatNumber", Value: request.VATNumber})
	}
	if dbo4logist.RolesChanged(contactDto.Roles, request.Roles) {
		contactDto.Roles = dbo4logist.ConvertLogistTeamRolesToStringSlice(request.Roles)
		updates = append(updates, dal.Update{Field: "roles", Value: request.Roles})
	}
	return
}
