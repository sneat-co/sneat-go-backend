package maintainance

import (
	"net/http"
)

func RegisterMappers() {
	//mapperServer, _ := mapper.NewServer(
	//	mapper.DefaultPath,
	//	mapper.DefaultQueue(anybot.QUEUE_MAPREDUCE),
	//)
	//http.Handle(mapper.DefaultPath, mapperServer)

	//registerAsyncJob := func(job interface {
	//	mapper.JobSpec
	//	mapper.SliceLifecycle
	//	mapper.JobLifecycle
	//}) {
	//	//mapper.RegisterJob(job)
	//}
	//_ = registerAsyncJob
	//registerAsyncJob(&verifyUsers{})
	//registerAsyncJob(&verifyContacts{})
	//registerAsyncJob(&verifyTransfers{})
	//registerAsyncJob(&migrateTransfers{})
	//registerAsyncJob(&verifyContactTransfers{})
	//registerAsyncJob(&transfersRecreateContacts{})
	//registerAsyncJob(&verifyTelegramUserAccounts{})

	http.HandleFunc("/_ah/merge-contacts", mergeContactsHandler)
}

//func filterByUserParam(query *mapper.Query, pv, prop string) (q *mapper.Query, filtered bool, err error) {
//	return filterByIntParam(query, pv, prop)
//}

//func filterByContactParam(r *http.Request, query *mapper.Query, prop string) (*mapper.Query, error) {
//	return filterByIntParam(r, query, "contact", prop)
//}

//func filterByIntParam(query *mapper.Query, pv, prop string) (q *mapper.Query, filtered bool, err error) {
//	q = query
//	if pv != "" {
//		var v int64
//		if v, err = strconv.ParseInt(pv, 10, 64); err != nil {
//			return
//		} else if v != 0 {
//			return query.Filter(prop+" =", v), true, nil
//		}
//	}
//	return
//}
