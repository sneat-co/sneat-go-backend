package healthcheck

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-core/apicore"
	"github.com/sneat-co/sneat-go-core/apicore/verify"
	"github.com/sneat-co/sneat-go-core/facade"
	"log"
	"net/http"
	"time"
)

// httpGetPage renders health-check page
func httpGetPage(w http.ResponseWriter, r *http.Request) {
	ctx, _, err := apicore.VerifyRequestAndCreateUserContext(w, r, verify.DefaultJsonWithAuthRequired)
	if err != nil {
		log.Println("VerifyRequest error:", err)
		//api4meetingus.ReturnError(ctx, w, err)
		return
	}
	data := healthCheck{
		At: time.Now(),
	}
	key := dal.NewKeyWithID("health_checks", "firestore-write")
	record := dal.NewRecordWithData(key, data)

	db := facade.GetDatabase(ctx)
	err = db.RunReadwriteTransaction(ctx, func(ctx context.Context, tx dal.ReadwriteTransaction) error {
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
		log.Println("Failed to write to output stream: ", err)
	}
}
