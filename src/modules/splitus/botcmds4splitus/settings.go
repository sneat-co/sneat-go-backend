package botcmds4splitus

import (
	"bytes"
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/crediterra/money"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/debtstracker-translations/emoji"
	"github.com/sneat-co/debtstracker-translations/trans"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/const4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/dal4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/debtusbots/profiles/shared_all"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/debtusbots/profiles/shared_space"
	"github.com/sneat-co/sneat-go-backend/src/modules/spaceus/dbo4spaceus"
	"github.com/sneat-co/sneat-go-backend/src/modules/splitus/models4splitus"
	"net/url"
)

func SpaceSettingsAction(whc botsfw.WebhookContext, space dbo4spaceus.SpaceEntry, isEdit bool) (m botsfw.MessageFromBot, err error) {

	splitusSpace := models4splitus.NewSplitusSpaceEntry(space.ID)
	contactusSpace := dal4contactus.NewContactusSpaceEntry(space.ID)

	db := whc.DB()
	if err = db.GetMulti(whc.Context(), []dal.Record{splitusSpace.Record, contactusSpace.Record}); err != nil {
		return
	}

	var membersCount int
	if contactusSpace.Record.Exists() {
		membersCount = contactusSpace.Data.GetContactsCount(const4contactus.SpaceMemberRoleMember)
	}
	var buf bytes.Buffer
	buf.WriteString(whc.Translate(trans.MT_GROUP_LABEL, space.Data.Title))
	buf.WriteString("\n")
	buf.WriteString(whc.Translate(trans.MT_TEXT_MEMBERS_COUNT, membersCount))
	m.Format = botsfw.MessageFormatHTML
	m.Text = buf.String()
	defaultCurrency := space.Data.PrimaryCurrency
	if defaultCurrency == "" {
		defaultCurrency = money.CurrencyCode(whc.Translate(trans.NOT_SET))
	}
	m.Keyboard = tgbotapi.NewInlineKeyboardMarkup(
		[]tgbotapi.InlineKeyboardButton{
			{
				Text:         whc.Translate(trans.BUTTON_TEXT_MANAGE_MEMBERS),
				CallbackData: GroupMembersCommandCode + "?space=" + space.ID,
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
				Text:         whc.Translate(trans.BUTTON_TEXT_SPLIT_MODE, whc.Translate(string(splitusSpace.Data.GetSplitMode()))),
				CallbackData: shared_space.SpaceCallbackCommandData(spaceSplitCommandCode, space.ID),
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
	var isInGroup bool
	if isInGroup, err = whc.IsInGroup(); err != nil {
		return
	} else if isInGroup {
		groupAction := shared_space.NewSpaceAction(func(whc botsfw.WebhookContext, space dbo4spaceus.SpaceEntry) (m botsfw.MessageFromBot, err error) {
			return SpaceSettingsAction(whc, space, false)
		})
		return groupAction(whc)
	} else {
		m, _, err = shared_all.SettingsMainTelegram(whc)
		return
	}
}
