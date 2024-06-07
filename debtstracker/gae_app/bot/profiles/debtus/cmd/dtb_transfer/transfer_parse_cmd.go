package dtb_transfer

import (
	"bytes"
	"fmt"
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/crediterra/money"
	"github.com/sneat-co/debtstracker-translations/emoji"
	"github.com/sneat-co/debtstracker-translations/trans"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/facade"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/facade/dto"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/models"
	"github.com/strongo/decimal"
	"github.com/strongo/log"
	"strings"
)

var ParseTransferCommand = botsfw.Command{
	Code: "parse-transfer",
	Matcher: func(c botsfw.Command, whc botsfw.WebhookContext) bool {
		input := whc.Input()
		switch input := input.(type) {
		case botsfw.WebhookTextMessage:
			return transferRegex.MatchString(input.Text())
		default:
			return false
		}
	},
	Action: func(whc botsfw.WebhookContext) (m botsfw.MessageFromBot, err error) {
		match := transferRegex.FindStringSubmatch(whc.Input().(botsfw.WebhookTextMessage).Text())
		var verb, valueS, counterpartyName, when string
		var direction models.TransferDirection
		var currency money.CurrencyCode

		for i, name := range transferRegex.SubexpNames() {
			if i != 0 && len(name) > 0 {
				v := strings.TrimSpace(match[i])
				if len(v) > 0 {
					switch name {
					case "verb":
						verb = v
					case "value":
						valueS = v
					case "currency":
						if string(v) == "" {
							currency = money.CurrencyUSD //TODO: Replace with user's default currency
						} else {
							currency = money.CurrencyCode(strings.ToUpper(v))
						}
					case "direction":
						direction = models.TransferDirection(v)
					case "contact":
						counterpartyName = v
					case "when":
						when = v
					}
				}
			}
		}
		if verb == "" {
			switch direction {
			case models.TransferDirectionUser2Counterparty:
				verb = "got"
			case models.TransferDirectionCounterparty2User:
				verb = "gave"
			}
		} else {
			verb = strings.ToLower(verb)
			switch verb {
			case "send":
				verb = "sent"
			case "return":
				verb = "returned"
			}
		}

		m = whc.NewMessage("")

		value, _ := decimal.ParseDecimal64p2(valueS)

		const isReturn = false

		creatorInfo := models.TransferCounterpartyInfo{
			UserID:      whc.AppUserID(),
			ContactName: counterpartyName,
		}
		c := whc.Context()

		from, to := facade.TransferCounterparties(direction, creatorInfo)

		//var botUserEntity botsfwmodels.AppUserData
		//if botUserEntity, err = whc.AppUserData(); err != nil {
		//	return m, err
		//}
		creatorUser := models.NewAppUser(whc.AppUserID(), nil /*botUserEntity.(*models.DebutsAppUserDataOBSOLETE)*/)

		request := dto.CreateTransferRequest{
			IsReturn: isReturn,
			Amount:   money.Amount{Currency: currency, Value: value},
		}
		env := whc.Environment()
		source := GetTransferSource(whc)
		newTransfer := dto.NewTransferInput(
			env,
			source,
			creatorUser,
			request,
			from, to,
		)

		output, err := facade.Transfers.CreateTransfer(c, newTransfer)

		//transferKey, err = nds.Put(ctx, transferKey, transfer)

		if err != nil {
			log.Errorf(c, "Failed to save transfer & counterparty to datastore: %v", err)
			return m, err
		}

		whc.ChatData().SetAwaitingReplyTo(fmt.Sprintf("ask-for-deadline:transferID=%v", output.Transfer.ID))

		m.Keyboard = tgbotapi.NewReplyKeyboardUsingStrings([][]string{
			{whc.Translate(trans.COMMAND_TEXT_YES_IT_HAS_RETURN_DEADLINE) + " " + emoji.ALARM_CLOCK_ICON},
			{whc.Translate(trans.COMMAND_TEXT_NO_IT_CAN_BE_RETURNED_ANYTIME)},
		})

		var buffer bytes.Buffer
		buffer.WriteString(fmt.Sprintf("You've %v %v %v %v %v", verb, valueS, currency, direction, counterpartyName))
		if when != "" {
			//TODO: Convert to time.Time
			buffer.WriteString(" " + when)
		}
		var counterparty models.ContactEntry
		switch direction {
		case models.TransferDirectionUser2Counterparty:
			counterparty = output.To.Contact
		case models.TransferDirectionCounterparty2User:
			counterparty = output.From.Contact
		}
		counterpartyBalance := counterparty.Data.Balance()
		buffer.WriteString(fmt.Sprintf(".\nTotal balance: %v", counterpartyBalance))
		//switch {
		//case counterparty.BalanceJson > 0: buffer.WriteString(fmt.Sprintf(".\nTotal balance: %v ows to you %v %v", contact, counterparty.BalanceJson, currency))
		//case counterparty.BalanceJson < 0: buffer.WriteString(fmt.Sprintf(".\nTotal balance: You owe to %v %v %v", contact, counterparty.BalanceJson, currency))
		//default:
		//}

		switch direction {
		case models.TransferDirectionCounterparty2User:
			buffer.WriteString("\n\nDo you need to return it on a specific date?")
		case models.TransferDirectionUser2Counterparty:
			buffer.WriteString(fmt.Sprintf("\n\nDoes %v have to return it on a specific date?", counterpartyName))
		}
		m.Text = buffer.String()

		return m, nil
	},
}
