package facade4linkage

import (
	"context"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/linkage/models4linkage"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/models4userus"
	"github.com/strongo/slice"
)

func updateUserWithRelatedTx(
	ctx context.Context,
	tx dal.ReadwriteTransaction,
	userID string,
	users map[string]models4userus.User,
	itemRef models4linkage.TeamModuleItemRef,
	relatedItem models4linkage.RelatedItem,
) (err error) {
	if users == nil {
		panic("users == nil")
	}

	var user models4userus.User
	var ok bool

	if user, ok = users[userID]; !ok {
		user := models4userus.NewUser(userID)
		if err = tx.Get(ctx, user.Record); err != nil {
			return err
		}
		users[userID] = user
	}

	if slice.Contains(user.Data.TeamIDs, itemRef.TeamID) {
		return nil
	}

	userRelated := models4linkage.GetRelatedItemByRef(user.Data.Related, itemRef, true)

	var updates []dal.Update

	for roleID, role := range relatedItem.RolesToItem {
		if userRelated.RolesOfItem[roleID] != role {
			userRelated.RolesOfItem[roleID] = role
			updates = append(updates, dal.Update{Field: "related." + itemRef.ID() + ".rolesOfItem." + roleID, Value: role})
		}
	}

	return tx.Update(ctx, user.Record.Key(), updates)
}
