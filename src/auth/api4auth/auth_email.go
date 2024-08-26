package api4auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/auth/facade4auth"
	"github.com/sneat-co/sneat-go-backend/src/auth/models4auth"
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

func HandleSignUpWithEmail(c context.Context, w http.ResponseWriter, r *http.Request) {
	email := strings.TrimSpace(r.PostFormValue("email"))
	userName := strings.TrimSpace(r.PostFormValue("name"))

	if email == "" {
		api4debtus.BadRequestMessage(c, w, "Missing required value: email")
		return
	}

	if err := validateEmail(email); err != nil {
		api4debtus.ErrorAsJson(c, w, http.StatusBadRequest, err)
		return
	}

	if _, err := facade4auth.UserEmail.GetUserEmailByID(c, nil, email); err != nil {
		if !dal.IsNotFound(err) {
			api4debtus.ErrorAsJson(c, w, http.StatusInternalServerError, err)
			return
		} else {
			api4debtus.ErrorAsJson(c, w, http.StatusConflict, facade4debtus.ErrEmailAlreadyRegistered)
			return
		}
	}

	if user, userEmail, err := facade4debtus.User.CreateUserByEmail(c, email, userName); err != nil {
		if errors.Is(err, facade4debtus.ErrEmailAlreadyRegistered) {
			api4debtus.ErrorAsJson(c, w, http.StatusConflict, err)
			return
		} else {
			api4debtus.ErrorAsJson(c, w, http.StatusInternalServerError, err)
			return
		}
	} else {
		if err = emailing.CreateConfirmationEmailAndQueueForSending(c, user, userEmail); err != nil {
			api4debtus.ErrorAsJson(c, w, http.StatusInternalServerError, err)
			return
		}
		api4debtus.ReturnToken(c, w, user.ID, r.Referer(), true, false /*user.Data.Email == "alexander.trakhimenok@gmail.com"*/)
	}
}

func HandleSignInWithEmail(c context.Context, w http.ResponseWriter, r *http.Request) {
	email := strings.TrimSpace(r.PostFormValue("email"))
	password := strings.TrimSpace(r.PostFormValue("password"))
	//logus.Debugf(c, "Email: %s", email)
	if email == "" || password == "" {
		api4debtus.ErrorAsJson(c, w, http.StatusBadRequest, errors.New("Missing required value"))
		return
	}

	if err := validateEmail(email); err != nil {
		api4debtus.JsonToResponse(c, w, map[string]string{"error": err.Error()})
		return
	}

	userEmail, err := facade4auth.UserEmail.GetUserEmailByID(c, nil, email)
	if err != nil {
		if dal.IsNotFound(err) {
			api4debtus.ErrorAsJson(c, w, http.StatusForbidden, errors.New("unknown email"))
		} else {
			api4debtus.ErrorAsJson(c, w, http.StatusInternalServerError, err)
		}
		return
	} else if err = userEmail.Data.CheckPassword(password); err != nil {
		logus.Debugf(c, "Invalid password: %v", err)
		api4debtus.ErrorAsJson(c, w, http.StatusForbidden, errors.New("invalid password"))
		return
	}

	api4debtus.ReturnToken(c, w, userEmail.Data.AppUserID, r.Referer(), false, userEmail.ID == "alexander.trakhimenok@gmail.com")
}

func HandleRequestPasswordReset(c context.Context, w http.ResponseWriter, r *http.Request) {
	email := r.PostFormValue("email")
	userEmail, err := facade4auth.UserEmail.GetUserEmailByID(c, nil, email)
	if dal.IsNotFound(err) {
		api4debtus.ErrorAsJson(c, w, http.StatusForbidden, errors.New("Unknown email"))
		return
	}

	now := time.Now()

	pwdResetEntity := models4auth.PasswordResetData{
		Email:             userEmail.ID,
		Status:            "created",
		OwnedByUserWithID: appuser.NewOwnedByUserWithID(userEmail.Data.AppUserID, now),
	}

	err = facade.RunReadwriteTransaction(c, func(c context.Context, tx dal.ReadwriteTransaction) (err error) {
		_, err = facade4auth.PasswordReset.CreatePasswordResetByID(c, tx, &pwdResetEntity)
		return err
	})
	if err != nil {
		api4debtus.ErrorAsJson(c, w, http.StatusInternalServerError, err)
		return
	}
}

func HandleChangePasswordAndSignIn(c context.Context, w http.ResponseWriter, r *http.Request) {
	var (
		err           error
		passwordReset models4auth.PasswordReset
	)

	if passwordReset.ID, err = strconv.Atoi(r.PostFormValue("pin")); err != nil {
		api4debtus.ErrorAsJson(c, w, http.StatusBadRequest, err)
		return
	}

	pwd := r.PostFormValue("pwd")
	if pwd == "" {
		api4debtus.ErrorAsJson(c, w, http.StatusBadRequest, errors.New("Empty password"))
		return
	}

	if passwordReset, err = facade4auth.PasswordReset.GetPasswordResetByID(c, nil, passwordReset.ID); err != nil {
		if dal.IsNotFound(err) {
			api4debtus.ErrorAsJson(c, w, http.StatusForbidden, errors.New("Unknown pin"))
			return
		}
		api4debtus.ErrorAsJson(c, w, http.StatusInternalServerError, err)
		return
	}

	isAdmin := api4debtus.IsAdmin(passwordReset.Data.Email)

	if err = facade.RunReadwriteTransaction(c, func(c context.Context, tx dal.ReadwriteTransaction) error {

		now := time.Now()
		appUser := dbo4userus.NewUserEntry(passwordReset.Data.AppUserID)
		userEmail := models4auth.NewUserEmail(passwordReset.Data.Email, nil)

		records := []dal.Record{appUser.Record, userEmail.Record, passwordReset.Record}

		//var db dal.DB
		//if db, err = facade.GetDatabase(c); err != nil {
		//	return err
		//}
		if err = tx.GetMulti(c, records); err != nil {
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
		userEmail.Data.SetLastLogin(now)
		appUser.Data.SetLastLogin(now)

		if err = tx.SetMulti(c, records); err != nil {
			return err
		}
		return err
	}); err != nil {
		api4debtus.ErrorAsJson(c, w, http.StatusInternalServerError, err)
		return
	}

	api4debtus.ReturnToken(c, w, passwordReset.Data.AppUserID, r.Referer(), false, isAdmin)
}

var errInvalidEmailConformationPin = errors.New("email confirmation pin is not valid")

func HandleConfirmEmailAndSignIn(c context.Context, w http.ResponseWriter, r *http.Request) {
	var (
		err       error
		userEmail models4auth.UserEmailEntry
		pin       string
	)

	userEmail.ID, pin = r.PostFormValue("email"), r.PostFormValue("pin")

	if userEmail.ID == "" {
		api4debtus.ErrorAsJson(c, w, http.StatusBadRequest, errors.New("Empty email"))
		return
	}
	if pin == "" {
		api4debtus.ErrorAsJson(c, w, http.StatusBadRequest, errors.New("Empty pin"))
		return
	}

	if err = facade.RunReadwriteTransaction(c, func(c context.Context, tx dal.ReadwriteTransaction) error {
		now := time.Now()

		if userEmail, err = facade4auth.UserEmail.GetUserEmailByID(c, tx, userEmail.ID); err != nil {
			return err
		}

		var appUser dbo4userus.UserEntry
		if appUser, err = dal4userus.GetUserByID(c, tx, userEmail.Data.AppUserID); err != nil {
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
		userEmail.Data.SetLastLogin(now)
		appUser.Data.SetLastLogin(now)

		entities := []dal.Record{appUser.Record, userEmail.Record}
		if err = tx.SetMulti(c, entities); err != nil {
			return err
		}
		return err
	}); err != nil {
		if dal.IsNotFound(err) {
			api4debtus.ErrorAsJson(c, w, http.StatusBadRequest, err)
			return
		} else if err == errInvalidEmailConformationPin {
			api4debtus.ErrorAsJson(c, w, http.StatusForbidden, err)
			return
		}
		api4debtus.ErrorAsJson(c, w, http.StatusInternalServerError, err)
		return
	}

	api4debtus.ReturnToken(c, w, userEmail.Data.AppUserID, r.Referer(), false, api4debtus.IsAdmin(userEmail.ID))
}
