package facade4retrospectus

import (
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-core/coretypes"
)

func getSpaceRetroDocKey(spaceID coretypes.SpaceID, meeting string) *dal.Key {
	parentKey := dal.NewKeyWithID("api4meetingus", spaceID)
	return dal.NewKeyWithParentAndID(parentKey, "api4meetingus", meeting)
}
