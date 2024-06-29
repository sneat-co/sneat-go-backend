package splitus

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/crediterra/money"
	"github.com/sneat-co/debtstracker-translations/trans"
	"github.com/strongo/i18n"
	"github.com/strongo/logus"
	"html"
	"net/url"
	"strconv"

	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/models"
	"github.com/strongo/decimal"
)

func editSplitCallbackAction(
	whc botsfw.WebhookContext,
	callbackUrl *url.URL,
	billID string,
	editCommandPrefix, backCommandPrefix string,
	msgTextAskToSplit string,
	members []models.BillMemberJson,
	totalAmount money.Amount,
	writeTitle func(buffer *bytes.Buffer) error,
	addShares func(memberID string, addValue int) (member models.BillMemberJson, err error),
) (m botsfw.MessageFromBot, err error) {
	c := whc.Context()

	q := callbackUrl.Query()

	var (
		addValue int
		member   models.BillMemberJson
	)

	if member, addValue, err = getSplitParamsAndCurrentMember(q, members); err != nil {
		return
	}

	logus.Debugf(c, "current member: %v", member)

	if addValue != 0 {
		logus.Debugf(c, "add=%d", addValue)

		if member, err = addShares(member.ID, addValue); err != nil {
			return
		}
	}

	buffer := new(bytes.Buffer)

	if writeTitle != nil {
		if err = writeTitle(buffer); err != nil {
			return
		}
		buffer.WriteString("\n\n")
	}

	fmt.Fprintf(buffer, "<b>%v</b>\n\n", whc.Translate(msgTextAskToSplit))

	writeSplitMembers(buffer, members, member.ID, totalAmount.Currency)

	if len(members) > 1 {
		writeSplitInstructions(buffer, member.TgUserID, member.Name)
	}

	m.Text = buffer.String()
	m.Format = botsfw.MessageFormatHTML

	tgKeyboard := &tgbotapi.InlineKeyboardMarkup{}
	tgKeyboard.InlineKeyboard = addEditSplitInlineKeyboardButtons(tgKeyboard.InlineKeyboard, whc,
		len(members),
		billID,
		editCommandPrefix+"&m="+member.ID+"&",
		backCommandPrefix,
	)
	m.Keyboard = tgKeyboard

	m.IsEdit = true
	return
}

func getSplitParamsAndCurrentMember(q url.Values, members []models.BillMemberJson) (member models.BillMemberJson, add int, err error) {
	if len(members) == 0 {
		err = errors.New("len(members) == 0")
		return
	}

	if memberID := q.Get("m"); memberID == "" {
		member = members[0]
	} else if memberID == "0" {
		err = errors.New("parameter 'm' is 0")
		return
	} else {
		member.ID = q.Get("m")
		var (
			i    int
			m    models.BillMemberJson
			move string
		)
		for i, m = range members {
			if m.ID == member.ID {
				break
			}
		}

		if move = q.Get("move"); move != "" {
			switch move {
			case "up":
				if i -= 1; i < 0 {
					if i = len(members) - 1; i < 0 {
						i = 0
					}
				}
			case "down":
				if i += 1; i >= len(members) {
					i = 0
				}
			default:
				err = fmt.Errorf("unknown move: %v", q.Get("move"))
				return
			}
			member = members[i]
		} else {
			if addStr := q.Get("add"); addStr != "" {
				if add, err = strconv.Atoi(addStr); err != nil {
					return
				}
			}
		}
	}

	return
}

func writeSplitInstructions(buffer *bytes.Buffer, tgUserID, memberName string) {
	buffer.WriteString("Use ⬆ & ⬇ to choose a member.")
	buffer.WriteString("\n\n")
	if tgUserID == "" {
		buffer.WriteString(fmt.Sprintf("<b>Selected:</b> %v", memberName))
	} else {
		buffer.WriteString(fmt.Sprintf(`<b>Selected:</b> <a href="tg://user?id=%v">%v</a>`, tgUserID, memberName))
	}
}

func writeSplitMembers(buffer *bytes.Buffer, members []models.BillMemberJson, currentMemberID string, currency money.CurrencyCode) {
	var totalShares int
	for _, m := range members {
		totalShares += m.Shares
	}
	if totalShares == 0 {
		totalShares = 1
	}
	for i, m := range members {
		if m.ID == currentMemberID {
			buffer.WriteString(fmt.Sprintf("  <b>%d. %v</b>\n", i+1, html.EscapeString(m.Name)))
		} else {
			buffer.WriteString(fmt.Sprintf("  %d. %v\n", i+1, html.EscapeString(m.Name)))
		}
		buffer.WriteString(fmt.Sprintf("     <i>Shares: %d</i> — <code>%v%%</code>", m.Shares, decimal.Decimal64p2(m.Shares*100*100/totalShares)))
		if m.Owes != 0 {
			buffer.WriteString(" = " + money.Amount{Currency: currency, Value: m.Owes}.String())
		}
		buffer.WriteString("\n\n")
	}
}

func addEditSplitInlineKeyboardButtons(kb [][]tgbotapi.InlineKeyboardButton, translator i18n.SingleLocaleTranslator, membersCount int, billID, callbackDataPrefix, backCallbackData string) [][]tgbotapi.InlineKeyboardButton {
	var lastRow []tgbotapi.InlineKeyboardButton
	if membersCount > 1 {
		kb = append(kb, // TODO: Move to Telegram specific package
			[]tgbotapi.InlineKeyboardButton{
				{
					Text:         "-10",
					CallbackData: callbackDataPrefix + "add=-10",
				},
				{
					Text:         "-1",
					CallbackData: callbackDataPrefix + "add=-1",
				},
				{
					Text:         "==",
					CallbackData: callbackDataPrefix + "set=50x50",
				},
				{
					Text:         "+1",
					CallbackData: callbackDataPrefix + "add=1",
				},
				{
					Text:         "+10",
					CallbackData: callbackDataPrefix + "add=10",
				},
			},
		)
		lastRow = append(lastRow,
			tgbotapi.InlineKeyboardButton{
				Text:         "⬆️",
				CallbackData: callbackDataPrefix + "move=up",
			},
			tgbotapi.InlineKeyboardButton{
				Text:         "⬇️",
				CallbackData: callbackDataPrefix + "move=down",
			},
		)
	} else {
		lastRow = append(lastRow,
			tgbotapi.InlineKeyboardButton{
				Text:         translator.Translate(trans.BUTTON_TEXT_JOIN),
				CallbackData: billCallbackCommandData(joinBillCommandCode, billID),
			},
		)
	}
	lastRow = append(lastRow, tgbotapi.InlineKeyboardButton{
		Text:         translator.Translate("✅ Done"),
		CallbackData: backCallbackData,
	})

	return append(kb, lastRow)
}
