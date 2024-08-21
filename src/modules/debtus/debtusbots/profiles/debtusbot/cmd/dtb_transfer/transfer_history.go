package dtb_transfer

import (
	"bytes"
	"fmt"
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/sneat-co/debtstracker-translations/trans"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/common4debtus"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/debtusbots/profiles/debtusbot/cmd/dtb_general"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/gae_app/debtstracker/dtdal"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/models4debtus"
	"github.com/strongo/i18n"
	"net/url"
	"strings"
	"time"

	"github.com/sneat-co/debtstracker-translations/emoji"
	"github.com/yaa110/go-persian-calendar"
)

const HistoryTopLimit = 5
const HistoryMoreLimit = 10

const HISTORY_COMMAND = "history"

var HistoryCommand = botsfw.Command{
	Code:     HISTORY_COMMAND,
	Icon:     emoji.HISTORY_ICON,
	Title:    trans.COMMAND_TEXT_HISTORY,
	Commands: trans.Commands(trans.COMMAND_HISTORY, emoji.HISTORY_ICON), // TODO: Check icon!
	Titles:   map[string]string{botsfw.ShortTitle: emoji.HISTORY_ICON},  // TODO: Check icon!
	Action: func(whc botsfw.WebhookContext) (m botsfw.MessageFromBot, err error) {
		return showHistoryCard(whc, HistoryTopLimit)
	},
}

func showHistoryCard(whc botsfw.WebhookContext, limit int) (m botsfw.MessageFromBot, err error) {
	c := whc.Context()

	var transfers []models4debtus.TransferEntry
	var hasMore bool

	if appUserID := whc.AppUserID(); appUserID != "" {
		if transfers, hasMore, err = dtdal.Transfer.LoadTransfersByUserID(c, appUserID, 0, limit); err != nil {
			return m, err
		}
	}

	if len(transfers) == 0 {
		m = whc.NewMessage(whc.Translate(trans.MESSAGE_TEXT_HISTORY_NO_RECORDS) + common4debtus.HORIZONTAL_LINE + dtb_general.AdSlot(whc, UTM_CAMPAIGN_TRANSFER_HISTORY))
	} else {
		m = whc.NewMessage(whc.Translate(
			trans.MESSAGE_TEXT_HISTORY_LIST,
			whc.Translate(trans.MESSAGE_TEXT_HISTORY_HEADER),
			len(transfers),
			transferHistoryRows(whc, transfers),
		) + common4debtus.HORIZONTAL_LINE + dtb_general.AdSlot(whc, UTM_CAMPAIGN_TRANSFER_HISTORY))
		if hasMore {
			//api4transfers = api4transfers[:limit]
			utmParams := common4debtus.FillUtmParams(whc, common4debtus.UtmParams{Campaign: UTM_CAMPAIGN_TRANSFER_HISTORY})
			m.Keyboard = &tgbotapi.InlineKeyboardMarkup{
				InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{
					{
						tgbotapi.NewInlineKeyboardButtonURL(
							whc.Translate(trans.INLINE_BUTTON_SHOW_FULL_HISTORY),
							//fmt.Sprintf("transfer-history?offset=%v", len(api4transfers)),
							fmt.Sprintf("https://debtstracker.io/%v/history?user=%v#%v", whc.Locale().SiteCode(), whc.AppUserID(), utmParams),
						),
					},
				},
			}
		}
	}
	m.Format = botsfw.MessageFormatHTML
	m.DisableWebPagePreview = true
	return m, nil
}

const (
	UTM_CAMPAIGN_TRANSFER_HISTORY = "transfer-history"
)

func transferHistoryRows(whc botsfw.WebhookContext, transfers []models4debtus.TransferEntry) string {
	var s bytes.Buffer
	for _, transfer := range transfers {
		isCreator := whc.AppUserID() == transfer.Data.CreatorUserID
		var counterpartyName string
		if isCreator {
			counterpartyName = transfer.Data.Counterparty().ContactName
		} else {
			counterpartyName = transfer.Data.Creator().ContactName
		}
		amount := fmt.Sprintf(`<a href="%v">%s</a>`,
			common4debtus.GetTransferUrlForUser(
				transfer.ID,
				whc.AppUserID(),
				whc.Locale(),
				common4debtus.NewUtmParams(whc, "history"),
			),
			transfer.Data.GetAmount(),
		)
		if (isCreator && transfer.Data.Direction() == models4debtus.TransferDirectionUser2Counterparty) || (!isCreator && transfer.Data.Direction() == models4debtus.TransferDirectionCounterparty2User) {
			s.WriteString(whc.Translate(trans.MESSAGE_TEXT_HISTORY_ROW_FROM_USER_WITH_NAME, shortDate(transfer.Data.DtCreated, whc), counterpartyName, amount))
		} else {
			s.WriteString(whc.Translate(trans.MESSAGE_TEXT_HISTORY_ROW_TO_USER_WITH_NAME, shortDate(transfer.Data.DtCreated, whc), counterpartyName, amount))
		}

		if transfer.Data.HasInterest() {
			s.WriteString("\n")
			common4debtus.WriteTransferInterest(&s, transfer, whc)
		}
		s.WriteString("\n\n")
	}
	return strings.TrimSpace(s.String())
}

var TransferHistoryCallbackCommand = botsfw.NewCallbackCommand("transfer-history", callbackTransferHistory)

func callbackTransferHistory(whc botsfw.WebhookContext, _ *url.URL) (botsfw.MessageFromBot, error) {
	return whc.NewMessage("TODO: Show more history records"), nil
}

func shortDate(t time.Time, translator i18n.SingleLocaleTranslator) string {
	switch translator.Locale().Code5 {
	case i18n.LocaleCodeEnUS:
		return t.Format("02 Jan 2006")
	case i18n.LocaleCodeFaIR:
		pt := ptime.New(t)
		return pt.Format("dd MMM yyyy")
	default:
		month := t.Format("Jan")
		return fmt.Sprintf("%v %v %v", t.Format("02"), translator.Translate(month), t.Format("2006"))
	}
}
