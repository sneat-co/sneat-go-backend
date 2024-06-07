package gaedal

import (
	"context"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/dtdal"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/models"
)

func NewRewardDalGae() rewardDalGae {
	return rewardDalGae{}
}

type rewardDalGae struct {
}

var _ dtdal.RewardDal = (*rewardDalGae)(nil)

func (rewardDalGae) InsertReward(c context.Context, tx dal.ReadwriteTransaction, rewardEntity *models.RewardDbo) (reward models.Reward, err error) {
	reward = models.NewRewardWithIncompleteKey(nil)
	if err = tx.Insert(c, reward.Record); err != nil {
		return
	}
	reward.ID = reward.Record.Key().ID.(string)
	return
}
