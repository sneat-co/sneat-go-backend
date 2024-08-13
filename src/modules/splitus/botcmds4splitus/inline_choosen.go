package botcmds4splitus

import (
	"context"
	"errors"
	"fmt"
	"github.com/bots-go-framework/bots-fw-telegram"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/crediterra/money"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/debtstracker-translations/trans"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/const4debtus"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/gae_app/bot/profiles/shared_space"
	"github.com/sneat-co/sneat-go-backend/src/modules/splitus/facade4splitus"
	"github.com/sneat-co/sneat-go-backend/src/modules/splitus/models4splitus"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/dbo4userus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/decimal"
	"github.com/strongo/logus"
	"net/url"
	"regexp"
	"strings"
)

var chosenInlineResultCommand = botsfw.Command{
	Code:       "chosen-inline-result-command",
	InputTypes: []botsfw.WebhookInputType{botsfw.WebhookInputChosenInlineResult},
	Action: func(whc botsfw.WebhookContext) (m botsfw.MessageFromBot, err error) {
		logus.Debugf(whc.Context(), "splitus.chosenInlineResultHandler.Action()")
		chosenResult := whc.Input().(botsfw.WebhookChosenInlineResult)
		resultID := chosenResult.GetResultID()
		if strings.HasPrefix(resultID, "bill?") {
			return createBillFromInlineChosenResult(whc, chosenResult)
		}
		return
	},
}

var reDecimal = regexp.MustCompile(`\d+(\.\d+)?`)

func createBillFromInlineChosenResult(whc botsfw.WebhookContext, chosenResult botsfw.WebhookChosenInlineResult) (m botsfw.MessageFromBot, err error) {
	c := whc.Context()
	logus.Debugf(c, "createBillFromInlineChosenResult()")

	resultID := chosenResult.GetResultID()

	const prefix = "bill?"

	if !strings.HasPrefix(resultID, prefix) {
		err = errors.New("Unexpected resultID: " + resultID)
		return
	}

	switch {
	case true:
		userID := whc.AppUserID()
		var values url.Values
		if values, err = url.ParseQuery(resultID[len(prefix):]); err != nil {
			return
		}
		if lang := values.Get("lang"); lang != "" {
			if err = whc.SetLocale(lang); err != nil {
				return
			}
		}
		var billName string
		if reMatches := reInlineQueryNewBill.FindStringSubmatch(chosenResult.GetQuery()); reMatches != nil {
			billName = strings.TrimSpace(reMatches[3])
		} else {
			billName = whc.Translate(trans.NO_NAME)
		}

		amountStr := values.Get("amount")
		amountIdx := reDecimal.FindStringIndex(amountStr)
		amountNum := amountStr[:amountIdx[1]]
		amountCcy := money.CurrencyCode(amountStr[amountIdx[1]:])

		var amount decimal.Decimal64p2
		if amount, err = decimal.ParseDecimal64p2(amountNum); err != nil {
			return
		}
		bill := models4splitus.BillEntry{
			Data: &models4splitus.BillDbo{
				BillCommon: models4splitus.BillCommon{

					TgInlineMessageIDs: []string{chosenResult.GetInlineMessageID()},
					Name:               billName,
					AmountTotal:        amount,
					Status:             const4debtus.StatusDraft,
					CreatorUserID:      userID,
					UserIDs:            []string{userID},
					SplitMode:          models4splitus.SplitModeEqually,
					Currency:           amountCcy,
				},
			},
		}

		//var (
		//	user          botsfw.BotAppUser
		//	appUserEntity *models.DebutsAppUserDataOBSOLETE
		//)
		//if user, err = whc.GetAppUser(); err != nil {
		//	return
		//}
		//appUserEntity = user.(*models.DebutsAppUserDataOBSOLETE)
		//_, _, _, _, members := bill.AddOrGetMember(userID, 0, appUserEntity.GetFullName())
		//if err = bill.setBillMembers(members); err != nil {
		//	return
		//}
		//billMember.Paid = bill.AmountTotal
		//switch values.Get("i") {
		//case "paid":
		//	billMember.Paid = bill.AmountTotal
		//case "owe":
		//default:
		//	err = fmt.Errorf("unknown value of 'i' parameter: %v", query.Get("i"))
		//	return
		//}

		defer func() {
			if r := recover(); r != nil {
				whc.LogRequest()
				panic(r)
			}
		}()

		user := dbo4userus.NewUserEntry(userID)
		spaceID := user.Data.GetFamilySpaceID()

		err = facade.RunReadwriteTransaction(c, func(c context.Context, tx dal.ReadwriteTransaction) (err error) {
			if bill, err = facade4splitus.CreateBill(c, tx, spaceID, bill.Data); err != nil {
				return
			}
			return
		})
		if err != nil {
			err = fmt.Errorf("failed to call facade4debtus.BillEntry.CreateBill(): %w", err)
			return
		}
		logus.Infof(c, "createBillFromInlineChosenResult() => BillEntry created")

		botCode := whc.GetBotCode()

		logus.Infof(c, "createBillFromInlineChosenResult() => suxx 0")

		footer := strings.Repeat("―", 15) + "\n" + whc.Translate(trans.MESSAGE_TEXT_ASK_BILL_PAYER)

		if m.Text, err = getBillCardMessageText(c, botCode, whc, bill, false, footer); err != nil {
			logus.Errorf(c, "Failed to create bill card")
			return
		} else if strings.TrimSpace(m.Text) == "" {
			err = errors.New("getBillCardMessageText() returned empty string")
			logus.Errorf(c, err.Error())
			return
		}

		logus.Infof(c, "createBillFromInlineChosenResult() => suxx 1")

		if m, err = whc.NewEditMessage(m.Text, botsfw.MessageFormatHTML); err != nil { // TODO: Unnecessary hack?
			logus.Infof(c, "createBillFromInlineChosenResult() => suxx 1.2")
			logus.Errorf(c, err.Error())
			return
		}

		logus.Infof(c, "createBillFromInlineChosenResult() => suxx 2")

		m.Keyboard = getWhoPaidInlineKeyboard(whc, bill.ID)

		var response botsfw.OnMessageSentResponse
		logus.Debugf(c, "createBillFromInlineChosenResult() => Sending bill card: %v", m)

		if response, err = whc.Responder().SendMessage(c, m, botsfw.BotAPISendMessageOverHTTPS); err != nil {
			logus.Errorf(c, "createBillFromInlineChosenResult() => %v", err)
			return
		}

		logus.Debugf(c, "response: %v", response)
		m.Text = botsfw.NoMessageToSend
	}

	return
}

var reBillUrl = regexp.MustCompile(`\?start=bill-(\d+)$`)

func getBillIDFromUrlInEditedMessage(whc botsfw.WebhookContext) (billID string) {
	tgInput, ok := whc.Input().(telegram.TgWebhookInput)
	if !ok {
		return
	}
	tgUpdate := tgInput.TgUpdate()
	if tgUpdate.EditedMessage == nil {
		return
	}
	if tgUpdate.EditedMessage.Entities == nil {
		return
	}
	for _, entity := range *tgUpdate.EditedMessage.Entities {
		if entity.Type == "text_link" {
			if s := reBillUrl.FindStringSubmatch(entity.URL); len(s) != 0 {
				billID = s[1]
				if billID == "" {
					logus.Errorf(whc.Context(), "Missing bill ContactID")
				}
				return
			}
		}
	}
	return
}

var EditedBillCardHookCommand = botsfw.Command{ // TODO: seems to be not used anywhere
	Code: "edited-bill-card",
	Action: func(whc botsfw.WebhookContext) (m botsfw.MessageFromBot, err error) {
		whc.LogRequest()
		c := whc.Context()
		billID := getBillIDFromUrlInEditedMessage(whc)
		logus.Debugf(c, "editedBillCardHookCommand.Action() => billID: %s", billID)
		if billID == "" {
			panic("billID is empty string")
		}

		m.Text = botsfw.NoMessageToSend

		var groupID string
		if groupID, err = shared_space.GetUserGroupID(whc); err != nil {
			return
		} else if groupID == "" {
			logus.Warningf(c, "group.ContactID is empty string")
			return
		}

		changed := false

		err = facade.RunReadwriteTransaction(c, func(tc context.Context, tx dal.ReadwriteTransaction) (err error) {
			var bill models4splitus.BillEntry
			if bill, err = facade4splitus.GetBillByID(c, tx, billID); err != nil {
				return err
			}

			if groupID != "" && bill.Data.GetUserGroupID() != groupID { // TODO: Should we check for empty bill.GetUserGroupID() or better fail?
				if bill, _, err = facade4splitus.AssignBillToGroup(c, tx, bill, groupID, whc.AppUserID()); err != nil {
					return err
				}
				changed = true
			}

			if changed {
				return facade4splitus.SaveBill(c, tx, bill)
			}

			return err
		})
		if err != nil {
			return
		}
		if changed {
			logus.Debugf(c, "BillEntry updated with group ContactID")
		}
		return
	},
	Matcher: func(command botsfw.Command, whc botsfw.WebhookContext) (result bool) {
		result = whc.IsInGroup() && getBillIDFromUrlInEditedMessage(whc) != ""
		logus.Debugf(whc.Context(), "editedBillCardHookCommand.Matcher(): %v", result)
		return
	},
}
