package dalmocks

//
//import (
//	"github.com/sneat-co/sneat-go-backend/debtusbot/gae_app/debtusbot/dtdal"
//	"github.com/sneat-co/sneat-go-backend/debtusbot/gae_app/debtusbot/models"
//	"context"
//	"github.com/bots-go-framework/bots-fw/botsfw"
//)
//
//var _ dtdal.ContactDal = (*ContactDalMock)(nil)
//
//type ContactDalMock struct {
//	LastContactID int64
//	Contacts      map[int64]*models.ContactEntity
//}
//
//func NewContactDalMock() *ContactDalMock {
//	return &ContactDalMock{Contacts: make(map[int64]*models.ContactEntity)}
//}
//
//func (mock *ContactDalMock) GetLatestContacts(whc botsfw.WebhookContext, limit, totalCount int) (contacts []models.DebtusSpaceContactEntry, err error) {
//	_, _, _ = whc, limit, totalCount
//	return
//}
//
//func (mock *ContactDalMock) InsertContact(_ context.Context, contactEntity *models.ContactEntity) (contact models.DebtusSpaceContactEntry, err error) {
//	if contactEntity == nil {
//		panic("contactEntity == nil")
//	}
//	mock.LastContactID += 1
//	contact.ContactID = mock.LastContactID
//	contact.Data = contactEntity
//	mock.Contacts[mock.LastContactID] = contact.Data
//	return
//}
//
////CreateContact(ctx context.Context, userID int64, contactDetails models.ContactDetails) (contact models.DebtusSpaceContactEntry, user models.AppUser, err error)
////CreateContactWithinTransaction(ctx context.Context, user models.AppUser, contactUserID, counterpartyCounterpartyID int64, contactDetails models.ContactDetails, balanced money.Balanced) (contact models.DebtusSpaceContactEntry, err error)
////UpdateContact(ctx context.Context, contactID int64, values map[string]string) (contactEntity *models.ContactEntity, err error)
//
//func (mock *ContactDalMock) SaveContact(_ context.Context, contact models.DebtusSpaceContactEntry) (err error) {
//	mock.Contacts[contact.ContactID] = contact.Data
//	return
//}
//
//func (mock *ContactDalMock) DeleteContact(_ context.Context, contactID int64) (err error) {
//	delete(mock.Contacts, contactID)
//	return
//}
//
//func (mock *ContactDalMock) GetContactIDsByTitle(_ context.Context, userID int64, title string, caseSensitive bool) (contactIDs []int64, err error) {
//	return
//}
//
//func (mock *ContactDalMock) GetContactsWithDebts(_ context.Context, userID int64) (contacts []models.DebtusSpaceContactEntry, err error) {
//	return
//}
