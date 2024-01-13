package facade4teamus

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
//func (dal dalFirestore) GetTeamByID(ctx context.Context, id string) (TeamIDs *models2spotbuddies.TeamIDs, err error) {
//	return GetTeamByID(ctx, dal.fsClient, id)
//}
