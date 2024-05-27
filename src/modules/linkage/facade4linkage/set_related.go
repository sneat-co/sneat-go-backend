package facade4linkage

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	"github.com/sneat-co/sneat-go-backend/src/modules/linkage/models4linkage"
	"strings"
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
	itemRef models4linkage.TeamModuleItemRef,
	rolesCommand models4linkage.RelationshipRolesCommand,
) (
	itemUpdates []dal.Update,
	teamModuleUpdates []dal.Update,
	err error,
) {

	{
		const invalidArgPrefix = "facade4linkage.SetRelated got invalid argument"
		if err = objectRef.Validate(); err != nil {
			return nil, nil, fmt.Errorf("%s `objectRef models4linkage.TeamModuleItemRef`: %w", invalidArgPrefix, err)
		}
		if err = rolesCommand.Validate(); err != nil {
			return nil, nil, err
		}
	}

	var relUpdates []dal.Update

	objectWithRelated := object.Data.GetRelated()
	if objectWithRelated.Related == nil {
		objectWithRelated.Related = make(models4linkage.RelatedByModuleID, 1)
	}
	getRelationships := func(ids []string) (relationships models4linkage.RelationshipRoles) {
		relationships = make(models4linkage.RelationshipRoles, len(ids))
		for _, r := range ids {
			relationships[r] = &models4linkage.RelationshipRole{
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
		return itemUpdates, teamModuleUpdates, err
	}
	itemUpdates = append(itemUpdates, relUpdates...)

	for _, itemUpdate := range itemUpdates {
		if strings.HasSuffix(itemUpdate.Field, "relatedIDs") {
			continue // Ignore relatedIDs for teamModuleUpdates
		}
		teamModuleUpdates = append(teamModuleUpdates, dal.Update{
			Field: fmt.Sprintf("%s.%s.%s", objectRef.Collection, objectRef.ItemID, itemUpdate.Field),
			Value: itemUpdate.Value,
		})
	}

	return itemUpdates, teamModuleUpdates, err
}
