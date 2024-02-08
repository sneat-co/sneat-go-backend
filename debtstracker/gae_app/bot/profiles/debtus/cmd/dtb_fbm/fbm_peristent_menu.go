package dtb_fbm

//import (
//	"fmt"
//	"net/http"
//	"strings"
//
//	"context"
//	"github.com/sneat-co/debtstracker-translations/emoji"
//	"github.com/strongo/bots-api-fbm"
//	"github.com/bots-go-framework/bots-fw/botsfw"
//	"github.com/strongo/log"
//)
//
//var EM_SPACE = strings.Repeat("\u00A0", 2)
//
//func SetPersistentMenu(c context.Context, r *http.Request, bot botsfw.BotSettings, api fbmbotapi.GraphAPI) (err error) {
//	url := fmt.Sprintf("https://%v/app/#fbm%v", r.Host, bot.ID)
//
//	//menuItemWebUrl := func(icon, title, hash string) fbm_api.MenuItemWebURL {
//	//	return fbm_api.NewMenuItemWebUrl(
//	//		icon + EM_SPACE + title,
//	//		url + hash, fbm_api.WebviewHeightRatioFull, false, true)
//	//}
//	menuItemPostback := func(icon, title, payload string) fbmbotapi.MenuItemPostback {
//		return fbmbotapi.NewMenuItemPostback(icon+EM_SPACE+title, payload)
//	}
//
//	log.Debugf(c, "url: %v", url)
//
//	persistentMenu := func(locale string) fbmbotapi.PersistentMenu {
//
//		//topMenuDebts := fbm_api.NewMenuItemNested(emoji.MEMO_ICON + EM_SPACE + "Debts",
//		//	menuItemWebUrl(emoji.TAKE_ICON, "I borrowed", "#new-debt=contact-to-user"),
//		//	menuItemWebUrl(emoji.GIVE_ICON, "I lent", "#new-debt=user-to-contact"),
//		//	menuItemWebUrl(emoji.RETURN_BACK_ICON, "Returned", "#debt-returned"),
//		//	menuItemWebUrl(emoji.BALANCE_ICON, "Balance", "#debts"),
//		//)
//		//
//		//topMenuBills := fbm_api.NewMenuItemNested(emoji.BILLS_ICON + " Bills",
//		//	menuItemWebUrl(emoji.DIVIDE_ICON, "Split bill", "#split-bill"),
//		//	menuItemWebUrl(emoji.MONEY_BAG_ICON, "Start collecting", "#start-collecting"),
//		//	menuItemWebUrl(emoji.OPEN_BOOK_ICON, "Outstanding bills", "#bills=outstanding"),
//		//	menuItemWebUrl(emoji.CALENDAR_ICON, "Recurring bills", "#bills=recurring"),
//		//)
//		//
//		//topMenuView := fbm_api.NewMenuItemNested(emoji.TOTAL_ICON + EM_SPACE + "More...",
//		//	menuItemPostback(emoji.HOME_ICON, "Get started", "fbm-get-started"),
//		//	menuItemWebUrl(emoji.CONTACTS_ICON, "Contacts", "#contacts"),
//		//	menuItemWebUrl(emoji.HISTORY_ICON, "History", "#history"),
//		//	menuItemWebUrl(emoji.SETTINGS_ICON, "Settings", "#settings"),
//		//)
//		//
//		//return fbm_api.NewPersistentMenu(locale, false,
//		//	topMenuDebts,
//		//	topMenuBills,
//		//	topMenuView,
//		//)
//		return fbmbotapi.NewPersistentMenu(locale, false,
//			menuItemPostback(emoji.HOME_ICON, "Main menu", FbmMainMenuCommand.Code),
//			menuItemPostback(emoji.MEMO_ICON, "Debt", FbmDebtsCommand.Code),
//			menuItemPostback(emoji.BILLS_ICON, "Bills", FbmBillsCommand.Code),
//		)
//	}
//
//	persistentMenuMessage := fbmbotapi.PersistentMenuMessage{
//		PersistentMenus: []fbmbotapi.PersistentMenu{
//			persistentMenu("default"),
//			//persistentMenu("ru_RU"),
//		},
//	}
//
//	if err = api.SetPersistentMenu(c, persistentMenuMessage); err != nil {
//		return
//	}
//	return
//}
