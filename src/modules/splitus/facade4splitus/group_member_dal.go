package facade4splitus

import (
	"context"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/splitus/models4splitus"
	"github.com/sneat-co/sneat-go-core/facade"
)

type GroupMemberDalGae struct {
}

func NewGroupMemberDalGae() GroupMemberDalGae {
	return GroupMemberDalGae{}
}

func (GroupMemberDalGae) CreateGroupMember(ctx context.Context, tx dal.ReadwriteTransaction, groupMemberData *models4splitus.GroupMemberData) (groupMember models4splitus.GroupMember, err error) {
	key := models4splitus.NewGroupMemberIncompleteKey()
	groupMember.Record = dal.NewRecordWithData(key, groupMemberData)
	if err = tx.Insert(ctx, groupMember.Record); err != nil {
		return
	}
	groupMember.ID = groupMember.Record.Key().ID.(int64)
	return
}

func (GroupMemberDalGae) GetGroupMemberByID(ctx context.Context, tx dal.ReadSession, groupMemberID int64) (groupMember models4splitus.GroupMember, err error) {
	groupMember = models4splitus.NewGroupMember(groupMemberID, nil)
	if tx == nil {
		if tx, err = facade.GetDatabase(ctx); err != nil {
			return
		}
	}
	return groupMember, tx.Get(ctx, groupMember.Record)
}
