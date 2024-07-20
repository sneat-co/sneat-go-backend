package facade4listus

import (
	"context"
	"errors"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/listus/const4listus"
	"github.com/sneat-co/sneat-go-backend/src/modules/listus/dal4listus"
	"github.com/sneat-co/sneat-go-backend/src/modules/listus/dbo4listus"
	"github.com/sneat-co/sneat-go-backend/src/modules/spaceus/dal4spaceus"
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
	err = dal4spaceus.CreateSpaceItem(ctx, user, request.SpaceRequest, const4listus.ModuleID, new(dbo4listus.ListusSpaceDbo),
		func(ctx context.Context, tx dal.ReadwriteTransaction, params *dal4spaceus.ModuleSpaceWorkerParams[*dbo4listus.ListusSpaceDbo]) (err error) {

			for id, brief := range params.SpaceModuleEntry.Data.Lists {
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
			for briefID := range params.SpaceModuleEntry.Data.Lists {
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
			list := dbo4listus.ListDbo{
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
				WithSpaceIDs: dbmodels.WithSpaceIDs{
					SpaceIDs: []string{request.SpaceID},
				},
				ListBase: dbo4listus.ListBase{
					Type:  request.Type,
					Title: request.Title,
				},
			}
			if err = list.Validate(); err != nil {
				return fmt.Errorf("formed list DTO struct is not valid: %w", err)
			}
			listKey := dal4listus.NewSpaceListKey(request.SpaceID, id)
			listRecord := dal.NewRecordWithData(listKey, &list)
			if err = tx.Insert(ctx, listRecord); err != nil {
				return fmt.Errorf("failed to insert list record")
			}
			if params.SpaceModuleEntry.Data.Lists == nil {
				params.SpaceModuleEntry.Data.Lists = make(map[string]*dbo4listus.ListBrief, 1)
			}
			listBrief := &dbo4listus.ListBrief{
				ListBase: dbo4listus.ListBase{
					Type:  request.Type,
					Title: request.Type,
				},
			}
			params.SpaceModuleEntry.Data.Lists[id] = listBrief
			if params.SpaceModuleEntry.Record.Exists() {
				if err = tx.Insert(ctx, params.SpaceModuleEntry.Record); err != nil {
					return fmt.Errorf("failed to insert team module entry record: %w", err)
				}
			} else {
				params.SpaceUpdates = append(params.SpaceUpdates, dal.Update{
					Field: "lists." + id,
					Value: listBrief,
				})
			}
			return err
		},
	)
	return
}
