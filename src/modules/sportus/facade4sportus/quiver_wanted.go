package facade4sportus

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/sportus/dbo4sportus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/validation"
)

// CreateWantedRequest defines request DTO
type CreateWantedRequest struct {
	dbo4sportus.Wanted
}

// Validate returns error if not valid
func (v *CreateWantedRequest) Validate() error {
	if err := v.Wanted.Validate(); err != nil {
		return err
	}
	return nil
}

func validateBrands(ctx context.Context, brands []string, db dal.DB) error {
	if len(brands) == 0 {
		return nil
	}
	brandRecords := make([]dal.Record, len(brands))
	for i, brand := range brands {
		key := dal.NewKeyWithID("Brand", brand)
		brandRecords[i] = dal.NewRecord(key)
	}
	if err := db.GetMulti(ctx, brandRecords); err != nil {
		return fmt.Errorf("failed to check brands: %w", err)
	}
	for _, brandRecord := range brandRecords {
		if !brandRecord.Exists() {
			return fmt.Errorf("unknown brand: %s", brandRecord.Key().ID)
		}
	}
	return nil
}

// CreateWanted creates wanted records
func CreateWanted(ctx facade.ContextWithUser, request CreateWantedRequest) (id string, err error) {
	userCtx := ctx.User()
	var db dal.DB
	if db, err = facade.GetSneatDB(ctx); err != nil {
		return "", err
	}
	if err = validateBrands(ctx, request.Brands, db); err != nil {
		return "", err
	}
	err = db.RunReadwriteTransaction(ctx, func(ctx context.Context, tx dal.ReadwriteTransaction) error {
		record := dal.NewRecordWithIncompleteKey(dbo4sportus.QuiverWantedCollection, reflect.String, &request.Wanted)
		request.UserID = userCtx.GetUserID()
		if err := tx.Insert(ctx, record); err != nil {
			return fmt.Errorf("failed to create wanted record: %w", err)
		}
		id = fmt.Sprintf("%v", record.Key().ID)
		return nil
	})
	return id, err
}

// DeleteWantedRequest defines delete w
type DeleteWantedRequest struct {
	ID string
}

// Validate returns error if not valid
func (v *DeleteWantedRequest) Validate() error {
	if strings.TrimSpace(v.ID) == "" {
		return validation.NewErrRequestIsMissingRequiredField("id")
	}
	return nil
}

// DeleteWanted deletes wanted records
func DeleteWanted(ctx facade.ContextWithUser, request DeleteWantedRequest) error {
	db, err := facade.GetSneatDB(ctx)
	if err != nil {
		return err
	}
	userCtx := ctx.User()
	return db.RunReadwriteTransaction(ctx, func(ctx context.Context, tx dal.ReadwriteTransaction) error {
		key := dal.NewKeyWithID(dbo4sportus.QuiverWantedCollection, request.ID)
		var wanted dbo4sportus.Wanted
		record := dal.NewRecordWithData(key, &wanted)
		if err := tx.Get(ctx, record); err != nil {
			return err
		}
		uid := userCtx.GetUserID()
		if wanted.UserID != uid {
			return fmt.Errorf("wanted.UserID != userCtx.ContactID(): %s != %s", wanted.UserID, uid)
		}
		if err := tx.Delete(ctx, key); err != nil {
			return fmt.Errorf("failed to delete wanted record: %v", err)
		}
		return nil
	})
}
