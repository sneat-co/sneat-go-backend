package sneatbot

import (
	"context"
	"fmt"
	"github.com/bots-go-framework/bots-fw-store/botsfwmodels"
	"github.com/bots-go-framework/bots-fw/botsdal"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	"github.com/sneat-co/debtstracker-translations/trans"
	"github.com/sneat-co/sneat-go-backend/src/modules/auth/facade4auth"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/strongo/i18n"
	"github.com/strongo/strongoapp/appuser"
	"reflect"
)

var _ botsfw.AppContext = (*sneatAppBotContext)(nil)

func NewSneatAppContextForBotsFW() botsfw.AppContext {
	return sneatAppBotContext{
		LocalesProvider: i18n.NewSupportedLocales([]string{
			i18n.LocaleCodeEnUK,
			i18n.LocaleCodeRuRU,
			i18n.LocaleCodeUaUA,
			i18n.LocaleCodeEsES,
			i18n.LocaleCodePtPT,
			i18n.LocaleCodeFrFR,
			i18n.LocaleCodeItIT,
			i18n.LocaleCodeDeDE,
			i18n.LocaleCodeFaIR,
			i18n.LocaleCodeZhCN,
			i18n.LocaleCodeJaJP,
			i18n.LocaleCodeKoKO,
			i18n.LocaleCodePtPT,
		}),
	}
}

var _ botsdal.AppUserDal = (*sneatAppBotDal)(nil)

type sneatAppBotDal struct {
}

func (s sneatAppBotDal) CreateAppUserFromBotUser(ctx context.Context,
	tx dal.ReadwriteTransaction, // TODO(fix): intentionally not using transaction as we can't have reads after writes
	bot botsdal.Bot,
) (
	appUser record.DataWithID[string, botsfwmodels.AppUserData],
	botUser record.DataWithID[string, botsfwmodels.PlatformUserData],
	err error,
) {
	botUserData := facade4auth.BotUserData{
		PlatformID:   bot.Platform,
		BotID:        bot.ID,
		BotUserID:    fmt.Sprintf("%v", bot.User.GetID()),
		FirstName:    bot.User.GetFirstName(),
		LastName:     bot.User.GetLastName(),
		Username:     bot.User.GetUserName(),
		LanguageCode: bot.User.GetLanguage(),
	}
	remoteClientInfo := dbmodels.RemoteClientInfo{
		HostOrApp: string(bot.Platform) + "@" + bot.ID,
	}
	var params facade4auth.CreateUserWorkerParams

	botUser, params, err = facade4auth.CreateBotUserAndAppUserRecords(ctx, tx, appUser.ID, bot.Platform, botUserData, remoteClientInfo)
	if err != nil {
		err = fmt.Errorf("failed to create user records: %w", err)
		return
	}
	appUser = record.DataWithID[string, botsfwmodels.AppUserData]{
		WithID: params.User.WithID,
		Data:   params.User.Data,
	}

	{ // Insert records except bot and app user records that botsfw will handle
		keysToExclude := make([]*dal.Key, 0, 2)
		for _, r := range params.RecordsToInsert() {
			// For some reason == does not work for botUser.Key
			if key := r.Key(); key.Equal(appUser.Key) || key.Equal(botUser.Key) {
				keysToExclude = append(keysToExclude, key)
			}
		}
		if err = params.ApplyChanges(ctx, tx, keysToExclude...); err != nil {
			err = fmt.Errorf("failed to apply changes: %w", err)
			return
		}
	}
	return
}

type sneatAppBotContext struct { // TODO: Duplication?!
	sneatAppBotDal
	i18n.LocalesProvider
}

//	func (s sneatAppBotContext) GetAppUserByBotUserID(ctx context.Context, tx dal.ReadwriteTransaction, platform, botID, botUserID string) (appUser record.DataWithID[string, botsfwmodels.AppUserData], err error) {
//		//TODO implement me
//		panic("implement me")
//	}
//
//	func (s sneatAppBotContext) AppUserCollectionName() string {
//		return "Users"
//	}
//
//	func (s sneatAppBotContext) AppUserEntityKind() string {
//		return "User"
//	}
func (s sneatAppBotContext) SetLocale(code5 string) error {
	panic(fmt.Sprintf("TODO: why we have this? should be removed?: code5=%s", code5))
}

func (s sneatAppBotContext) AppUserEntityType() reflect.Type {
	//TODO implement me
	panic("implement AppUserEntityType()")
}

func (s sneatAppBotContext) NewAppUserData() appuser.BaseUserData {
	//TODO implement me
	panic("implement NewAppUserData()")
}

func (s sneatAppBotContext) GetTranslator(ctx context.Context) i18n.Translator {
	return i18n.NewMapTranslator(ctx, i18n.LocaleCodeEnUK, trans.TRANS)
}

func (s sneatAppBotContext) NewBotAppUserEntity() botsfwmodels.AppUserData {
	//TODO implement me
	panic("implement NewBotAppUserEntity()")
}

func (s sneatAppBotContext) GetBotChatEntityFactory(platform string) func() botsfwmodels.BotChatData {
	//TODO implement me
	panic(fmt.Sprintf("implement me. platform=%s", platform))
}
