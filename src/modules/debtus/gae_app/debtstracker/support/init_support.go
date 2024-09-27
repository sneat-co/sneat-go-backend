package support

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func InitSupportHandlers(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, "/support/validate-users", ValidateUsersHandler)
	router.HandlerFunc(http.MethodGet, "/support/validate-user", ValidateUserHandler)
}

func ValidateUsersHandler(w http.ResponseWriter, r *http.Request) {
	panic("implement me")
	//c := r.Context()
	//fix := r.URL.Query().Get("fix")
	//query := datastore.NewQuery(models4debtus.AppUserKind).KeysOnly() //.Limit(25)
	//t := query.Run(c)
	//batchSize := 100
	//tasks := make([]*taskqueue.Task, 0, batchSize)
	//var (
	//	usersCount int
	//	params     url.Values
	//)
	//
	//addTasksToQueue := func() error {
	//	if _, err := taskqueue.AddMulti(c, tasks, "support"); err != nil {
	//		logus.Errorf(c, "Failed to add tasks: %v", err)
	//		return err
	//	}
	//	tasks = make([]*taskqueue.Task, 0, batchSize)
	//	return nil
	//}
	//
	//for {
	//	if key, err := t.Next(nil); err != nil {
	//		if err == datastore.Done {
	//			break
	//		}
	//		logus.Errorf(c, "Failed to fetch %v: %v", key, err)
	//		return
	//	} else {
	//		usersCount += 1
	//		taskUrl := fmt.Sprintf("/support/validate-user?id=%v", key.IntID())
	//		if fix != "" {
	//			taskUrl += "&fix=" + fix
	//		}
	//		tasks = append(tasks, taskqueue.NewPOSTTask(taskUrl, params))
	//		if len(tasks) == batchSize {
	//			if err = addTasksToQueue(); err != nil {
	//				return
	//			}
	//		}
	//	}
	//
	//}
	//if len(tasks) > 0 {
	//	if err := addTasksToQueue(); err != nil {
	//		return
	//	}
	//}
	//logus.Errorf(c, "(NOT error) Users count: %v", usersCount)
	//_, _ = w.Write([]byte(fmt.Sprintf("Users count: %v", usersCount)))
}

//type int64sortable []int64
//
//func (a int64sortable) Len() int           { return len(a) }
//func (a int64sortable) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
//func (a int64sortable) Less(i, j int) bool { return a[i] < a[j] }
