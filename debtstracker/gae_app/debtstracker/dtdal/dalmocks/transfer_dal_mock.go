package dalmocks

import (
	"context"
	"github.com/crediterra/money"
	"github.com/dal-go/mocks4dalgo/mocks4dal"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/dtdal"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/models"
	"time"
)

const NOT_IMPLEMENTED_YET = "NOT_IMPLEMENTED_YET"

type TransferDalMock struct {
	mockDB *mocks4dal.MockDatabase
}

func NewTransferDalMock(mockDB *mocks4dal.MockDatabase) *TransferDalMock {
	return &TransferDalMock{
		mockDB: mockDB,
	}
}

func (mock *TransferDalMock) DelayUpdateTransfersOnReturn(c context.Context, returntransferID int, transferReturnUpdates []dtdal.TransferReturnUpdate) (err error) {
	panic("not implemented yet")
}

func (mock *TransferDalMock) GetTransfersByID(c context.Context, transferIDs []int) ([]models.Transfer, error) {
	panic("not implemented yet")
}

func (mock *TransferDalMock) LoadTransfersByUserID(c context.Context, userID string, offset, limit int) (transfers []models.Transfer, hasMore bool, err error) {
	panic(NOT_IMPLEMENTED_YET)
}

func (mock *TransferDalMock) LoadTransferIDsByContactID(c context.Context, contactID string, limit int, startCursor string) (transferIDs []int, endCursor string, err error) {
	panic(NOT_IMPLEMENTED_YET)
}

func (mock *TransferDalMock) LoadTransfersByContactID(c context.Context, contactID string, offset, limit int) (transfers []models.Transfer, hasMore bool, err error) {
	panic(NOT_IMPLEMENTED_YET)
}

func (mock *TransferDalMock) LoadOverdueTransfers(c context.Context, userID string, limit int) (transfers []models.Transfer, err error) {
	panic(NOT_IMPLEMENTED_YET)
}

func (mock *TransferDalMock) LoadOutstandingTransfers(c context.Context, periodEnds time.Time, userID, contactID string, currency money.CurrencyCode, direction models.TransferDirection) (transfers []models.Transfer, err error) {
	panic("not implemented yet")
	//for _, entity := range mock.mockDB.EntitiesByKind[models.TransferKind] {
	//	t := entity.(*models.Transfer)
	//	if t.Direction() == direction && t.GetOutstandingValue(periodEnds) != 0 {
	//		transfers = append(transfers, *t)
	//	}
	//}
	//return
}

func (mock *TransferDalMock) LoadDueTransfers(c context.Context, userID string, limit int) (transfers []models.Transfer, err error) {
	panic(NOT_IMPLEMENTED_YET)
}

func (mock *TransferDalMock) LoadLatestTransfers(c context.Context, offset, limit int) ([]models.Transfer, error) {
	panic(NOT_IMPLEMENTED_YET)
}

func (mock *TransferDalMock) DelayUpdateTransferWithCreatorReceiptTgMessageID(c context.Context, botCode string, transferID int, creatorTgChatID, creatorTgReceiptMessageID int64) error {
	panic(NOT_IMPLEMENTED_YET)
}

func (mock *TransferDalMock) DelayUpdateTransfersWithCounterparty(c context.Context, creatorCounterpartyID, counterpartyCounterpartyID string) error {
	panic(NOT_IMPLEMENTED_YET)
}
