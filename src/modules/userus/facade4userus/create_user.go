package facade4userus

//import (
//	"context"
//	"fmt"
//	"github.com/dal-go/dalgo/dal"
//	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/briefs4contactus"
//	"github.com/sneat-co/sneat-go-backend/src/modules/userus/dto4userus"
//	"github.com/sneat-co/sneat-go-backend/src/modules/userus/models4userus"
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
//	user := models4userus.NewUserContext(userID)
//	if err := TxGetUserByID(ctx, tx, user.Record); !dal.IsNotFound(err) {
//		return err // Might be nil or not related to "record not found"
//	}
//
//	user.Dto.Created.Client = request.RemoteClient
//	user.Dto.CreatedAt = time.Now()
//	user.Dto.CreatedBy = request.RemoteClient.HostOrApp
//	if i := strings.Index(user.Dto.CreatedBy, ":"); i > 0 {
//		user.Dto.CreatedBy = user.Dto.CreatedBy[:i]
//	}
//
//	{ // Set user's names
//		user.Dto.Names.Full = briefs4contactus.CleanTitle(request.Title)
//		if strings.Contains(user.Dto.Names.Full, " ") {
//			user.Dto.Defaults = &models4userus.UserDefaults{
//				ShortNames: briefs4contactus.GetShortNames(user.Dto.Names.Full),
//			}
//		}
//	}
//	user.Dto.Email = strings.TrimSpace(request.Email)
//	user.Dto.Emails = []dbmodels.PersonEmail{
//		{Type: "primary", Address: user.Dto.Email},
//	}
//	if user.Dto.Gender == "" {
//		user.Dto.Gender = "unknown"
//	}
//	if user.Dto.AgeGroup == "" {
//		user.Dto.AgeGroup = "unknown"
//	}
//	if err := user.Dto.Validate(); err != nil {
//		return fmt.Errorf("not able to create user record: %w", err)
//	}
//	if err := tx.Insert(ctx, user.Record); err != nil {
//		return fmt.Errorf("failed to create user record: %w", err)
//	}
//	return nil
//}
