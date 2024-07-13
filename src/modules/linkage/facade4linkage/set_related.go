package facade4linkage

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/linkage/dbo4linkage"
)

//type RelatableAdapter[D dbo4linkage.Relatable] interface {
//	VerifyItem(ctx context.Context, tx dal.ReadTransaction, recordRef dbo4linkage.SpaceModuleItemRef) (err error)
//	//GetRecord(ctx context.Context, tx dal.ReadTransaction, recordRef dbo4linkage.SpaceModuleItemRef) (record.DataWithID[string, D], error)
//}
//type relatableAdapter[D dbo4linkage.Relatable] struct {
//	verifyItem func(ctx context.Context, tx dal.ReadTransaction, recordRef dbo4linkage.SpaceModuleItemRef) (err error)
//}
//
//func (v relatableAdapter[D]) VerifyItem(ctx context.Context, tx dal.ReadTransaction, recordRef dbo4linkage.SpaceModuleItemRef) (err error) {
//	return v.verifyItem(ctx, tx, recordRef)
//}
//
//func NewRelatableAdapter[D dbo4linkage.Relatable](
//	verifyItem func(ctx context.Context, tx dal.ReadTransaction, recordRef dbo4linkage.SpaceModuleItemRef) (err error),
//) RelatableAdapter[D] {
//	return relatableAdapter[D]{
//		verifyItem: verifyItem,
//	}
//}

//func (relatableAdapter[D]) GetRecord(ctx context.Context, tx dal.ReadTransaction, recordRef dbo4linkage.SpaceModuleItemRef) (record.DataWithID[string, D], error) {
//	return nil, nil
//}

type SetRelatedResult struct {
	ItemUpdates []dal.Update
}

// SetRelated updates related records to define relationships
func SetRelated(
	_ context.Context,
	_ dal.ReadwriteTransaction,
	object dbo4linkage.Relatable,
	objectRef dbo4linkage.SpaceModuleItemRef,
	itemRef dbo4linkage.SpaceModuleItemRef,
	rolesCommand dbo4linkage.RelationshipItemRolesCommand,
) (
	result SetRelatedResult,
	//teamModuleUpdates []dal.Update,
	err error,
) {

	{
		const invalidArgPrefix = "facade4linkage.SetRelated got invalid argument"
		if err = objectRef.Validate(); err != nil {
			err = fmt.Errorf("%s `objectRef dbo4linkage.SpaceModuleItemRef`: %w", invalidArgPrefix, err)
			return
		}
		if err = rolesCommand.Validate(); err != nil {
			return
		}
	}

	var relUpdates []dal.Update

	objectWithRelated := object.GetRelated()
	if objectWithRelated.Related == nil {
		objectWithRelated.Related = make(dbo4linkage.RelatedByModuleID, 1)
	}
	getRelationships := func(ids []string) (relationships dbo4linkage.RelationshipRoles) {
		relationships = make(dbo4linkage.RelationshipRoles, len(ids))
		for _, r := range ids {
			relationships[r] = &dbo4linkage.RelationshipRole{
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
	rolesOfItem := getRelationships(rolesCommand.Add.RolesOfItem)
	rolesToItem := getRelationships(rolesCommand.Add.RolesToItem)

	if relUpdates, err = objectWithRelated.AddRelationshipsAndIDs(
		itemRef,
		rolesOfItem,
		rolesToItem,
	); err != nil {
		return
	}
	result.ItemUpdates = append(result.ItemUpdates, relUpdates...)

	//for _, itemUpdate := range itemUpdates {
	//	if strings.HasSuffix(itemUpdate.Field, "relatedIDs") {
	//		continue // Ignore relatedIDs for teamModuleUpdates
	//	}
	//	teamModuleUpdates = append(teamModuleUpdates, dal.Update{
	//		Field: fmt.Sprintf("%s.%s.%s", objectRef.Collection, objectRef.ItemID, itemUpdate.Field),
	//		Value: itemUpdate.Value,
	//	})
	//}

	return
}
