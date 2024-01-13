package facade4retrospectus

import (
	"github.com/dal-go/dalgo/dal"
)

func getTeamRetroDocKey(team, meeting string) *dal.Key {
	parentKey := dal.NewKeyWithID("api4meetingus", team)
	return dal.NewKeyWithParentAndID(parentKey, "api4meetingus", meeting)
}
