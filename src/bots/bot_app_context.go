package bots

import (
	"context"
	"fmt"
	"github.com/bots-go-framework/bots-fw-store/botsfwmodels"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/strongo/i18n"
	"github.com/strongo/strongoapp/appuser"
	"reflect"
)

var _ botsfw.BotAppContext = (*sneatAppBotContext)(nil)

type sneatAppBotContext struct {
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
	panic("implement me")
}

func (s sneatAppBotContext) NewAppUserData() appuser.BaseUserData {
	//TODO implement me
	panic("implement me")
}

func (s sneatAppBotContext) GetTranslator(c context.Context) i18n.Translator {
	return i18n.NewMapTranslator(c, nil)
}

func (s sneatAppBotContext) SupportedLocales() i18n.LocalesProvider {
	return nil
}

func (s sneatAppBotContext) NewBotAppUserEntity() botsfwmodels.AppUserData {
	//TODO implement me
	panic("implement me")
}

func (s sneatAppBotContext) GetBotChatEntityFactory(platform string) func() botsfwmodels.BotChatData {
	//TODO implement me
	panic(fmt.Sprintf("implement me. platform=%s", platform))
}
