package unsorted

import (
	"github.com/bots-go-framework/bots-fw-store/botsfwmodels"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/sneat-co/sneat-go-backend/src/coremodules/userus/dbo4userus"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/gae_app/debtstracker/dtdal"
	"github.com/strongo/logus"
	"net/http"

	"github.com/bots-go-framework/bots-fw-telegram"
)

type ApiWebhookContext struct {
	appUser    *dbo4userus.UserDbo
	appUserID  string
	botChatID  int64
	chatEntity botsfwmodels.BotChatData
	*botsfw.WebhookContextBase
}

var _ botsfw.WebhookContext = (*ApiWebhookContext)(nil)

func (ApiWebhookContext) IsInGroup() (bool, error) {
	panic("not supported")
}

func NewApiWebhookContext(r *http.Request, appUser *dbo4userus.UserDbo, userID string, botChatID int64, chatData botsfwmodels.BotChatData) ApiWebhookContext {
	var botSettings botsfw.BotSettings
	botContext := botsfw.NewBotContext(dtdal.BotHost, &botSettings)
	args := botsfw.NewCreateWebhookContextArgs(
		r,
		nil, /*anybot.TheAppContext*/
		*botContext,
		nil,
		nil,
		nil,
	)
	whcb, err := botsfw.NewWebhookContextBase(
		args,
		telegram.Platform, // webhookInput
		nil,               // records fields setter
		func() (bool, error) { return false, nil },
		nil, // GaMeasurement
	)
	if err != nil {
		logus.Errorf(r.Context(), "failed to create WebhookContextBase: %v", err)
	}
	whc := ApiWebhookContext{
		appUser:            appUser,
		appUserID:          userID,
		botChatID:          botChatID,
		chatEntity:         chatData,
		WebhookContextBase: whcb,
	}
	if err := whc.SetLocale(chatData.GetPreferredLanguage()); err != nil {
		logus.Errorf(r.Context(), "failed to set locale: %v", err)
	}
	return whc
}

func (whc ApiWebhookContext) AppUserData() (botsfwmodels.AppUserData, error) {
	//TODO implement me
	panic("implement me")
}

func (whc ApiWebhookContext) BotChatIntID() int64 {
	return whc.botChatID
}

func (whc ApiWebhookContext) ChatEntity() botsfwmodels.BotChatData {
	return whc.chatEntity
}

func (whc ApiWebhookContext) GetAppUser() (botsfwmodels.AppUserData, error) {
	panic("implement me")
	//return nil /*whc.appUser*/, nil
}

func (whc ApiWebhookContext) Init(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (whc ApiWebhookContext) IsNewerThen(chatEntity botsfwmodels.BotChatData) bool {
	return true
}

func (whc ApiWebhookContext) MessageText() string {
	return ""
}

func (whc ApiWebhookContext) NewEditMessage(text string, format botsfw.MessageFormat) (m botsfw.MessageFromBot, err error) {
	panic("Not implemented")
}

func (whc ApiWebhookContext) Responder() botsfw.WebhookResponder {
	panic("Not implemented")
}

func (whc ApiWebhookContext) UpdateLastProcessed(chatEntity botsfwmodels.BotChatData) error {
	panic("Not implemented")
}
