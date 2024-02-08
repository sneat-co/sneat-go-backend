package models

import (
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	"reflect"
	"time"
)

const RewardKind = "Reward"

type Reward struct {
	record.WithID[int]
	Data *RewardData
}

func NewReward(id int, data *RewardData) Reward {
	key := dal.NewKeyWithID(RewardKind, id)
	if data == nil {
		data = new(RewardData)
	}
	return Reward{
		WithID: record.NewWithID(id, key, data),
		Data:   data,
	}
}

func NewRewardWithIncompleteKey(data *RewardData) Reward {
	key := dal.NewIncompleteKey(RewardKind, reflect.Int, nil)
	return Reward{
		WithID: record.NewWithID(0, key, data),
		Data:   data,
	}
}

type RewardReason string

const (
	RewardReasonInvitedUserJoined         RewardReason = "InvitedUserJoined"
	RewardReasonFriendOfInvitedUserJoined RewardReason = "FriendOfInvitedUserJoined"
)

type RewardData struct {
	UserID       int64
	DtCreated    time.Time
	Reason       RewardReason `datastore:",noindex"`
	JoinedUserID int64        `datastore:",noindex"`
	Points       int          `datastore:",noindex"`
}

type UserRewardBalance struct {
	RewardPoints   int
	RewardOptedOut time.Time
	RewardIDs      []int64 `datastore:",noindex"`
}

//func (UserRewardBalance) cleanProperties(properties []datastore.Property) ([]datastore.Property, error) {
//	return gaedb.CleanProperties(properties, map[string]gaedb.IsOkToRemove{
//		"RewardPoints":   gaedb.IsZeroInt,
//		"RewardOptedOut": gaedb.IsZeroTime,
//	})
//}

func (rewardBalance *UserRewardBalance) AddRewardPoints(rewardID int64, rewardPoints int) (changed bool) {
	for _, id := range rewardBalance.RewardIDs {
		if id == rewardID {
			return
		}
	}
	rewardBalance.RewardPoints += rewardPoints
	rewardBalance.RewardIDs = append([]int64{rewardID}, rewardBalance.RewardIDs...)
	return true
}
