package maintainance

//import (
//	"context"
//	"github.com/sanity-io/litter"
//	"google.golang.org/appengine/v2/datastore"
//	"google.golang.org/appengine/v2/log"
//)
//
//func logJobCompletion(c context.Context, id string) {
//	logus.Infof(c, "Job completed: %v", id)
//	key := datastore.NewKey(c, "MP_job", id, 0, nil)
//
//	var props datastore.PropertyList
//	if err := datastore.Get(c, key, &props); err != nil {
//		logus.Errorf(c, "Failed to get job entity by ID=%v: %v", id, err)
//	} else {
//		logus.Debugf(c, "Job entity: %v", litter.Sdump(props))
//	}
//}
