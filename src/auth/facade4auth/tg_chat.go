package facade4auth

import (
	"fmt"
	"github.com/bots-go-framework/bots-fw-telegram-models/botsfwtgmodels"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	"github.com/sneat-co/sneat-go-backend/src/auth"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/models4debtus"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/dbo4userus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/logus"
	"strconv"
	"strings"
	"sync"

	"context"
)

type TgChatDalGae struct {
}

func NewTgChatDalGae() TgChatDalGae {
	return TgChatDalGae{}
}

func (TgChatDalGae) GetTgChatByID(c context.Context, tgBotID string, tgChatID int64) (tgChat models4debtus.DebtusTelegramChat, err error) {
	tgChatFullID := fmt.Sprintf("%s:%d", tgBotID, tgChatID)
	key := dal.NewKeyWithID(botsfwtgmodels.TgChatCollection, tgChatFullID)
	data := new(models4debtus.DebtusTelegramChatData)
	tgChat = models4debtus.DebtusTelegramChat{
		WithID: record.NewWithID(tgChatFullID, key, data),
		Data:   data,
	}
	//tgChat.SetID(tgBotID, tgChatID)

	var db dal.DB
	if db, err = facade.GetDatabase(c); err != nil {
		return
	}
	err = db.Get(c, tgChat.Record)
	return
}

func (TgChatDalGae) /* TODO: rename properly! */ DoSomething(c context.Context,
	userTask *sync.WaitGroup, currency string, tgChatID int64, authInfo auth.AuthInfo, user dbo4userus.UserEntry,
	sendToTelegram func(tgChat botsfwtgmodels.TgChatData) error,
) (err error) {
	var isSentToTelegram bool // Needed in case of failed to save to DB and is auto-retry
	debtusTgChatData := &models4debtus.DebtusTelegramChatData{}

	id := strconv.FormatInt(tgChatID, 10)
	debtusTgChat := models4debtus.DebtusTelegramChat{
		WithID: record.NewWithID(id, dal.NewKeyWithID(botsfwtgmodels.TgChatCollection, id), debtusTgChatData),
		Data:   debtusTgChatData,
	}

	if err = facade.RunReadwriteTransaction(c, func(tc context.Context, tx dal.ReadwriteTransaction) (err error) {
		if err = tx.Get(tc, debtusTgChat.Record); err != nil {
			return err
		}
		if debtusTgChat.Data.BotID == "" {
			logus.Errorf(c, "Data inconsistency issue - TgChat(%v).BotID is empty string", tgChatID)
			if strings.Contains(authInfo.Issuer, ":") {
				issuer := strings.Split(authInfo.Issuer, ":")
				if strings.ToLower(issuer[0]) == "telegram" {
					debtusTgChat.Data.BotID = issuer[1]
					logus.Infof(c, "Data inconsistency fixed, set to: %v", debtusTgChat.Data.BotID)
				}
			}
		}
		debtusTgChat.Data.AddWizardParam("currency", string(currency))

		if !isSentToTelegram {
			if err = sendToTelegram(debtusTgChat.Data); err != nil { // This is some serious architecture sheet. Too sleepy to make it right, just make it working.
				return err
			}
			isSentToTelegram = true
		}
		if err = tx.Set(tc, debtusTgChat.Record); err != nil {
			return fmt.Errorf("failed to save Telegram chat record to db: %w", err)
		}
		return err
	}, nil); err != nil {
		err = fmt.Errorf("method TgChatDalGae.DoSomething() transaction failed: %w", err)
	}
	return
}
