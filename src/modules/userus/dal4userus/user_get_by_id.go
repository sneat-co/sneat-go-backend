package dal4userus

import (
	"context"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/dbo4userus"
	"github.com/sneat-co/sneat-go-core/facade"
)

var GetUserByID = func(ctx context.Context, tx dal.ReadSession, userID string) (user dbo4userus.UserEntry, err error) {
	user = dbo4userus.NewUserEntry(userID)
	return user, GetUser(ctx, tx, user)
}

func GetUser(ctx context.Context, tx dal.ReadSession, user dbo4userus.UserEntry) (err error) {
	return tx.Get(ctx, user.Record)
}

func GetUsersByIDs(c context.Context, userIDs []string) (users []dbo4userus.UserEntry, err error) {
	//logus.Debugf(c, "UserDalGae.GetUsersByIDs(%d)", userIDs)
	if len(userIDs) == 0 {
		return
	}

	appUsers := NewUserEntries(userIDs)
	records := UserRecords(appUsers)
	var db dal.DB
	if db, err = facade.GetDatabase(c); err != nil {
		return
	}
	if err = db.GetMulti(c, records); err != nil {
		return
	}
	return
}

func NewUserEntries(userIDs []string) []dbo4userus.UserEntry {
	users := make([]dbo4userus.UserEntry, len(userIDs))
	for i, id := range userIDs {
		users[i] = dbo4userus.NewUserEntry(id)
	}
	return users
}

func UserRecords(users []dbo4userus.UserEntry) (records []dal.Record) {
	records = make([]dal.Record, len(users))
	for i, u := range users {
		records[i] = u.Record
	}
	return
}
