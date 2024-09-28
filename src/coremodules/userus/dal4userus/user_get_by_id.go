package dal4userus

import (
	"context"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/coremodules/userus/dbo4userus"
	"github.com/sneat-co/sneat-go-core/facade"
)

// GetUserByID returns user entry by its ID
var GetUserByID = func(ctx context.Context, tx dal.ReadSession, userID string) (user dbo4userus.UserEntry, err error) {
	user = dbo4userus.NewUserEntry(userID)
	return user, GetUser(ctx, tx, user)
}

// GetUser returns user by its ID
func GetUser(ctx context.Context, tx dal.ReadSession, user dbo4userus.UserEntry) (err error) {
	return tx.Get(ctx, user.Record)
}

// GetUsersByIDs returns users by their IDs
func GetUsersByIDs(ctx context.Context, userIDs []string) (users []dbo4userus.UserEntry, err error) {
	//logus.Debugf(ctx, "UserDalGae.GetUsersByIDs(%d)", userIDs)
	if len(userIDs) == 0 {
		return
	}

	appUsers := NewUserEntries(userIDs)
	records := UserRecords(appUsers)
	var db dal.DB
	if db, err = facade.GetSneatDB(ctx); err != nil {
		return
	}
	if err = db.GetMulti(ctx, records); err != nil {
		return
	}
	return
}

// NewUserEntries creates slice of user entries by user IDs
func NewUserEntries(userIDs []string) []dbo4userus.UserEntry {
	users := make([]dbo4userus.UserEntry, len(userIDs))
	for i, id := range userIDs {
		users[i] = dbo4userus.NewUserEntry(id)
	}
	return users
}

// UserRecords converts user entries to slice of Dalgo records
func UserRecords(users []dbo4userus.UserEntry) (records []dal.Record) {
	records = make([]dal.Record, len(users))
	for i, u := range users {
		records[i] = u.Record
	}
	return
}
