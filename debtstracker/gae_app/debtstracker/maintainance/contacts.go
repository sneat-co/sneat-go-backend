package maintainance

import (
	"context"
	//"github.com/captaincodeman/datastore-mapper"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/models"
)

//type contactsAsyncJob struct {
//	//asyncMapper
//	entity *models.DebtusContactDbo
//}

//var _ mapper.JobEntity = (*contactsAsyncJob)(nil)

//func (m *contactsAsyncJob) Make() interface{} {
//	m.entity = new(models.DebtusContactDbo)
//	return m.entity
//}
//
//func (m *contactsAsyncJob) Query(r *http.Request) (query any /* *mapper.Query */, err error) {
//	return nil, errors.New("contactsAsyncJob.Query() is not implemented")
//	//return applyIDAndUserFilters(r, "contactsAsyncJob", models.DebtusContactsCollection, filterByIntID, "UserID")
//}
//
//func (m *contactsAsyncJob) ContactEntry(key *datastore.Key) (contact models.ContactEntry) {
//	contact = models.NewDebtusContact(key.StringID(), nil)
//	if m.entity != nil {
//		entity := *m.entity
//		contact.Data = &entity
//	}
//	return
//}

type ContactWorker func(c context.Context, counters any /* *asyncCounters*/, contact models.ContactEntry) error

//func (m *contactsAsyncJob) startContactWorker(c context.Context, counters mapper.Counters, key *datastore.Key, contactWorker ContactWorker) error {
//	//log.Debugf(c, "*contactsAsyncJob.startContactWorker()")
//	contact := m.ContactEntry(key)
//	createContactWorker := func() Worker {
//		//log.Debugf(c, "createContactWorker()")
//		return func(counters *asyncCounters) error {
//			//log.Debugf(c, "asyncContactWorker() => contact.ID: %v", contact.ID)
//			return contactWorker(c, counters, contact)
//		}
//	}
//	return m.startWorker(c, counters, createContactWorker)
//}
