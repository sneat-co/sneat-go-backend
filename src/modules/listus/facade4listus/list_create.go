package facade4listus

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/update"
	dal4spaceus2 "github.com/sneat-co/sneat-core-modules/spaceus/dal4spaceus"
	"github.com/sneat-co/sneat-go-backend/src/modules/listus/const4listus"
	"github.com/sneat-co/sneat-go-backend/src/modules/listus/dal4listus"
	"github.com/sneat-co/sneat-go-backend/src/modules/listus/dbo4listus"
	"github.com/sneat-co/sneat-go-backend/src/modules/listus/dto4listus"
	"github.com/sneat-co/sneat-go-core/coretypes"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/strongo/random"
	"github.com/strongo/strongoapp/with"
	"strings"
	"time"
)

// CreateList creates a new list
func CreateList(ctx context.Context, userCtx facade.UserContext, request dto4listus.CreateListRequest) (response dto4listus.CreateListResponse, err error) {
	request.Title = strings.TrimSpace(request.Title)
	if err = request.Validate(); err != nil {
		return
	}
	err = dal4spaceus2.CreateSpaceItem(ctx, userCtx, request.SpaceRequest, const4listus.ModuleID, new(dbo4listus.ListusSpaceDbo),
		func(ctx context.Context, tx dal.ReadwriteTransaction, params *dal4spaceus2.ModuleSpaceWorkerParams[*dbo4listus.ListusSpaceDbo]) (err error) {

			for id, brief := range params.SpaceModuleEntry.Data.Lists {
				if brief.Title == request.Title {
					return fmt.Errorf("an attempt to create a new listDbo with duplicate title [%s] equal to listDbo {listType=%s, type=%s}", request.Title, id, brief.Type)
				}
			}
			listType := request.Type
			idLen := 2               // must be before checkId label
			idGenerationAttempt := 0 // must be before checkId label
			var listSubID string
			for {
				listSubID = random.ID(idLen)
				listID := dbo4listus.NewListKey(listType, listSubID)
				if _, found := params.SpaceModuleEntry.Data.Lists[string(listID)]; !found {
					break
				}
				idGenerationAttempt++
				if idGenerationAttempt > 1000 {
					err = fmt.Errorf("too many attempts to generate random listDbo ID")
					return
				}
			}

			listID := dbo4listus.NewListKey(listType, listSubID)

			modified := dbmodels.Modified{
				By: userCtx.GetUserID(),
				At: time.Now(),
			}
			listDbo := dbo4listus.ListDbo{
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
					SpaceIDs: []coretypes.SpaceID{request.SpaceID},
				},
				ListBase: dbo4listus.ListBase{
					Type:  request.Type,
					Title: request.Title,
				},
			}
			if err = listDbo.Validate(); err != nil {
				return fmt.Errorf("formed listDbo DTO struct is not valid: %w", err)
			}
			listKey := dal4listus.NewListKey(request.SpaceID, listID)
			listRecord := dal.NewRecordWithData(listKey, &listDbo)
			if err = tx.Insert(ctx, listRecord); err != nil {
				return fmt.Errorf("failed to insert listDbo record")
			}
			if params.SpaceModuleEntry.Data.Lists == nil {
				params.SpaceModuleEntry.Data.Lists = make(dbo4listus.ListBriefs, 1)
			}
			listBrief := &dbo4listus.ListBrief{
				ListBase: dbo4listus.ListBase{
					Type:  request.Type,
					Title: request.Type,
				},
			}
			params.SpaceModuleEntry.Data.Lists[string(listID)] = listBrief
			if params.SpaceModuleEntry.Record.Exists() {
				if err = tx.Insert(ctx, params.SpaceModuleEntry.Record); err != nil {
					return fmt.Errorf("failed to insert team module entry record: %w", err)
				}
			} else {
				params.SpaceUpdates = append(params.SpaceUpdates, update.ByFieldName("lists."+listType, listBrief))
			}
			return err
		},
	)
	return
}
