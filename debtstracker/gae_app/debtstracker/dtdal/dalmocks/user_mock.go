package dalmocks

import (
	"github.com/crediterra/money"
	"time"

	"context"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/models"
)

type UserDalMock struct {
	LastUserID int64
	Users      map[int64]*models.DebutsAppUserDataOBSOLETE
}

func (mock *UserDalMock) DelayUpdateUserWithContact(c context.Context, userID, contactID string) error {
	panic("not implemented yet")
}

func NewUserDalMock() *UserDalMock {
	return &UserDalMock{
		Users: make(map[int64]*models.DebutsAppUserDataOBSOLETE),
	}
}

func (mock *UserDalMock) SetLastCurrency(c context.Context, userID int64, currency money.CurrencyCode) error {
	panic("Not implemented yet")
}

func (mock *UserDalMock) GetUserByStrID(c context.Context, userID string) (user models.AppUser, err error) {
	panic("not implemented yet due to import cycle")
	// if user.ID, err = strconv.ParseInt(userID, 10, 64); err != nil {
	// 	return
	// }
	// return facade.User.GetUserByID(c, user.ID)
}

func (mock *UserDalMock) GetUserByEmail(c context.Context, email string) (models.AppUser, error) {
	panic("Not implemented yet")
}

func (mock *UserDalMock) CreateUser(c context.Context, userEntity *models.DebutsAppUserDataOBSOLETE) (models.AppUser, error) {
	panic("Not implemented yet")
}

func (mock *UserDalMock) GetUserByVkUserID(c context.Context, vkUserID int64) (models.AppUser, error) {
	panic("Not implemented yet")
}
func (mock *UserDalMock) CreateAnonymousUser(c context.Context) (models.AppUser, error) {
	panic("Not implemented yet")
}
func (mock *UserDalMock) DelaySetUserPreferredLocale(c context.Context, delay time.Duration, userID int64, localeCode5 string) error {
	return nil
}
func (mock *UserDalMock) DelayUpdateUserHasDueTransfers(c context.Context, userID int64) error {
	return nil
}
func (mock *UserDalMock) DelayUpdateUserWithBill(c context.Context, groupID, billID string) error {
	panic(NOT_IMPLEMENTED_YET)
}
