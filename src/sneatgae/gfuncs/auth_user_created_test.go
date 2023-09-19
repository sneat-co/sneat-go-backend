package gfuncs

import (
	"cloud.google.com/go/firestore"
	"context"
	"testing"
)

type dbClientMock struct {
	txRunCount int
}

func (db *dbClientMock) RunTransaction(_ context.Context, _ func(ctx context.Context, transaction *firestore.Transaction) error, opts ...firestore.TransactionOption) error {
	db.txRunCount++
	return nil
}

func (db *dbClientMock) Collection(path string) *firestore.CollectionRef {
	//TODO implement me
	panic("implement me")
}

var _ dbClient = (*dbClientMock)(nil)

func TestAuthUserCreated(t *testing.T) {
	dbMock := new(dbClientMock)
	db = dbMock
	if err := AuthUserCreated(context.Background(), AuthEvent{}); err != nil {
		t.Fatalf("failed on a call with empty AuthEvent{}: %v", err)
	}
	if dbMock.txRunCount > 0 {
		t.Fatalf("was not expecting calls to RunTransaction but got %v", dbMock.txRunCount)
	}
	if err := AuthUserCreated(context.Background(), AuthEvent{UID: "test-user"}); err != nil {
		t.Fatalf("failed on a call with empty AuthEvent{}: %v", err)
	}
	if dbMock.txRunCount == 0 {
		t.Fatal("transaction not called")
	}
	if dbMock.txRunCount > 1 {
		t.Fatalf("db transaction expected to be called just once, got caller %v times", dbMock.txRunCount)
	}
}
