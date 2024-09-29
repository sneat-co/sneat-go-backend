package unsorted4auth

import "time"

type UserRewardBalance struct {
	RewardPoints   int
	RewardOptedOut time.Time
	RewardIDs      []int64 `firestore:",omitempty"`
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
