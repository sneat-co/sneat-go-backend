package api4auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/auth/facade4auth"
	"github.com/sneat-co/sneat-go-backend/src/modules/auth/models4auth"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/facade4debtus"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/gae_app/debtstracker/api4debtus"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/gae_app/debtstracker/emailing"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/dal4userus"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/dbo4userus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/logus"
	"github.com/strongo/strongoapp/appuser"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	reEmail         = regexp.MustCompile(`.+@.+\.\w+`)
	ErrInvalidEmail = errors.New("Invalid email")
)

func validateEmail(email string) error {
	if !reEmail.MatchString(email) {
		return ErrInvalidEmail
	}
	return nil
}

func HandleSignUpWithEmail(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	email := strings.TrimSpace(r.PostFormValue("email"))
	userName := strings.TrimSpace(r.PostFormValue("name"))

	if email == "" {
		api4debtus.BadRequestMessage(ctx, w, "Missing required value: email")
		return
	}

	if err := validateEmail(email); err != nil {
		api4debtus.ErrorAsJson(ctx, w, http.StatusBadRequest, err)
		return
	}

	if _, err := facade4auth.UserEmail.GetUserEmailByID(ctx, nil, email); err != nil {
		if !dal.IsNotFound(err) {
			api4debtus.ErrorAsJson(ctx, w, http.StatusInternalServerError, err)
			return
		} else {
			api4debtus.ErrorAsJson(ctx, w, http.StatusConflict, facade4debtus.ErrEmailAlreadyRegistered)
			return
		}
	}

	if user, userEmail, err := facade4debtus.User.CreateUserByEmail(ctx, email, userName); err != nil {
		if errors.Is(err, facade4debtus.ErrEmailAlreadyRegistered) {
			api4debtus.ErrorAsJson(ctx, w, http.StatusConflict, err)
			return
		} else {
			api4debtus.ErrorAsJson(ctx, w, http.StatusInternalServerError, err)
			return
		}
	} else {
		if err = emailing.CreateConfirmationEmailAndQueueForSending(ctx, user, userEmail); err != nil {
			api4debtus.ErrorAsJson(ctx, w, http.StatusInternalServerError, err)
			return
		}
		ReturnToken(ctx, w, r, user.ID, r.Referer() /*, user.Data.Email == "alexander.trakhimenok@gmail.com"*/)
	}
}

func HandleSignInWithEmail(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	email := strings.TrimSpace(r.PostFormValue("email"))
	password := strings.TrimSpace(r.PostFormValue("password"))
	//logus.Debugf(ctx, "Email: %s", email)
	if email == "" || password == "" {
		api4debtus.ErrorAsJson(ctx, w, http.StatusBadRequest, errors.New("Missing required value"))
		return
	}

	if err := validateEmail(email); err != nil {
		api4debtus.JsonToResponse(ctx, w, map[string]string{"error": err.Error()})
		return
	}

	userEmail, err := facade4auth.UserEmail.GetUserEmailByID(ctx, nil, email)
	if err != nil {
		if dal.IsNotFound(err) {
			api4debtus.ErrorAsJson(ctx, w, http.StatusForbidden, errors.New("unknown email"))
		} else {
			api4debtus.ErrorAsJson(ctx, w, http.StatusInternalServerError, err)
		}
		return
	} else if err = userEmail.Data.CheckPassword(password); err != nil {
		logus.Debugf(ctx, "Invalid password: %v", err)
		api4debtus.ErrorAsJson(ctx, w, http.StatusForbidden, errors.New("invalid password"))
		return
	}

	ReturnToken(ctx, w, r, userEmail.Data.AppUserID, r.Referer() /*, userEmail.ID == "alexander.trakhimenok@gmail.com"*/)
}

func HandleRequestPasswordReset(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	email := r.PostFormValue("email")
	userEmail, err := facade4auth.UserEmail.GetUserEmailByID(ctx, nil, email)
	if dal.IsNotFound(err) {
		api4debtus.ErrorAsJson(ctx, w, http.StatusForbidden, errors.New("Unknown email"))
		return
	}

	now := time.Now()

	pwdResetEntity := models4auth.PasswordResetData{
		Email:             userEmail.ID,
		Status:            "created",
		OwnedByUserWithID: appuser.NewOwnedByUserWithID(userEmail.Data.AppUserID, now),
	}

	err = facade.RunReadwriteTransaction(ctx, func(ctx context.Context, tx dal.ReadwriteTransaction) (err error) {
		_, err = facade4auth.PasswordReset.CreatePasswordResetByID(ctx, tx, &pwdResetEntity)
		return err
	})
	if err != nil {
		api4debtus.ErrorAsJson(ctx, w, http.StatusInternalServerError, err)
		return
	}
}

func HandleChangePasswordAndSignIn(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var (
		err           error
		passwordReset models4auth.PasswordReset
	)

	if passwordReset.ID, err = strconv.Atoi(r.PostFormValue("pin")); err != nil {
		api4debtus.ErrorAsJson(ctx, w, http.StatusBadRequest, err)
		return
	}

	pwd := r.PostFormValue("pwd")
	if pwd == "" {
		api4debtus.ErrorAsJson(ctx, w, http.StatusBadRequest, errors.New("Empty password"))
		return
	}

	if passwordReset, err = facade4auth.PasswordReset.GetPasswordResetByID(ctx, nil, passwordReset.ID); err != nil {
		if dal.IsNotFound(err) {
			api4debtus.ErrorAsJson(ctx, w, http.StatusForbidden, errors.New("Unknown pin"))
			return
		}
		api4debtus.ErrorAsJson(ctx, w, http.StatusInternalServerError, err)
		return
	}

	if err = facade.RunReadwriteTransaction(ctx, func(ctx context.Context, tx dal.ReadwriteTransaction) error {

		now := time.Now()
		appUser := dbo4userus.NewUserEntry(passwordReset.Data.AppUserID)
		userEmail := models4auth.NewUserEmail(passwordReset.Data.Email, nil)

		records := []dal.Record{appUser.Record, userEmail.Record, passwordReset.Record}

		//var db dal.DB
		//if db, err = facade.GetSneatDB(ctx); err != nil {
		//	return err
		//}
		if err = tx.GetMulti(ctx, records); err != nil {
			return err
		}

		if err = userEmail.Data.SetPassword(pwd); err != nil {
			return err
		}

		passwordReset.Data.Status = "changed"
		passwordReset.Data.Email = "" // Clean email as we don't need it anymore
		passwordReset.Data.UpdatedAt = now
		if changed := userEmail.Data.AddProvider("password-reset"); changed {
			userEmail.Data.UpdatedAt = now
		}
		userEmail.Data.SetLastLoginAt(now)
		appUser.Data.SetLastLoginAt(now)

		if err = tx.SetMulti(ctx, records); err != nil {
			return err
		}
		return err
	}); err != nil {
		api4debtus.ErrorAsJson(ctx, w, http.StatusInternalServerError, err)
		return
	}

	ReturnToken(ctx, w, r, passwordReset.Data.AppUserID, r.Referer())
}

var errInvalidEmailConformationPin = errors.New("email confirmation pin is not valid")

func HandleConfirmEmailAndSignIn(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var (
		err       error
		userEmail models4auth.UserEmailEntry
		pin       string
	)

	userEmail.ID, pin = r.PostFormValue("email"), r.PostFormValue("pin")

	if userEmail.ID == "" {
		api4debtus.ErrorAsJson(ctx, w, http.StatusBadRequest, errors.New("Empty email"))
		return
	}
	if pin == "" {
		api4debtus.ErrorAsJson(ctx, w, http.StatusBadRequest, errors.New("Empty pin"))
		return
	}

	if err = facade.RunReadwriteTransaction(ctx, func(ctx context.Context, tx dal.ReadwriteTransaction) error {
		now := time.Now()

		if userEmail, err = facade4auth.UserEmail.GetUserEmailByID(ctx, tx, userEmail.ID); err != nil {
			return err
		}

		var appUser dbo4userus.UserEntry
		if appUser, err = dal4userus.GetUserByID(ctx, tx, userEmail.Data.AppUserID); err != nil {
			return err
		}

		if userEmail.Data.ConfirmationPin() != pin {
			return errInvalidEmailConformationPin
		}

		userEmail.Data.IsConfirmed = true
		if err = userEmail.Data.SetUpdatedTime(now); err != nil {
			return fmt.Errorf("failed to set update time stamp: %w", err)
		}
		userEmail.Data.PasswordBcryptHash = []byte{}
		userEmail.Data.SetLastLoginAt(now)
		appUser.Data.SetLastLoginAt(now)

		entities := []dal.Record{appUser.Record, userEmail.Record}
		if err = tx.SetMulti(ctx, entities); err != nil {
			return err
		}
		return err
	}); err != nil {
		if dal.IsNotFound(err) {
			api4debtus.ErrorAsJson(ctx, w, http.StatusBadRequest, err)
			return
		} else if errors.Is(err, errInvalidEmailConformationPin) {
			api4debtus.ErrorAsJson(ctx, w, http.StatusForbidden, err)
			return
		}
		api4debtus.ErrorAsJson(ctx, w, http.StatusInternalServerError, err)
		return
	}

	ReturnToken(ctx, w, r, userEmail.Data.AppUserID, r.Referer())
}
