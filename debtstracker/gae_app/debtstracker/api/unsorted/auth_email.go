package unsorted

import (
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/api"
	"github.com/strongo/logus"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"context"
	"errors"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/dtdal"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/emailing"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/facade"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/models"
	"github.com/strongo/strongoapp/appuser"
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
		api.BadRequestMessage(c, w, "Missing required value: email")
		return
	}

	if err := validateEmail(email); err != nil {
		api.ErrorAsJson(c, w, http.StatusBadRequest, err)
		return
	}

	if _, err := dtdal.UserEmail.GetUserEmailByID(c, nil, email); err != nil {
		if !dal.IsNotFound(err) {
			api.ErrorAsJson(c, w, http.StatusInternalServerError, err)
			return
		} else {
			api.ErrorAsJson(c, w, http.StatusConflict, facade.ErrEmailAlreadyRegistered)
			return
		}
	}

	if user, userEmail, err := facade.User.CreateUserByEmail(c, email, userName); err != nil {
		if errors.Is(err, facade.ErrEmailAlreadyRegistered) {
			api.ErrorAsJson(c, w, http.StatusConflict, err)
			return
		} else {
			api.ErrorAsJson(c, w, http.StatusInternalServerError, err)
			return
		}
	} else {
		if err = emailing.CreateConfirmationEmailAndQueueForSending(c, user, userEmail); err != nil {
			api.ErrorAsJson(c, w, http.StatusInternalServerError, err)
			return
		}
		ReturnToken(c, w, user.ID, true, user.Data.EmailAddress == "alexander.trakhimenok@gmail.com")
	}
}

func HandleSignInWithEmail(c context.Context, w http.ResponseWriter, r *http.Request) {
	email := strings.TrimSpace(r.PostFormValue("email"))
	password := strings.TrimSpace(r.PostFormValue("password"))
	//logus.Debugf(c, "Email: %s", email)
	if email == "" || password == "" {
		api.ErrorAsJson(c, w, http.StatusBadRequest, errors.New("Missing required value"))
		return
	}

	if err := validateEmail(email); err != nil {
		api.JsonToResponse(c, w, map[string]string{"error": err.Error()})
		return
	}

	userEmail, err := dtdal.UserEmail.GetUserEmailByID(c, nil, email)
	if err != nil {
		if dal.IsNotFound(err) {
			api.ErrorAsJson(c, w, http.StatusForbidden, errors.New("Unknown email"))
		} else {
			api.ErrorAsJson(c, w, http.StatusInternalServerError, err)
		}
		return
	} else if err = userEmail.Data.CheckPassword(password); err != nil {
		logus.Debugf(c, "Invalid password: %v", err)
		api.ErrorAsJson(c, w, http.StatusForbidden, errors.New("invalid password"))
		return
	}

	ReturnToken(c, w, userEmail.Data.AppUserID, false, userEmail.ID == "alexander.trakhimenok@gmail.com")
}

func HandleRequestPasswordReset(c context.Context, w http.ResponseWriter, r *http.Request) {
	email := r.PostFormValue("email")
	userEmail, err := dtdal.UserEmail.GetUserEmailByID(c, nil, email)
	if dal.IsNotFound(err) {
		api.ErrorAsJson(c, w, http.StatusForbidden, errors.New("Unknown email"))
		return
	}

	now := time.Now()

	pwdResetEntity := models.PasswordResetData{
		Email:             userEmail.ID,
		Status:            "created",
		OwnedByUserWithID: appuser.NewOwnedByUserWithID(userEmail.Data.AppUserID, now),
	}

	var db dal.DB
	if db, err = facade.GetDatabase(c); err != nil {
		return
	}
	err = db.RunReadwriteTransaction(c, func(c context.Context, tx dal.ReadwriteTransaction) (err error) {
		_, err = dtdal.PasswordReset.CreatePasswordResetByID(c, tx, &pwdResetEntity)
		return err
	})
	if err != nil {
		api.ErrorAsJson(c, w, http.StatusInternalServerError, err)
		return
	}
}

func HandleChangePasswordAndSignIn(c context.Context, w http.ResponseWriter, r *http.Request) {
	var (
		err           error
		passwordReset models.PasswordReset
	)

	if passwordReset.ID, err = strconv.Atoi(r.PostFormValue("pin")); err != nil {
		api.ErrorAsJson(c, w, http.StatusBadRequest, err)
		return
	}

	pwd := r.PostFormValue("pwd")
	if pwd == "" {
		api.ErrorAsJson(c, w, http.StatusBadRequest, errors.New("Empty password"))
		return
	}

	if passwordReset, err = dtdal.PasswordReset.GetPasswordResetByID(c, nil, passwordReset.ID); err != nil {
		if dal.IsNotFound(err) {
			api.ErrorAsJson(c, w, http.StatusForbidden, errors.New("Unknown pin"))
			return
		}
		api.ErrorAsJson(c, w, http.StatusInternalServerError, err)
		return
	}

	isAdmin := IsAdmin(passwordReset.Data.Email)

	var db dal.DB
	if db, err = facade.GetDatabase(c); err != nil {
		api.ErrorAsJson(c, w, http.StatusInternalServerError, err)
		return
	}
	if err = db.RunReadwriteTransaction(c, func(c context.Context, tx dal.ReadwriteTransaction) error {

		now := time.Now()
		appUser := models.NewAppUser(passwordReset.Data.AppUserID, nil)
		userEmail := models.NewUserEmail(passwordReset.Data.Email, nil)

		records := []dal.Record{appUser.Record, userEmail.Record, passwordReset.Record}

		var db dal.DB
		if db, err = facade.GetDatabase(c); err != nil {
			return err
		}
		if err = db.GetMulti(c, records); err != nil {
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
		api.ErrorAsJson(c, w, http.StatusInternalServerError, err)
		return
	}

	ReturnToken(c, w, passwordReset.Data.AppUserID, false, isAdmin)
}

var errInvalidEmailConformationPin = errors.New("email confirmation pin is not valid")

func HandleConfirmEmailAndSignIn(c context.Context, w http.ResponseWriter, r *http.Request) {
	var (
		err       error
		userEmail models.UserEmailEntry
		pin       string
	)

	userEmail.ID, pin = r.PostFormValue("email"), r.PostFormValue("pin")

	if userEmail.ID == "" {
		api.ErrorAsJson(c, w, http.StatusBadRequest, errors.New("Empty email"))
		return
	}
	if pin == "" {
		api.ErrorAsJson(c, w, http.StatusBadRequest, errors.New("Empty pin"))
		return
	}

	var db dal.DB
	if db, err = facade.GetDatabase(c); err != nil {
		return
	}
	if err = db.RunReadwriteTransaction(c, func(c context.Context, tx dal.ReadwriteTransaction) error {
		now := time.Now()

		if userEmail, err = dtdal.UserEmail.GetUserEmailByID(c, tx, userEmail.ID); err != nil {
			return err
		}

		var appUser models.AppUser
		if appUser, err = facade.User.GetUserByID(c, tx, userEmail.Data.AppUserID); err != nil {
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
			api.ErrorAsJson(c, w, http.StatusBadRequest, err)
			return
		} else if err == errInvalidEmailConformationPin {
			api.ErrorAsJson(c, w, http.StatusForbidden, err)
			return
		}
		api.ErrorAsJson(c, w, http.StatusInternalServerError, err)
		return
	}

	ReturnToken(c, w, userEmail.Data.AppUserID, false, IsAdmin(userEmail.ID))
}
