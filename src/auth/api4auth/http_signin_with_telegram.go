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
	"github.com/sneat-co/sneat-go-backend/src/bots/models4bots"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/gae_app/debtstracker/api4debtus"
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
	if tgBotUser, err = botsdal.GetPlatformUser(ctx, db, telegram.PlatformID, botUserID, new(models4bots.TelegramUserDbo)); err != nil {
		if !dal.IsNotFound(err) {
			apicore.ReturnError(ctx, w, r, err)
		}
		if tgBotUser, err = createTelegramUserAndAppUserRecords(ctx, db, initData); err != nil {
			return
		}
		isNewUser = true
	}

	appUserID := tgBotUser.Data.GetAppUserID()

	api4debtus.ReturnToken(ctx, w, appUserID, isNewUser, false)
}

func createTelegramUserAndAppUserRecords(ctx context.Context, db dal.DB, initData twainitdata.InitData) (tgBotUser record.DataWithID[string, botsfwmodels.PlatformUserData], err error) {
	telegramUserData := new(models4bots.TelegramUserDbo)
	telegramUserData.FirstName = initData.User.FirstName
	telegramUserData.LastName = initData.User.LastName
	telegramUserData.UserName = initData.User.Username

	tgBotUser.ID = strconv.FormatInt(initData.User.ID, 10)
	key := botsdal.NewPlatformUserKey(telegram.PlatformID, tgBotUser.ID)
	tgBotUser.Record = dal.NewRecordWithData(key, telegramUserData)

	if err = db.RunReadwriteTransaction(ctx, func(ctx context.Context, tx dal.ReadwriteTransaction) error {
		if err = tx.Insert(ctx, tgBotUser.Record); err != nil {
			err = fmt.Errorf("failed to insert telegram user record: %w", err)
			return err
		}
		return nil
	}); err != nil {
		return
	}
	return
}
