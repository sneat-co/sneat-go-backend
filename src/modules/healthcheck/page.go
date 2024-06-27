package healthcheck

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-core/apicore"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/log"
	"net/http"
	"time"
)

// httpGetPage renders health-check page
func httpGetPage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	data := healthCheck{}
	key := dal.NewKeyWithID("health_checks", "firestore-write")
	record := dal.NewRecordWithData(key, &data)

	db := facade.GetDatabase(ctx)
	err := db.RunReadwriteTransaction(ctx, func(ctx context.Context, tx dal.ReadwriteTransaction) error {
		err := tx.Get(ctx, record)
		if err != nil && !dal.IsNotFound(err) {
			return fmt.Errorf("failed to get health check record: %w", err)
		}
		data.At = time.Now()
		return tx.Set(ctx, record)
	})
	if err != nil {
		apicore.ReturnError(ctx, w, r, err)
	}

	w.WriteHeader(http.StatusOK)

	_, err = w.Write([]byte(fmt.Sprintf("Firestore write: OK at %v: id=%v;",
		data.At,
		key.String(),
	)))
	if err != nil {
		log.Errorf(ctx, "Failed to write to output stream: ", err)
	}
}
