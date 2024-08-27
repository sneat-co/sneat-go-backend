package common4debtus

//
//import (
//	"context"
//	"github.com/bots-go-framework/bots-fw-store/botsfwmodels"
//	"github.com/bots-go-framework/bots-fw/botsfw"
//	"github.com/sneat-co/debtusbot-translations/trans"
//	"github.com/sneat-co/sneat-go-backend/debtusbot/gae_app/debtusbot/models"
//	"github.com/strongo/i18n"
//	"github.com/strongo/strongoapp/appuser"
//	"reflect"
//	"time"
//)
//
//type DebtusAppContext struct {
//}
//
//func (appCtx DebtusAppContext) AppUserCollectionName() string {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (appCtx DebtusAppContext) SetLocale(code5 string) error {
//	//TODO implement me
//	panic("implement me")
//}
//
//var _ botsfw.BotAppContext = (*DebtusAppContext)(nil)
//
//func (appCtx DebtusAppContext) AppUserEntityKind() string {
//	return models.AppUserKind
//}
//
//func (appCtx DebtusAppContext) AppUserEntityType() reflect.ExtraType {
//	return reflect.TypeOf(&models.DebutsAppUserDataOBSOLETE{})
//}
//
//func (appCtx DebtusAppContext) NewBotAppUserEntity() botsfwmodels.AppUserData {
//	return &models.DebutsAppUserDataOBSOLETE{
//		ContactDetails: models.ContactDetails{
//			PhoneContact: models.PhoneContact{},
//		},
//		DtCreated: time.Now(),
//	}
//}
//
//func (appCtx DebtusAppContext) GetBotChatEntityFactory(platform string) func() botsfwmodels.BotChatData {
//	switch platform {
//	case "telegram":
//		panic("not implemented")
//		//return func() botsfwmodels.ChatBaseData {
//		//	return &models.DebtusTelegramChatData{
//		//		TgChatBase: *botsfwtgmodels.NewTelegramChatEntity(),
//		//	}
//		//}
//	default:
//		panic("Unknown platform: " + platform)
//	}
//}
//
//func (appCtx DebtusAppContext) NewAppUserData() appuser.BaseUserData {
//	return appCtx.NewBotAppUserEntity()
//}
//
//func (appCtx DebtusAppContext) GetTranslator(ctx context.Context) i18n.Translator {
//	return i18n.NewMapTranslator(ctx, trans.TRANS)
//}
//
//func (appCtx DebtusAppContext) SupportedLocales() i18n.LocalesProvider {
//	return trans.DebtsTrackerLocales{}
//}
//
//var _ botsfw.BotAppContext = (*DebtusAppContext)(nil)
//
//var TheAppContext = DebtusAppContext{}
