package gaedal

import (
	"context"
	"errors"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/dtdal"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/models"
)

type splitDalGae struct {
}

var _ dtdal.SplitDal = (*splitDalGae)(nil) // Make sure we implement interface

func (splitDalGae) GetSplitByID(c context.Context, splitID int) (split models.Split, err error) {
	split.ID = splitID
	//err = dtdal.DB.Get(c, &split)
	return split, errors.New("TODO: implement")
}

func (splitDalGae) InsertSplit(c context.Context, splitEntity models.SplitEntity) (split models.Split, err error) {
	split.SplitEntity = &splitEntity
	//if err = dtdal.DB.Update(c, &split); err != nil {
	//	return
	//}
	return split, errors.New("TODO: implement")
}
