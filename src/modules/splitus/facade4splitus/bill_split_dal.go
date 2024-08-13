package facade4splitus

import (
	"context"
	"errors"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/gae_app/debtstracker/dtdal"
	"github.com/sneat-co/sneat-go-backend/src/modules/splitus/models4splitus"
)

type splitDalGae struct {
}

var _ dtdal.SplitDal = (*splitDalGae)(nil) // Make sure we implement interface

func (splitDalGae) GetSplitByID(c context.Context, splitID int) (split models4splitus.Split, err error) {
	split.ID = splitID
	//err = dtdal.DB.Get(c, &split)
	return split, errors.New("TODO: implement")
}

func (splitDalGae) InsertSplit(c context.Context, splitEntity models4splitus.SplitEntity) (split models4splitus.Split, err error) {
	split.SplitEntity = &splitEntity
	//if err = dtdal.DB.Update(c, &split); err != nil {
	//	return
	//}
	return split, errors.New("TODO: implement")
}
