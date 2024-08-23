package sneatbot

import (
	"context"
	"errors"
	"fmt"
	"github.com/bots-go-framework/bots-fw-store/botsfwmodels"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/dal-go/dalgo/record"
	"github.com/sneat-co/debtstracker-translations/trans"
	"github.com/strongo/i18n"
	"github.com/strongo/strongoapp/appuser"
	"reflect"
)

var _ botsfw.BotAppContext = (*sneatAppBotContext)(nil)

func NewSneatAppBotContext() botsfw.BotAppContext {
	return sneatAppBotContext{
		LocalesProvider: i18n.NewSupportedLocales([]string{i18n.LocaleCodeEnUS, i18n.LocalCodeRuRu}),
	}
}

type sneatAppBotContext struct { // TODO: Duplication?!
	i18n.LocalesProvider
}

func (s sneatAppBotContext) GetAppUserByBotUserID(ctx context.Context, platform, botID, botUserID string) (appUser record.DataWithID[string, botsfwmodels.AppUserData], err error) {
	err = errors.New("GetAppUserByBotUserID() is not implemented in sneatAppBotContext")
	return
}

func (s sneatAppBotContext) AppUserCollectionName() string {
	return "Users"
}

func (s sneatAppBotContext) SetLocale(code5 string) error {
	panic(fmt.Sprintf("TODO: why we have this? should be removed?: code5=%s", code5))
}

func (s sneatAppBotContext) AppUserEntityKind() string {
	return "User"
}

func (s sneatAppBotContext) AppUserEntityType() reflect.Type {
	//TODO implement me
	panic("implement AppUserEntityType()")
}

func (s sneatAppBotContext) NewAppUserData() appuser.BaseUserData {
	//TODO implement me
	panic("implement NewAppUserData()")
}

func (s sneatAppBotContext) GetTranslator(c context.Context) i18n.Translator {
	return i18n.NewMapTranslator(c, trans.TRANS)
}

func (s sneatAppBotContext) NewBotAppUserEntity() botsfwmodels.AppUserData {
	//TODO implement me
	panic("implement NewBotAppUserEntity()")
}

func (s sneatAppBotContext) GetBotChatEntityFactory(platform string) func() botsfwmodels.BotChatData {
	//TODO implement me
	panic(fmt.Sprintf("implement me. platform=%s", platform))
}
