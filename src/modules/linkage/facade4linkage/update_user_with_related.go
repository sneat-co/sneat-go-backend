package facade4linkage

import (
	"context"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/linkage/dbo4linkage"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/dbo4userus"
	"slices"
)

func updateUserWithRelatedTx(
	ctx context.Context,
	tx dal.ReadwriteTransaction,
	userID string,
	users map[string]dbo4userus.User,
	itemRef dbo4linkage.SpaceModuleItemRef,
	relatedItem dbo4linkage.RelatedItem,
) (err error) {
	if users == nil {
		panic("users == nil")
	}

	var user dbo4userus.User
	var ok bool

	if user, ok = users[userID]; !ok {
		user := dbo4userus.NewUser(userID)
		if err = tx.Get(ctx, user.Record); err != nil {
			return err
		}
		users[userID] = user
	}

	if slices.Contains(user.Data.SpaceIDs, itemRef.Space) {
		return nil
	}

	userRelated := dbo4linkage.GetRelatedItemByRef(user.Data.Related, itemRef, true)

	var updates []dal.Update

	for roleID, role := range relatedItem.RolesToItem {
		if userRelated.RolesOfItem[roleID] != role {
			userRelated.RolesOfItem[roleID] = role
			updates = append(updates, dal.Update{Field: "related." + itemRef.ID() + ".rolesOfItem." + roleID, Value: role})
		}
	}

	return tx.Update(ctx, user.Record.Key(), updates)
}
