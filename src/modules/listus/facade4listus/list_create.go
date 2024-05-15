package facade4listus

import (
	"context"
	"errors"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/listus/const4listus"
	"github.com/sneat-co/sneat-go-backend/src/modules/listus/dal4listus"
	"github.com/sneat-co/sneat-go-backend/src/modules/listus/models4listus"
	"github.com/sneat-co/sneat-go-backend/src/modules/teamus/dal4teamus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/strongo/random"
	"github.com/strongo/strongoapp/with"
	"strings"
	"time"
)

// CreateList creates a new list
func CreateList(ctx context.Context, user facade.User, request CreateListRequest) (response CreateListResponse, err error) {
	request.Title = strings.TrimSpace(request.Title)
	if err = request.Validate(); err != nil {
		return
	}
	err = dal4teamus.CreateTeamItem(ctx, user, "", request.TeamRequest, const4listus.ModuleID, new(models4listus.ListusTeamDto),
		func(ctx context.Context, tx dal.ReadwriteTransaction, params *dal4teamus.ModuleTeamWorkerParams[*models4listus.ListusTeamDto]) (err error) {

			for id, brief := range params.TeamModuleEntry.Data.Lists {
				if brief.Title == request.Title {
					return fmt.Errorf("an attempt to create a new list with duplicate title [%s] equal to list {id=%s, type=%s}", request.Title, id, brief.Type)
				}
			}
			id := request.Type
			idTryNumber := 0
			idLen := 2               // must be before checkId label
			idGenerationAttempt := 0 // must be before checkId label
		checkId:
			if idTryNumber++; idTryNumber == 10 {
				return errors.New("too many attempts to generate random list ContactID")
			}
			for briefID := range params.TeamModuleEntry.Data.Lists {
				if briefID == id {
					id = random.ID(idLen)
					idGenerationAttempt++
					if idGenerationAttempt > 1000 {
						idLen++
						idGenerationAttempt = 0
					}
					goto checkId
				}
			}
			modified := dbmodels.Modified{
				By: user.GetID(),
				At: time.Now(),
			}
			list := models4listus.ListDto{
				WithModified: dbmodels.WithModified{
					CreatedFields: with.CreatedFields{
						CreatedAtField: with.CreatedAtField{
							CreatedAt: modified.At,
						},
						CreatedByField: with.CreatedByField{
							CreatedBy: modified.By,
						},
					},
					UpdatedFields: with.UpdatedFields{
						UpdatedAt: modified.At,
						UpdatedBy: modified.By,
					},
				},
				WithTeamIDs: dbmodels.WithTeamIDs{
					TeamIDs: []string{request.TeamID},
				},
				ListBase: models4listus.ListBase{
					Type:  request.Type,
					Title: request.Title,
				},
			}
			if err = list.Validate(); err != nil {
				return fmt.Errorf("formed list DTO struct is not valid: %w", err)
			}
			listKey := dal4listus.NewTeamListKey(request.TeamID, id)
			listRecord := dal.NewRecordWithData(listKey, &list)
			if err = tx.Insert(ctx, listRecord); err != nil {
				return fmt.Errorf("failed to insert list record")
			}
			if params.TeamModuleEntry.Data.Lists == nil {
				params.TeamModuleEntry.Data.Lists = make(map[string]*models4listus.ListBrief, 1)
			}
			listBrief := &models4listus.ListBrief{
				ListBase: models4listus.ListBase{
					Type:  request.Type,
					Title: request.Type,
				},
			}
			params.TeamModuleEntry.Data.Lists[id] = listBrief
			if params.TeamModuleEntry.Record.Exists() {
				if err = tx.Insert(ctx, params.TeamModuleEntry.Record); err != nil {
					return fmt.Errorf("failed to insert team module entry record: %w", err)
				}
			} else {
				params.TeamUpdates = append(params.TeamUpdates, dal.Update{
					Field: "lists." + id,
					Value: listBrief,
				})
			}
			return err
		},
	)
	return
}
