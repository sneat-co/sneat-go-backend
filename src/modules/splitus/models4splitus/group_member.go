package models4splitus

import (
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
)

const GroupMemberKind = "GroupMember"

type GroupMember struct {
	record.WithID[int64]
	Data *GroupMemberData
}

func NewGroupMember(id int64, data *GroupMemberData) GroupMember {
	key := NewGroupMemberKey(id)
	if data == nil {
		data = new(GroupMemberData)
	}
	return GroupMember{
		WithID: record.NewWithID(id, key, data),
		Data:   data,
	}
}

func NewGroupMemberKey(groupMemberID int64) *dal.Key {
	if groupMemberID == 0 {
		panic("groupMemberID == 0")
	}
	return dal.NewKeyWithID(GroupMemberKind, groupMemberID)
}

func NewGroupMemberIncompleteKey() *dal.Key {
	return dal.NewKeyWithID(GroupMemberKind, 0)
}

//var _ db.EntityHolder = (*GroupMember)(nil)

//func (GroupMember) Kind() string {
//	return GroupMemberKind
//}

//func (gm GroupMember) Entity() interface{} {
//	return gm.Data
//}
//
//func (GroupMember) NewEntity() interface{} {
//	return new(GroupMemberData)
//}
//
//func (gm *GroupMember) SetEntity(entity interface{}) {
//	if entity == nil {
//		gm.GroupMemberData = nil
//	} else {
//		gm.GroupMemberData = entity.(*GroupMemberData)
//	}
//}

type GroupMemberData struct {
	GroupID    int64
	UserID     int64
	ContactIDs []int64
	Name       string `firestore:",omitempty"`
}
