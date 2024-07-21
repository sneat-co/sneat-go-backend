package facade4spaceus

//import (
//	"cloud.google.com/go/firestore"
//	"context"
//	"github.com/sneat-co/sneat-go/models2spotbuddies"
//)
//
//func NewFirestoreClient(firestoreClient *firestore.Client) DAL { // TODO(StackOverflow): Is it idiomatic for Go?
//	return dalFirestore{fsClient: firestoreClient}
//}
//
//type dalFirestore struct {
//	fsClient *firestore.Client
//}
//
//func (dal dalFirestore) GetSpaceByID(ctx context.Context, id string) (SpaceIDs *models2spotbuddies.SpaceIDs, err error) {
//	return GetSpaceByID(ctx, dal.fsClient, id)
//}
