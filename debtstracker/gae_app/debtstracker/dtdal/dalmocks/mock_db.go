package dalmocks

//
// import (
// 	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/dtdal"
// 	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/models"
// 	"context"
// 	"github.com/dal-go/dalgo"
// )
//
// const NOT_IMPLEMENTED_YET = "Not implemented yet"
//
// type MockDB struct {
// 	mockdb
// 	ContactMock  *ContactDalMock
// 	BillMock     *BillDalMock
// 	UserMock     *UserDalMock
// 	TransferMock *TransferDalMock
// 	ReminderMock *ReminderDalMock
// 	//TaskQueueMock TaskQueueDalMock
// }
//
// var _ db.Database = (*MockDB)(nil)
//
// func NewMockDB() MockDB {
// 	mockDB := MockDB{
// 		ContactMock:  NewContactDalMock(),
// 		BillMock:     NewBillDalMock(),
// 		UserMock:     NewUserDalMock(),
// 		TransferMock: NewTransferDalMock(),
// 		ReminderMock: NewReminderDalMock(),
// 		//TaskQueueMock: NewTaskQueueDalMock(),
// 	}
//
// 	dtdal.ContactEntry = mockDB.ContactMock
// 	dtdal.Bill = mockDB.BillMock
// 	dtdal.User = mockDB.UserMock
// 	dtdal.Transfer = mockDB.TransferMock
// 	dtdal.Reminder = mockDB.ReminderMock
// 	//dtdal.TaskQueue = mockDB.TaskQueueMock
//
// 	return mockDB
// }
//
// func (mockDB MockDB) Get(c context.Context, entityHolder db.EntityHolder) error {
// 	panic("not implemented yet")
// }
//
// func (mockDB MockDB) Update(c context.Context, entityHolder db.EntityHolder) error {
// 	panic("not implemented yet")
// }
//
// func (mockDB MockDB) IsInTransaction(c context.Context) bool {
// 	panic("not implemented yet")
// }
//
// func (mockDB MockDB) NonTransactionalContext(c context.Context) context.Context {
// 	panic("not implemented yet")
// }
//
// func (mockDB MockDB) Delete(c context.Context, entityHolder db.EntityHolder) error {
// 	panic("not implemented yet")
// }
//
// func (mockDB MockDB) GetMulti(c context.Context, entityHolders []db.EntityHolder) error {
// 	for _, entityHolder := range entityHolders {
// 		switch entityHolder.Kind() {
// 		//case models.CounterpartyKind:
// 		//	if newEntityHolder, err := mockDB.CounterpartyMock.GetCounterpartyByID(c, entityHolder.IntegerID()); err != nil {
// 		//		return err
// 		//	} else {
// 		//		entityHolder.SetEntity(newEntityHolder.Entity())
// 		//	}
// 		case models.BillKind:
// 			if newEntityHolder, err := mockDB.BillMock.GetBillByID(c, entityHolder.StrID()); err != nil {
// 				return err
// 			} else {
// 				entityHolder.SetEntity(newEntityHolder.Entity())
// 			}
// 		case models.AppUserKind:
// 			if newEntityHolder, err := mockDB.UserMock.GetUserByID(c, entityHolder.IntID()); err != nil {
// 				return err
// 			} else {
// 				entityHolder.SetEntity(newEntityHolder.Entity())
// 			}
// 		case models.ContactKind:
// 			if newEntityHolder, err := mockDB.ContactMock.GetContactByID(c, entityHolder.IntID()); err != nil {
// 				return err
// 			} else {
// 				entityHolder.SetEntity(newEntityHolder.Entity())
// 			}
// 		case models.TransferKind:
// 			if newEntityHolder, err := mockDB.TransferMock.GetTransferByID(c, entityHolder.IntID()); err != nil {
// 				return err
// 			} else {
// 				entityHolder.SetEntity(newEntityHolder.Entity())
// 			}
// 		default:
// 			panic("Unsupported kind: " + entityHolder.Kind())
// 		}
// 	}
// 	return nil
// }
//
// func (mockDB MockDB) UpdateMulti(c context.Context, entityHolders []db.EntityHolder) error {
// 	for _, entityHolder := range entityHolders {
// 		switch entityHolder.Kind() {
// 		case models.BillKind:
// 			mockDB.BillMock.Bills[entityHolder.StrID()] = entityHolder.Entity().(*models.BillEntity)
// 		case models.AppUserKind:
// 			mockDB.UserMock.Users[entityHolder.IntID()] = entityHolder.Entity().(*models.AppUserEntity)
// 		case models.ContactKind:
// 			mockDB.ContactMock.Contacts[entityHolder.IntID()] = entityHolder.Entity().(*models.ContactEntity)
// 		case models.TransferKind:
// 			mockDB.TransferMock.Transfers[entityHolder.IntID()] = entityHolder.Entity().(*models.TransferEntity)
// 		default:
// 			panic("Unsupported kind: " + entityHolder.Kind())
// 		}
// 	}
// 	return nil
// }
//
// func (MockDB) RunInTransaction(c context.Context, f func(c context.Context) error, options db.RunOptions) error {
// 	return f(context.WithValue(c, "IsInTransaction", true))
// }
