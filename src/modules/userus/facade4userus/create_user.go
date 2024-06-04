package facade4userus

//import (
//	"context"
//	"fmt"
//	"github.com/dal-go/dalgo/dal"
//	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/briefs4contactus"
//	"github.com/sneat-co/sneat-go-backend/src/modules/userus/dto4userus"
//	"github.com/sneat-co/sneat-go-backend/src/modules/userus/dbo4userus"
//	"github.com/sneat-co/sneat-go-core/facade"
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
//	db := facade.GetDatabase(ctx)
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
//	user.Dbo.Created.Client = request.RemoteClient
//	user.Dbo.CreatedAt = time.Now()
//	user.Dbo.CreatedBy = request.RemoteClient.HostOrApp
//	if i := strings.Index(user.Dbo.CreatedBy, ":"); i > 0 {
//		user.Dbo.CreatedBy = user.Dbo.CreatedBy[:i]
//	}
//
//	{ // Set user's names
//		user.Dbo.Names.Full = briefs4contactus.CleanTitle(request.Title)
//		if strings.Contains(user.Dbo.Names.Full, " ") {
//			user.Dbo.Defaults = &dbo4userus.UserDefaults{
//				ShortNames: briefs4contactus.GetShortNames(user.Dbo.Names.Full),
//			}
//		}
//	}
//	user.Dbo.Email = strings.TrimSpace(request.Email)
//	user.Dbo.Emails = []dbmodels.PersonEmail{
//		{Type: "primary", Address: user.Dbo.Email},
//	}
//	if user.Dbo.Gender == "" {
//		user.Dbo.Gender = "unknown"
//	}
//	if user.Dbo.AgeGroup == "" {
//		user.Dbo.AgeGroup = "unknown"
//	}
//	if err := user.Dbo.Validate(); err != nil {
//		return fmt.Errorf("not able to create user record: %w", err)
//	}
//	if err := tx.Insert(ctx, user.Record); err != nil {
//		return fmt.Errorf("failed to create user record: %w", err)
//	}
//	return nil
//}
