package maintainance

import (
	"context"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/models4debtus"
)

//type contactsAsyncJob struct {
//	//asyncMapper
//	entity *models.DebtusSpaceContactDbo
//}

//var _ mapper.JobEntity = (*contactsAsyncJob)(nil)

//func (m *contactsAsyncJob) Make() interface{} {
//	m.entity = new(models.DebtusSpaceContactDbo)
//	return m.entity
//}
//
//func (m *contactsAsyncJob) Query(r *http.Request) (query any /* *mapper.Query */, err error) {
//	return nil, errors.New("contactsAsyncJob.Query() is not implemented")
//	//return applyIDAndUserFilters(r, "contactsAsyncJob", models.DebtusContactsCollection, filterByIntID, "UserID")
//}
//
//func (m *contactsAsyncJob) DebtusSpaceContactEntry(key *datastore.Key) (contact models.DebtusSpaceContactEntry) {
//	contact = models.NewDebtusSpaceContactEntry(key.StringID(), nil)
//	if m.entity != nil {
//		entity := *m.entity
//		contact.Data = &entity
//	}
//	return
//}

type ContactWorker func(ctx context.Context, counters any /* *asyncCounters*/, contact models4debtus.DebtusSpaceContactEntry) error

//func (m *contactsAsyncJob) startContactWorker(ctx context.Context, counters mapper.Counters, key *datastore.Key, contactWorker ContactWorker) error {
//	//logus.Debugf(ctx, "*contactsAsyncJob.startContactWorker()")
//	contact := m.DebtusSpaceContactEntry(key)
//	createContactWorker := func() Worker {
//		//logus.Debugf(c, "createContactWorker()")
//		return func(counters *asyncCounters) error {
//			//logus.Debugf(c, "asyncContactWorker() => contact.ContactID: %v", contact.ContactID)
//			return contactWorker(ctx, counters, contact)
//		}
//	}
//	return m.startWorker(c, counters, createContactWorker)
//}
