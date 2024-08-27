package dalmocks

import (
	"context"
	"github.com/crediterra/money"
	"github.com/dal-go/mocks4dalgo/mocks4dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/gae_app/debtstracker/dtdal"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/models4debtus"
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

func (mock *TransferDalMock) DelayUpdateTransfersOnReturn(_ context.Context, returntransferID int, transferReturnUpdates []dtdal.TransferReturnUpdate) (err error) {
	panic("not implemented yet")
}

func (mock *TransferDalMock) GetTransfersByID(_ context.Context, transferIDs []int) ([]models4debtus.TransferEntry, error) {
	panic("not implemented yet")
}

func (mock *TransferDalMock) LoadTransfersByUserID(_ context.Context, userID string, offset, limit int) (transfers []models4debtus.TransferEntry, hasMore bool, err error) {
	panic(NOT_IMPLEMENTED_YET)
}

func (mock *TransferDalMock) LoadTransferIDsByContactID(_ context.Context, contactID string, limit int, startCursor string) (transferIDs []int, endCursor string, err error) {
	panic(NOT_IMPLEMENTED_YET)
}

func (mock *TransferDalMock) LoadTransfersByContactID(_ context.Context, contactID string, offset, limit int) (transfers []models4debtus.TransferEntry, hasMore bool, err error) {
	panic(NOT_IMPLEMENTED_YET)
}

func (mock *TransferDalMock) LoadOverdueTransfers(_ context.Context, userID string, limit int) (transfers []models4debtus.TransferEntry, err error) {
	panic(NOT_IMPLEMENTED_YET)
}

func (mock *TransferDalMock) LoadOutstandingTransfers(_ context.Context, periodEnds time.Time, userID, contactID string, currency money.CurrencyCode, direction models4debtus.TransferDirection) (transfers []models4debtus.TransferEntry, err error) {
	panic("not implemented yet")
	//for _, entity := range mock.mockDB.EntitiesByKind[models.TransfersCollection] {
	//	t := entity.(*models.TransferEntry)
	//	if t.Direction() == direction && t.GetOutstandingValue(periodEnds) != 0 {
	//		api4transfers = append(api4transfers, *t)
	//	}
	//}
	//return
}

func (mock *TransferDalMock) LoadDueTransfers(_ context.Context, userID string, limit int) (transfers []models4debtus.TransferEntry, err error) {
	panic(NOT_IMPLEMENTED_YET)
}

func (mock *TransferDalMock) LoadLatestTransfers(_ context.Context, offset, limit int) ([]models4debtus.TransferEntry, error) {
	panic(NOT_IMPLEMENTED_YET)
}

func (mock *TransferDalMock) DelayUpdateTransferWithCreatorReceiptTgMessageID(_ context.Context, botCode string, transferID int, creatorTgChatID, creatorTgReceiptMessageID int64) error {
	panic(NOT_IMPLEMENTED_YET)
}

func (mock *TransferDalMock) DelayUpdateTransfersWithCounterparty(_ context.Context, creatorCounterpartyID, counterpartyCounterpartyID string) error {
	panic(NOT_IMPLEMENTED_YET)
}
