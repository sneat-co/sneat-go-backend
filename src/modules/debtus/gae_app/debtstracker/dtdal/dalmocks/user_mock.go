package dalmocks

import (
	"context"
	dbo4userus2 "github.com/sneat-co/sneat-core-modules/userus/dbo4userus"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/models4debtus"
)

type UserDalMock struct {
	LastUserID int64
	Users      map[int64]*models4debtus.DebutsAppUserDataOBSOLETE
}

func NewUserDalMock() *UserDalMock {
	return &UserDalMock{
		Users: make(map[int64]*models4debtus.DebutsAppUserDataOBSOLETE),
	}
}

func (mock *UserDalMock) GetUserByStrID(_ context.Context, userID string) (user dbo4userus2.UserEntry, err error) {
	panic("not implemented yet due to import cycle")
	// if user.ContactID, err = strconv.ParseInt(userID, 10, 64); err != nil {
	// 	return
	// }
	// return dal4userus.GetUserByID(ctx, user.ContactID)
}

func (mock *UserDalMock) CreateUser(_ context.Context, userEntity *dbo4userus2.UserDbo) (dbo4userus2.UserEntry, error) {
	panic("Not implemented yet")
}

func (mock *UserDalMock) GetUserByVkUserID(_ context.Context, vkUserID int64) (dbo4userus2.UserEntry, error) {
	panic("Not implemented yet")
}
func (mock *UserDalMock) CreateAnonymousUser(_ context.Context) (dbo4userus2.UserEntry, error) {
	panic("Not implemented yet")
}
