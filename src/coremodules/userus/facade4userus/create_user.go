package facade4userus

//import (
//	"context"
//	"fmt"
//	"github.com/dal-go/dalgo/dal"
//	"github.com/sneat-co/sneat-go-backend/src/coremodules/contactus/briefs4contactus"
//	"github.com/sneat-co/sneat-go-backend/src/coremodules/userus/dto4userus"
//	"github.com/sneat-co/sneat-go-backend/src/coremodules/userus/dbo4userus"
//	"github.com/sneat-co/sneat-go-core/facade4debtus"
//	"github.com/sneat-co/sneat-go-core/models/dbmodels"
//	"strings"
//	"time"
//)
//
//// CreateUser creates user record in DB
//func CreateUser(ctx context.Context, userID string, request dto4userus.CreateUserRequestWithRemoteClientInfo) error {
//	if request.Creator != "" { // TODO: document why we do this
//		request.RemoteClient.HostOrApp = request.Creator
//	}
//	db := facade4debtus.GetDatabase(ctx)
//
//	err := db.RunReadwriteTransaction(ctx, func(ctx context.Context, tx dal.ReadwriteTransaction) error {
//		return createUserTx(ctx, tx, userID, request)
//	})
//	if err != nil {
//		return fmt.Errorf("failed to create user record in database: %w", err)
//	}
//	return nil
//}
//
//func createUserTx(ctx context.Context, tx dal.ReadwriteTransaction, userID string, request dto4userus.CreateUserRequestWithRemoteClientInfo) error {
//	user := dbo4userus.NewUserContext(userID)
//	if err := TxGetUserByID(ctx, tx, user.Record); !dal.IsNotFound(err) {
//		return err // Might be nil or not related to "record not found"
//	}
//
//	user.Data.Created.Client = request.RemoteClient
//	user.Data.CreatedAt = time.Now()
//	user.Data.CreatedBy = request.RemoteClient.HostOrApp
//	if i := strings.Index(user.Data.CreatedBy, ":"); i > 0 {
//		user.Data.CreatedBy = user.Data.CreatedBy[:i]
//	}
//
//	{ // Set user's names
//		user.Data.Names.Full = briefs4contactus.CleanTitle(request.Title)
//		if strings.Contains(user.Data.Names.Full, " ") {
//			user.Data.Defaults = &dbo4userus.UserDefaults{
//				ShortNames: briefs4contactus.GetShortNames(user.Data.Names.Full),
//			}
//		}
//	}
//	user.Data.Email = strings.TrimSpace(request.Email)
//	user.Data.Emails = []dbmodels.PersonEmail{
//		{ExtraType: "primary", Address: user.Data.Email},
//	}
//	if user.Data.Gender == "" {
//		user.Data.Gender = "unknown"
//	}
//	if user.Data.AgeGroup == "" {
//		user.Data.AgeGroup = "unknown"
//	}
//	if err := user.Data.Validate(); err != nil {
//		return fmt.Errorf("not able to create user record: %w", err)
//	}
//	if err := tx.Insert(ctx, user.Record); err != nil {
//		return fmt.Errorf("failed to create user record: %w", err)
//	}
//	return nil
//}
