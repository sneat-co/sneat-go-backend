package mocks4logist

import (
	"context"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/mocks4dalgo/mocks4dal"
	"github.com/golang/mock/gomock"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/dbo4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/linkage/dbo4linkage"
	"github.com/sneat-co/sneat-go-backend/src/modules/logistus/dbo4logist"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"testing"
	"time"
)

func MockTx(t *testing.T) (tx *mocks4dal.MockReadwriteTransaction) {
	mockCtrl := gomock.NewController(t)
	tx = mocks4dal.NewMockReadwriteTransaction(mockCtrl)
	tx.EXPECT().
		Get(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, record dal.Record) error {
			record.SetError(nil)
			switch record.Key().Collection() {
			case dbo4logist.OrdersCollection:
				orderDto := record.Data().(*dbo4logist.OrderDbo)
				orderDto.Status = "active"
			case models4contactus.TeamContactsCollection:
				contactDto := record.Data().(*models4contactus.ContactDbo)
				contactDto.Status = "active"
				contactDto.CountryID = "IE"
				contactDto.CreatedAt = time.Now()
				contactDto.CreatedBy = "u1"
				switch record.Key().ID {
				case Dispatcher1ContactID:
					contactDto.Type = "company"
					contactDto.Title = Dispatcher1ContactTitle
				case Dispatcher1warehouse1ContactID:
					contactDto.ParentID = Dispatcher1ContactID
					contactDto.Type = "location"
					contactDto.Title = "WarehouseOperator 1"
					contactDto.Address = &dbmodels.Address{
						CountryID: "IE",
						Lines:     "WarehouseOperator 1\nDispatcher1\nIreland",
					}
				case Dispatcher2ContactID:
					contactDto.Type = "company"
					contactDto.Title = Dispatcher2ContactTitle
				case Dispatcher2warehouse1ContactID:
					contactDto.Type = "location"
					contactDto.ParentID = Dispatcher2ContactID
					contactDto.Title = "WarehouseOperator 1"
					contactDto.Address = &dbmodels.Address{
						CountryID: "IE",
						Lines:     "WarehouseOperator 1\nDispatcher2\nIreland",
					}
				case Dispatcher2warehouse2ContactID:
					contactDto.Type = "location"
					contactDto.ParentID = Dispatcher2ContactID
					contactDto.Title = "WarehouseOperator 2"
					contactDto.Address = &dbmodels.Address{
						CountryID: "IE",
						Lines:     "WarehouseOperator 2\nDispatcher2\nIreland",
					}
				case Port1ContactID:
					contactDto.Type = "company"
					contactDto.Title = Port1ContactTitle
				case Port1dock1ContactID:
					contactDto.Type = "location"
					contactDto.Title = "Port 1 dock"
					contactDto.ParentID = Port1ContactID
					contactDto.Address = &dbmodels.Address{
						CountryID: "IE",
						Lines:     "Dock 1\nPort 1\nIreland",
					}
				case Port2ContactID:
					contactDto.Type = "company"
					contactDto.Title = Port1ContactTitle
				case Port2dock1ContactID:
					contactDto.Type = "location"
					contactDto.Title = "Port 2 dock 1"
					contactDto.ParentID = Port2ContactID
					contactDto.Address = &dbmodels.Address{
						CountryID: "IE",
						Lines:     "Dock 1\nPort 2\nIreland",
					}
				case Port2dock2ContactID:
					contactDto.Type = "location"
					contactDto.Title = "Port 2 dock 2"
					contactDto.ParentID = Port2ContactID
					contactDto.Address = &dbmodels.Address{
						CountryID: "IE",
						Lines:     "Dock 2\nPort 2\nIreland",
					}
				case Trucker1ContactID:
					contactDto.Type = "company"
					contactDto.Title = "Trucker 1"
				default:
					return dal.ErrRecordNotFound
				}
				dbo4linkage.UpdateRelatedIDs(&contactDto.WithRelated, &contactDto.WithRelatedIDs)
			default:
				t.Fatalf("Unexpected collection: %v", record.Key())
			}
			return nil
		}).AnyTimes()
	return tx
}
