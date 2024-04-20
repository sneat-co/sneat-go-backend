package facade4linkage

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	"github.com/sneat-co/sneat-go-backend/src/modules/linkage/models4linkage"
	"github.com/strongo/slice"
)

type RelatableAdapter[D models4linkage.Relatable] interface {
	VerifyItem(ctx context.Context, tx dal.ReadTransaction, recordRef models4linkage.TeamModuleItemRef) (err error)
	//GetRecord(ctx context.Context, tx dal.ReadTransaction, recordRef models4linkage.TeamModuleItemRef) (record.DataWithID[string, D], error)
}
type relatableAdapter[D models4linkage.Relatable] struct {
	verifyItem func(ctx context.Context, tx dal.ReadTransaction, recordRef models4linkage.TeamModuleItemRef) (err error)
}

func (v relatableAdapter[D]) VerifyItem(ctx context.Context, tx dal.ReadTransaction, recordRef models4linkage.TeamModuleItemRef) (err error) {
	return v.verifyItem(ctx, tx, recordRef)
}

func NewRelatableAdapter[D models4linkage.Relatable](
	verifyItem func(ctx context.Context, tx dal.ReadTransaction, recordRef models4linkage.TeamModuleItemRef) (err error),
) RelatableAdapter[D] {
	return relatableAdapter[D]{
		verifyItem: verifyItem,
	}
}

//func (relatableAdapter[D]) GetRecord(ctx context.Context, tx dal.ReadTransaction, recordRef models4linkage.TeamModuleItemRef) (record.DataWithID[string, D], error) {
//	return nil, nil
//}

// SetRelated updates related records to define relationships
func SetRelated[D models4linkage.Relatable](
	_ context.Context,
	_ dal.ReadwriteTransaction,
	/*adapter*/ _ RelatableAdapter[D],
	object record.DataWithID[string, D],
	objectRef models4linkage.TeamModuleItemRef,
	relatedTo models4linkage.Link,
) (updates []dal.Update, err error) {

	{
		const invalidArgPrefix = "facade4linkage.SetRelated got invalid argument"
		if err = objectRef.Validate(); err != nil {
			return nil, fmt.Errorf("%s `objectRef models4linkage.TeamModuleItemRef`: %w", invalidArgPrefix, err)
		}
		if err = relatedTo.Validate(); err != nil {
			return nil, fmt.Errorf("%s 'relatedTo models4linkage.Link': %w", invalidArgPrefix, err)
		}
	}

	var updatedFields []string

	var relUpdates []dal.Update

	objectWithRelated := object.Data.GetRelated()
	if objectWithRelated.Related == nil {
		objectWithRelated.Related = make(models4linkage.RelatedByModuleID, 1)
	}
	getRelationships := func(ids []string) (relationships models4linkage.Relationships) {
		relationships = make(models4linkage.Relationships, len(ids))
		for _, r := range ids {
			relationships[r] = &models4linkage.Relationship{
				//CreatedField: with.CreatedField{
				//	Created: with.Created{
				//		At: now.Format(time.DateOnly),
				//		By: userID,
				//	},
				//},
			}
		}
		return
	}
	relatedAs := getRelationships(relatedTo.RelatedAs)
	relatesAs := getRelationships(relatedTo.RelatesAs)

	if relUpdates, err = objectWithRelated.AddRelationshipsAndIDs(
		objectRef,
		relatedTo.TeamModuleItemRef,
		relatedAs,
		relatesAs,
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
