package dtb_settings

import (
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/sneat-co/debtstracker-translations/trans"
	"net/url"

	"github.com/sneat-co/debtstracker-translations/emoji"
)

const CONTACTS_LIST_COMMAND = "contacts-list"

var ContactsListCommand = botsfw.Command{
	Code:     CONTACTS_LIST_COMMAND,
	Commands: trans.Commands(trans.COMMAND_TEXT_CONTACTS, emoji.MAN_AND_WOMAN),
	Action:   contactsAction,
	CallbackAction: func(whc botsfw.WebhookContext, _ *url.URL) (m botsfw.MessageFromBot, err error) {
		return contactsAction(whc)
	},
}

func contactsAction(whc botsfw.WebhookContext) (m botsfw.MessageFromBot, err error) {
	panic("implement me")
	//var user *models.DebutsAppUserDataOBSOLETE
	//if appUser, err := whc.AppUserData(); err != nil {
	//	return m, err
	//} else {
	//	user = appUser.(*models.DebutsAppUserDataOBSOLETE)
	//}
	//var buffer bytes.Buffer
	//buffer.WriteString(fmt.Sprintf("<b>%v</b>\n", whc.Translate(trans.COMMAND_TEXT_CONTACTS)))
	//linker := common.NewLinkerFromWhc(whc)
	//contacts := user.Contacts()
	//numFormat := "%0" + strconv.Itoa(len(strconv.Itoa(len(contacts)))) + "d. "
	//if len(contacts) == 0 {
	//	buffer.WriteString(whc.Translate(trans.MESSAGE_TEXT_YOU_HAVE_NO_CONTACTS))
	//} else {
	//	for i, contact := range contacts {
	//		buffer.WriteString(fmt.Sprintf(numFormat, i+1))
	//		buffer.WriteString(fmt.Sprintf(`<a href="%v">%v</a>`, linker.UrlToContact(contact.ID), html.EscapeString(contact.Name)))
	//		if contact.Status != "" && contact.Status != models.STATUS_ACTIVE {
	//			buffer.WriteString(" (")
	//			buffer.WriteString(contact.Status)
	//			buffer.WriteString(")")
	//		}
	//		buffer.WriteString("\n")
	//	}
	//}
	//keyboard := tgbotapi.NewInlineKeyboardMarkup(
	//	[]tgbotapi.InlineKeyboardButton{
	//		{
	//			Text:         whc.CommandText(trans.COMMAND_TEXT_REFRESH, emoji.REFRESH_ICON),
	//			CallbackData: CONTACTS_LIST_COMMAND + "?do=refresh",
	//		},
	//	},
	//)
	//buffer.WriteString(fmt.Sprintf("\n\nRefreshed on: %v", time.Now()))
	//m = whc.NewMessage(buffer.String())
	//m.Keyboard = keyboard
	//m.IsEdit = whc.InputType() == botsfw.WebhookInputCallbackQuery
	////if callbackUrl.Query().Get("do") == "refresh" {
	////	if m, err = bot.SendRefreshOrNothingChanged(whc, m); err != nil {
	////		return
	////	}
	////}
	//return
}

//const CONTACT_DETAILS_COMMAND = "contact-details"
//
//var ContactDetailsCommand = botsfw.Command{
//	Code:     CONTACTS_LIST_COMMAND,
//	Commands: trans.Commands(CONTACTS_LIST_COMMAND),
//	CallbackAction: func(whc botsfw.WebhookContext, _ *url.URL) (m botsfw.MessageFromBot, err error) {
//		keyboard := tgbotapi.NewInlineKeyboardMarkup(
//			[]tgbotapi.InlineKeyboardButton{
//				{
//					Text:         whc.CommandText(trans.COMMAND_TEXT_LANGUAGE, emoji.EARTH_ICON),
//					CallbackData: SETTINGS_LOCALE_LIST_CALLBACK_PATH,
//				},
//			},
//		)
//		messageText := whc.NewMessageByCode(trans.MESSAGE_TEXT_CONTACT_DETAILS)
//		m.TelegramEditMessageText = telegram.EditMessageOnCallbackQuery(whc.Input().(botsfw.WebhookCallbackQuery), "HTML", messageText)
//		m.TelegramEditMessageText.ReplyMarkup = keyboard
//		return
//	},
//}
//
//const DELETE_CONTACT_COMMAND = "delete-contact"
//
//var DeleteContactCommand = botsfw.Command{
//	Code:     DELETE_CONTACT_COMMAND,
//	Commands: trans.Commands(CONTACTS_LIST_COMMAND),
//	CallbackAction: func(whc botsfw.WebhookContext, _ *url.URL) (m botsfw.MessageFromBot, err error) {
//
//		return
//	},
//}
//
//const EDIT_CONTACT_NAME_COMMAND = "edit-contact-name"
//
//var EditContactNameCommand = botsfw.Command{
//	Code:     EDIT_CONTACT_NAME_COMMAND,
//	Commands: trans.Commands(CONTACTS_LIST_COMMAND),
//	CallbackAction: func(whc botsfw.WebhookContext, _ *url.URL) (m botsfw.MessageFromBot, err error) {
//
//		return
//	},
//}
