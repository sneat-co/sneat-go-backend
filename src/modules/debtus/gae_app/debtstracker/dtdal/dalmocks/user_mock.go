package dalmocks

import (
	"context"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/models4debtus"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/dbo4userus"
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

func (mock *UserDalMock) GetUserByStrID(c context.Context, userID string) (user dbo4userus.UserEntry, err error) {
	panic("not implemented yet due to import cycle")
	// if user.ContactID, err = strconv.ParseInt(userID, 10, 64); err != nil {
	// 	return
	// }
	// return dal4userus.GetUserByID(c, user.ContactID)
}

func (mock *UserDalMock) CreateUser(c context.Context, userEntity *dbo4userus.UserDbo) (dbo4userus.UserEntry, error) {
	panic("Not implemented yet")
}

func (mock *UserDalMock) GetUserByVkUserID(c context.Context, vkUserID int64) (dbo4userus.UserEntry, error) {
	panic("Not implemented yet")
}
func (mock *UserDalMock) CreateAnonymousUser(c context.Context) (dbo4userus.UserEntry, error) {
	panic("Not implemented yet")
}
