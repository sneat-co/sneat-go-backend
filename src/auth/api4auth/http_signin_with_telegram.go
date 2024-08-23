package api4auth

import (
	"context"
	"fmt"
	"github.com/bots-go-framework/bots-fw-store/botsfwmodels"
	telegram "github.com/bots-go-framework/bots-fw-telegram"
	"github.com/bots-go-framework/bots-fw-telegram-webapp/twainitdata"
	"github.com/bots-go-framework/bots-fw/botsfw/botsdal"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	"github.com/sneat-co/sneat-go-backend/src/auth/facade4auth"
	"github.com/sneat-co/sneat-go-backend/src/bots/models4bots"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/gae_app/debtstracker/api4debtus"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/dbo4userus"
	"github.com/sneat-co/sneat-go-core/apicore"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/validation"
	"io"
	"net/http"
	"strconv"
	"time"
)

func httpSignInFromTelegramMiniapp(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	// Read request body to a string
	bodyBytes, err := io.ReadAll(r.Body)

	if err != nil {
		err = fmt.Errorf("failed to read request body: %w", err)
		apicore.ReturnError(ctx, w, r, err)
		return
	}

	const botToken = "TODO: get bot token from config or env variable"

	initDataStr := string(bodyBytes)
	if err = twainitdata.Validate(initDataStr, botToken, 10*time.Second); err != nil {
		err = validation.NewBadRequestError(err)
		apicore.ReturnError(ctx, w, r, err)
		return
	}

	var initData twainitdata.InitData
	if initData, err = twainitdata.Parse(initDataStr); err != nil {
		err = validation.NewBadRequestError(err)
		apicore.ReturnError(ctx, w, r, err)
	}

	botUserID := fmt.Sprintf("TGU%d", initData.User.ID)

	var db dal.DB
	if db, err = facade.GetDatabase(ctx); err != nil {
		apicore.ReturnError(ctx, w, r, err)
		return
	}

	var isNewUser bool
	var tgBotUser record.DataWithID[string, botsfwmodels.PlatformUserData]
	err = db.RunReadwriteTransaction(ctx, func(ctx context.Context, tx dal.ReadwriteTransaction) (err error) {
		if tgBotUser, err = botsdal.GetPlatformUser(ctx, db, telegram.PlatformID, botUserID, new(models4bots.TelegramUserDbo)); err != nil {
			if !dal.IsNotFound(err) {
				apicore.ReturnError(ctx, w, r, err)
			}
			if tgBotUser, err = createTelegramUserAndAppUserRecords(ctx, tx, initData); err != nil {
				return
			}
			isNewUser = true
		}

		if appUserID := tgBotUser.Data.GetAppUserID(); appUserID == "" {
			if err = createAppUserRecordAndUpdateTelegramUserRecord(ctx, tx, tgBotUser, initData); err != nil {
				err = fmt.Errorf("failed in createAppUserRecordAndUpdateTelegramUserRecord(): %w", err)
				return
			}
		}
		return
	})
	if err != nil {
		apicore.ReturnError(ctx, w, r, err)
		return
	}

	api4debtus.ReturnToken(ctx, w, tgBotUser.Data.GetAppUserID(), isNewUser, false)
}

func createAppUserRecordAndUpdateTelegramUserRecord(ctx context.Context, tx dal.ReadwriteTransaction, tgBotUser record.DataWithID[string, botsfwmodels.PlatformUserData], initData twainitdata.InitData) (err error) {
	var user dbo4userus.UserEntry
	if user, err = createUserRecordWithRandomID(ctx, tx); err != nil {
		return
	}
	tgBotUser.Data.SetAppUserID(user.ID)
	tgUserDbo := tgBotUser.Data.(*models4bots.TelegramUserDbo)
	tgUserDbo.FirstName = initData.User.FirstName
	tgUserDbo.LastName = initData.User.LastName
	tgUserDbo.UserName = initData.User.Username

	if err = tx.Set(ctx, tgBotUser.Record); err != nil { // TODO: Implement update
		err = fmt.Errorf("failed to update telegram user record: %w", err)
		return
	}
	return
}

func createUserRecordWithRandomID(ctx context.Context, tx dal.ReadwriteTransaction) (user dbo4userus.UserEntry, err error) {
	if user.ID, err = facade4auth.GenerateRandomUserID(ctx, tx); err != nil {
		return
	}
	user = dbo4userus.NewUserEntry(user.ID)
	if err = tx.Insert(ctx, user.Record); err != nil {
		err = fmt.Errorf("failed to insert user record: %w", err)
	}
	return
}

func createTelegramUserAndAppUserRecords(ctx context.Context, tx dal.ReadwriteTransaction, initData twainitdata.InitData) (tgBotUser record.DataWithID[string, botsfwmodels.PlatformUserData], err error) {
	telegramUserData := new(models4bots.TelegramUserDbo)
	telegramUserData.FirstName = initData.User.FirstName
	telegramUserData.LastName = initData.User.LastName
	telegramUserData.UserName = initData.User.Username

	tgBotUser.ID = strconv.FormatInt(initData.User.ID, 10)
	key := botsdal.NewPlatformUserKey(telegram.PlatformID, tgBotUser.ID)
	tgBotUser.Record = dal.NewRecordWithData(key, telegramUserData)

	var user dbo4userus.UserEntry
	if user, err = createUserRecordWithRandomID(ctx, tx); err != nil {
		return
	}
	telegramUserData.SetAppUserID(user.ID)
	if err = tx.Insert(ctx, tgBotUser.Record); err != nil {
		err = fmt.Errorf("failed to insert telegram user record: %w", err)
		return
	}
	return
}
