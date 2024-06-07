package splitus

import (
	"bytes"
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/crediterra/money"
	"github.com/sneat-co/debtstracker-translations/emoji"
	"github.com/sneat-co/debtstracker-translations/trans"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/bot/profiles/shared_all"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/bot/profiles/shared_group"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/models"
	"net/url"
)

func GroupSettingsAction(whc botsfw.WebhookContext, group models.GroupEntry, isEdit bool) (m botsfw.MessageFromBot, err error) {
	var buf bytes.Buffer
	buf.WriteString(whc.Translate(trans.MT_GROUP_LABEL, group.Data.Name))
	buf.WriteString("\n")
	buf.WriteString(whc.Translate(trans.MT_TEXT_MEMBERS_COUNT, group.Data.MembersCount))
	m.Format = botsfw.MessageFormatHTML
	m.Text = buf.String()
	defaultCurrency := group.Data.DefaultCurrency
	if defaultCurrency == "" {
		defaultCurrency = money.CurrencyCode(whc.Translate(trans.NOT_SET))
	}
	m.Keyboard = tgbotapi.NewInlineKeyboardMarkup(
		[]tgbotapi.InlineKeyboardButton{
			{
				Text:         whc.Translate(trans.BUTTON_TEXT_MANAGE_MEMBERS),
				CallbackData: GroupMembersCommandCode + "?group=" + group.ID,
			},
		},
		[]tgbotapi.InlineKeyboardButton{
			{
				Text:         whc.Translate(trans.BT_DEFAULT_CURRENCY, defaultCurrency),
				CallbackData: GroupSettingsChooseCurrencyCommandCode,
			},
		},
		[]tgbotapi.InlineKeyboardButton{
			{
				Text:         whc.Translate(trans.BUTTON_TEXT_SPLIT_MODE, whc.Translate(string(group.Data.GetSplitMode()))),
				CallbackData: shared_group.GroupCallbackCommandData(groupSplitCommandCode, group.ID),
			},
		},
		[]tgbotapi.InlineKeyboardButton{
			tgbotapi.NewInlineKeyboardButtonSwitchInlineQueryCurrentChat(
				emoji.CLIPBOARD_ICON+whc.Translate(trans.COMMAND_TEXT_NEW_BILL),
				"",
			),
		},
	)
	m.IsEdit = isEdit
	return
}

var settingsCommand = func() (settingsCommand botsfw.Command) {
	settingsCommand = shared_all.SettingsCommandTemplate
	settingsCommand.Action = settingsAction
	settingsCommand.CallbackAction = func(whc botsfw.WebhookContext, callbackUrl *url.URL) (m botsfw.MessageFromBot, err error) {
		m, err = settingsAction(whc)
		m.IsEdit = true
		return
	}
	return
}()

func settingsAction(whc botsfw.WebhookContext) (m botsfw.MessageFromBot, err error) {
	if whc.IsInGroup() {
		groupAction := shared_group.NewGroupAction(func(whc botsfw.WebhookContext, group models.GroupEntry) (m botsfw.MessageFromBot, err error) {
			return GroupSettingsAction(whc, group, false)
		})
		return groupAction(whc)
	} else {
		m, _, err = shared_all.SettingsMainTelegram(whc)
		return
	}
}
