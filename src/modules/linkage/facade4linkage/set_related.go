package facade4linkage

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	"github.com/sneat-co/sneat-go-backend/src/modules/linkage/models4linkage"
	"github.com/strongo/slice"
	"github.com/strongo/strongoapp/with"
	"time"
)

type RelatableAdapter[D models4linkage.Relatable] interface {
	VerifyItem(ctx context.Context, tx dal.ReadTransaction, recordRef models4linkage.TeamModuleDocRef) (err error)
	//GetRecord(ctx context.Context, tx dal.ReadTransaction, recordRef models4linkage.TeamModuleDocRef) (record.DataWithID[string, D], error)
}
type relatableAdapter[D models4linkage.Relatable] struct {
	verifyItem func(ctx context.Context, tx dal.ReadTransaction, recordRef models4linkage.TeamModuleDocRef) (err error)
}

func (v relatableAdapter[D]) VerifyItem(ctx context.Context, tx dal.ReadTransaction, recordRef models4linkage.TeamModuleDocRef) (err error) {
	return v.verifyItem(ctx, tx, recordRef)
}

func NewRelatableAdapter[D models4linkage.Relatable](
	verifyItem func(ctx context.Context, tx dal.ReadTransaction, recordRef models4linkage.TeamModuleDocRef) (err error),
) RelatableAdapter[D] {
	return relatableAdapter[D]{
		verifyItem: verifyItem,
	}
}

//func (relatableAdapter[D]) GetRecord(ctx context.Context, tx dal.ReadTransaction, recordRef models4linkage.TeamModuleDocRef) (record.DataWithID[string, D], error) {
//	return nil, nil
//}

// SetRelated updates related records to define relationships
func SetRelated[D models4linkage.Relatable](
	ctx context.Context,
	tx dal.ReadwriteTransaction,
	userID string,
	now time.Time,
	adapter RelatableAdapter[D],
	object record.DataWithID[string, D],
	objectRef models4linkage.TeamModuleDocRef,
	relatedTo models4linkage.Link,
) (updates []dal.Update, err error) {

	{
		const invalidArgPrefix = "facade4linkage.SetRelated got invalid argument"
		if err = objectRef.Validate(); err != nil {
			return nil, fmt.Errorf("%s `objectRef models4linkage.TeamModuleDocRef`: %w", invalidArgPrefix, err)
		}
		if err = relatedTo.Validate(); err != nil {
			return nil, fmt.Errorf("%s 'relatedTo models4linkage.Link': %w", invalidArgPrefix, err)
		}
	}

	var updatedFields []string

	var relUpdates []dal.Update

	objectWithRelated := object.Data.GetRelated()
	if objectWithRelated.Related == nil {
		objectWithRelated.Related = make(models4linkage.RelatedByTeamID, 1)
	}
	getRelationships := func(ids []string) (relationships models4linkage.Relationships) {
		relationships = make(models4linkage.Relationships, len(ids))
		for _, r := range ids {
			relationships[r] = &models4linkage.Relationship{
				CreatedField: with.CreatedField{
					Created: with.Created{
						At: now.Format(time.DateOnly),
						By: userID,
					},
				},
			}
		}
		return
	}
	relatedAs := getRelationships(relatedTo.RelatedAs)
	relatesAs := getRelationships(relatedTo.RelatesAs)

	if relUpdates, err = objectWithRelated.SetRelationshipsToItem(
		userID,
		objectRef,
		relatedTo.TeamModuleDocRef,
		relatedAs,
		relatesAs,
		now,
	); err != nil {
		return updates, err
	}
	updates = append(updates, relUpdates...)
	for _, update := range relUpdates {
		if !slice.Contains(updatedFields, update.Field) {
			updatedFields = append(updatedFields, update.Field)
		}
	}

	return updates, err
}
