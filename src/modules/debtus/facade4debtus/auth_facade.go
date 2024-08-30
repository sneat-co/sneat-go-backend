package facade4debtus

import (
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/auth/models4auth"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/gae_app/debtstracker/dtdal"
	"github.com/sneat-co/sneat-go-core/facade"
	"math/rand"
	"time"

	"context"
	"errors"
)

type authFacade struct {
}

var AuthFacade = authFacade{}

func (authFacade) AssignPinCode(ctx context.Context, loginID int, userID string) (loginPin models4auth.LoginPin, err error) {
	err = facade.RunReadwriteTransaction(ctx, func(ctx context.Context, tx dal.ReadwriteTransaction) error {
		if loginPin, err = dtdal.LoginPin.GetLoginPinByID(ctx, tx, loginID); err != nil {
			return fmt.Errorf("failed to get LoginPin entity by loginID=%d: %w", loginID, err)
		}
		if loginPin.Data.UserID != "" && loginPin.Data.UserID != userID {
			return errors.New("LoginPin.UserID != userID")
		}
		if !loginPin.Data.SignedIn.IsZero() {
			return errors.New("LoginPin.SignedIn.IsZero(): false")
		}
		random := rand.New(rand.NewSource(time.Now().UnixNano()))
		loginPin.Data.Code = random.Int31n(9000) + 1000
		loginPin.Data.UserID = userID
		loginPin.Data.Pinned = time.Now()
		if err = dtdal.LoginPin.SaveLoginPin(ctx, tx, loginPin); err != nil {
			return fmt.Errorf("failed to save LoginPin entity with ContactID=%v: %w", loginID, err)
		}
		return err
	}, nil)
	return
}

func (authFacade) SignInWithPin(ctx context.Context, loginID int, loginPinCode int32) (userID string, err error) {
	_ = loginPinCode
	err = facade.RunReadwriteTransaction(ctx, func(ctx context.Context, tx dal.ReadwriteTransaction) error {
		var loginPin models4auth.LoginPin
		if loginPin, err = dtdal.LoginPin.GetLoginPinByID(ctx, tx, loginID); err != nil {
			return fmt.Errorf("failed to get LoginPin entity by loginID=%d: %w", loginID, err)
		}
		if !loginPin.Data.SignedIn.IsZero() {
			return ErrLoginAlreadySigned
		}
		if loginPin.Data.Created.Add(time.Hour).Before(time.Now()) {
			return ErrLoginExpired
		}
		if userID = loginPin.Data.UserID; userID == "" {
			return errors.New("LoginPin.UserID == 0")
		}

		loginPin.Data.SignedIn = time.Now()
		if err = dtdal.LoginPin.SaveLoginPin(ctx, tx, loginPin); err != nil {
			return err
		}
		return err
	}, nil) // dtdal.CrossGroupTransaction)
	return
}
