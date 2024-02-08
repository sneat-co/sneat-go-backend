package dtb_transfer

import (
	"bytes"
	"fmt"
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/debtstracker-translations/trans"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/facade"
	"github.com/strongo/i18n"
	"html"
	"net/url"
	"strings"
	"time"

	"github.com/sneat-co/debtstracker-translations/emoji"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/dtdal"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/models"
	"github.com/strongo/log"
)

const DUE_RETURNS_COMMAND = "due-returns"

var DueReturnsCallbackCommand = botsfw.NewCallbackCommand(DUE_RETURNS_COMMAND, dueReturnsCallbackAction)

func dueReturnsCallbackAction(whc botsfw.WebhookContext, _ *url.URL) (m botsfw.MessageFromBot, err error) {
	c := whc.Context()

	userID := whc.AppUserID()
	var (
		overdueTransfers, dueTransfers []models.Transfer
	)

	er := make(chan error, 2)
	go func(er chan<- error) {
		var db dal.DB
		if db, err = facade.GetDatabase(c); err != nil {
			er <- err
			return
		}
		if overdueTransfers, err = dtdal.Transfer.LoadOverdueTransfers(c, db, userID, 5); err != nil {
			er <- fmt.Errorf("failed to get overdue transfers: %w", err)
		} else {
			log.Debugf(c, "Loaded %v overdue transfer", len(overdueTransfers))
			er <- nil
		}
	}(er)
	go func(er chan<- error) {
		var db dal.DB
		if db, err = facade.GetDatabase(c); err != nil {
			er <- err
			return
		}
		if dueTransfers, err = dtdal.Transfer.LoadDueTransfers(c, db, userID, 5); err != nil {
			er <- fmt.Errorf("failed to get due transfers: %w", err)
		} else {
			log.Debugf(c, "Loaded %v due transfer", len(dueTransfers))
			er <- nil
		}
	}(er)

	for i := 0; i < 2; i++ {
		if err = <-er; err != nil {
			return
		}
	}

	if len(overdueTransfers) == 0 || len(dueTransfers) == 0 {
		if m, err = whc.NewEditMessage(whc.Translate(trans.MESSAGE_TEXT_DUE_RETURNS_EMPTY), botsfw.MessageFormatHTML); err != nil {
			return
		}
	} else {
		var buffer bytes.Buffer

		now := time.Now()
		listTransfers := func(header string, transfers []models.Transfer) {
			if len(transfers) == 0 {
				return
			}
			buffer.WriteString(whc.Translate(header))
			buffer.WriteString("\n\n")
			for i, transfer := range transfers {
				switch transfer.Data.Direction() {
				case models.TransferDirectionCounterparty2User:
					buffer.WriteString(whc.Translate(trans.MESSAGE_TEXT_DUE_RETURNS_ROW_BY_USER, html.EscapeString(transfer.Data.Counterparty().ContactName), transfer.Data.GetAmount(), DurationToString(transfer.Data.DtDueOn.Sub(now), whc)))
				case models.TransferDirectionUser2Counterparty:
					buffer.WriteString(whc.Translate(trans.MESSAGE_TEXT_DUE_RETURNS_ROW_BY_COUNTERPARTY, html.EscapeString(transfer.Data.Counterparty().ContactName), transfer.Data.GetAmount(), DurationToString(transfer.Data.DtDueOn.Sub(now), whc)))
				default:
					panic(fmt.Sprintf("Unknown direction for transfer id=%v: %v", transfers[i].ID, transfer))
				}
				buffer.WriteString("\n")
			}
			buffer.WriteString("\n")
		}
		listTransfers(trans.MESSAGE_TEXT_OVERDUE_RETURNS_HEADER, overdueTransfers)
		listTransfers(trans.MESSAGE_TEXT_DUE_RETURNS_HEADER, dueTransfers)
		if m, err = whc.NewEditMessage(strings.TrimSuffix(buffer.String(), "\n"), botsfw.MessageFormatHTML); err != nil {
			return
		}
	}
	m.Keyboard = tgbotapi.NewInlineKeyboardMarkup(
		[]tgbotapi.InlineKeyboardButton{
			{
				Text:         whc.CommandText(trans.COMMAND_TEXT_BALANCE, emoji.BALANCE_ICON),
				CallbackData: BALANCE_COMMAND,
			},
		},
	)

	return m, err
}

func DurationToString(d time.Duration, translator i18n.SingleLocaleTranslator) string {
	hours := d.Hours()
	switch hours {
	case 0:
		switch d.Minutes() {
		case 0:
			return translator.Translate(trans.DUE_IN_NOW)
		case 1:
			return translator.Translate(trans.DUE_IN_A_MINUTE)
		default:
			return fmt.Sprintf(translator.Translate(trans.DUE_IN_X_MINUTES), d.Minutes())
		}
	case 1:
		return translator.Translate(trans.DUE_IN_AN_HOUR)
	default:
		if hours < 24 {
			return fmt.Sprintf(translator.Translate(trans.DUE_IN_X_HOURS), int(hours))
		}
		return fmt.Sprintf(translator.Translate(trans.DUE_IN_X_DAYS), int(hours/24))
	}
}
