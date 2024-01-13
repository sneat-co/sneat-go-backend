package facade4contactus

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-core/facade/db"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"time"
)

var txUpdate = func(ctx context.Context, tx dal.ReadwriteTransaction, key *dal.Key, data []dal.Update, opts ...dal.Precondition) error {
	return db.TxUpdate(ctx, tx, key, data, opts...)
}

// UpdateMemberGroup weird unclear method - TODO: remove/replace or document to make sense
func txUpdateMemberGroup(ctx context.Context, tx dal.ReadwriteTransaction, updatedAt time.Time, updatedBy string, membersGroup dbmodels.Versioned, key *dal.Key, data []dal.Update, opts ...dal.Precondition) error {
	version := membersGroup.IncreaseVersion(updatedAt, updatedBy)
	if err := membersGroup.Validate(); err != nil {
		return fmt.Errorf("membersGroup record is not valid: %w", err)
	}
	data = append(data,
		dal.Update{Field: "v", Value: version},
		dal.Update{Field: "updatedAt", Value: updatedAt},
		dal.Update{Field: "updatedBy", Value: updatedBy},
	)
	return txUpdate(ctx, tx, key, data, opts...)
}
