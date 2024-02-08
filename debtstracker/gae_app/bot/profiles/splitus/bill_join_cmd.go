package splitus

import (
	"context"
	"errors"
	"fmt"
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw-store/botsfwmodels"
	"github.com/bots-go-framework/bots-fw-telegram"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/crediterra/money"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/debtstracker-translations/trans"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/bot/profiles/shared_all"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/bot/profiles/shared_group"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/dtdal"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/facade"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/models"
	"github.com/strongo/decimal"
	"github.com/strongo/i18n"
	"github.com/strongo/log"
	"net/url"
	"strings"
	"time"
)

const joinBillCommandCode = "join_bill"
const leaveBillCommandCode = "leave_bill"

var joinBillCommand = botsfw.Command{
	Code: joinBillCommandCode,
	Action: func(whc botsfw.WebhookContext) (m botsfw.MessageFromBot, err error) {
		text := whc.Input().(botsfw.WebhookTextMessage).Text()
		var bill models.Bill
		if bill.ID = strings.Replace(text, "/start join_bill-", "", 1); bill.ID == "" {
			err = errors.New("Missing bill ID")
			return
		}
		var db dal.DB
		if db, err = facade.GetDatabase(whc.Context()); err != nil {
			return
		}
		if err = db.RunReadwriteTransaction(whc.Context(), func(c context.Context, tx dal.ReadwriteTransaction) (err error) {
			if bill, err = facade.GetBillByID(whc.Context(), tx, bill.ID); err != nil {
				return
			}
			m, err = joinBillAction(whc, tx, bill, "", false)
			return
		}, dal.TxWithCrossGroup()); err != nil {
			return
		}
		return
	},
	CallbackAction: func(whc botsfw.WebhookContext, callbackUrl *url.URL) (m botsfw.MessageFromBot, err error) {
		_ = whc.AppUserID() // Make sure we have user before transaction starts, TODO: it smells, should be refactored?
		//
		return shared_all.TransactionalCallbackAction(billCallbackAction(func(whc botsfw.WebhookContext, tx dal.ReadwriteTransaction, callbackUrl *url.URL, bill models.Bill) (m botsfw.MessageFromBot, err error) {
			c := whc.Context()
			log.Debugf(c, "joinBillCommand.CallbackAction()")
			memberStatus := callbackUrl.Query().Get("i")
			m, err = joinBillAction(whc, tx, bill, memberStatus, true)
			return
		}))(whc, callbackUrl)
	},
}

func joinBillAction(whc botsfw.WebhookContext, tx dal.ReadwriteTransaction, bill models.Bill, memberStatus string, isEditMessage bool) (m botsfw.MessageFromBot, err error) {

	if bill.ID == "" {
		panic("bill.ID is empty string")
	}
	c := whc.Context()
	log.Debugf(c, "joinBillAction(bill.ID=%v)", bill.ID)

	userID := whc.AppUserID()
	var appUserData botsfwmodels.AppUserData
	if appUserData, err = whc.AppUserData(); err != nil {
		return
	}

	type User interface {
		GetPrimaryCurrency() string
		GetLastCurrencies() []string
		FullName() string
	}

	isAlreadyMember := func(members []models.BillMemberJson) (member models.BillMemberJson, isMember bool) {
		for _, member = range bill.Data.GetBillMembers() {
			if isMember = member.UserID == userID; isMember {
				return
			}
		}
		return
	}

	_, isMember := isAlreadyMember(bill.Data.GetBillMembers())

	user, isUser := appUserData.(User)
	if !isUser {
		err = errors.New("failed to cast appUserData to User")
		return
	}

	userName := user.FullName()

	if userName == "" {
		err = errors.New("userName is empty string")
		return
	}

	if memberStatus == "" && isMember {
		log.Infof(c, "User is already member of the bill before transaction, memberStatus: "+memberStatus)
		callbackAnswer := tgbotapi.NewCallback("", whc.Translate(trans.MESSAGE_TEXT_ALREADY_BILL_MEMBER, userName))
		callbackAnswer.ShowAlert = true
		m.BotMessage = telegram.CallbackAnswer(callbackAnswer)
		whc.LogRequest()
		if update := whc.Input().(telegram.TgWebhookInput).TgUpdate(); update.CallbackQuery.Message != nil {
			if m2, err := ShowBillCard(whc, true, bill, ""); err != nil {
				return m2, err
			} else if m2.Text != update.CallbackQuery.Message.Text {
				log.Debugf(c, "Need to update bill card")
				if _, err = whc.Responder().SendMessage(c, m2, botsfw.BotAPISendMessageOverHTTPS); err != nil {
					return m2, err
				}
			} else {
				log.Debugf(c, "m.Text: %v", m2.Text)
			}
		}
		return
	}

	//if err = dtdal.DB.RunInTransaction(c, func(c context.Context) (err error) {
	//if bill, err = facade.GetBillByID(c, bill.ID); err != nil {
	//	return
	//}

	billChanged := false
	if bill.Data.Currency == "" {
		guessCurrency := func() money.CurrencyCode {
			switch whc.Locale().Code5 {
			case i18n.LocalCodeRuRu:
				return money.CurrencyRUB
			case i18n.LocaleCodeFrFR, i18n.LocaleCodeDeDE, i18n.LocaleCodeItIT, i18n.LocaleCodePtPT:
				return money.CurrencyEUR
			case i18n.LocaleCodeEnUK:
				return money.CurrencyGBP
			default:
				return money.CurrencyUSD
			}
		}

		if whc.IsInGroup() {
			var group models.Group
			if group, err = shared_group.GetGroup(whc, nil); err != nil {
				return
			}
			if group.Data != nil {
				if group.Data.DefaultCurrency != "" {
					bill.Data.Currency = group.Data.DefaultCurrency
				} else {
					bill.Data.Currency = guessCurrency()
				}
			}
		} else if primaryCurrency := user.GetPrimaryCurrency(); primaryCurrency != "" {
			bill.Data.Currency = money.CurrencyCode(primaryCurrency)
		} else if lastCurrencies := user.GetLastCurrencies(); len(lastCurrencies) > 0 {
			bill.Data.Currency = money.CurrencyCode(lastCurrencies[0])
		}
		if bill.Data.Currency == "" {
			bill.Data.Currency = guessCurrency()
		}
		billChanged = true
	}

	var isJoined bool

	var paid decimal.Decimal64p2
	switch memberStatus {
	case "paid":
		paid = bill.Data.AmountTotal
	case "owe":
	default:
	}

	billChanged2 := false
	if bill, _, billChanged2, isJoined, err = facade.Bill.AddBillMember(c, tx, userID, bill, "", userID, userName, paid); err != nil {
		return
	}
	if billChanged = billChanged2 || billChanged; billChanged {
		if err = dtdal.Bill.SaveBill(c, tx, bill); err != nil {
			return
		}
		if isJoined {
			if err = delayUpdateBillCardOnUserJoin(c, bill.ID, whc.Translate(fmt.Sprintf("%v: ", time.Now())+trans.MESSAGE_TEXT_USER_JOINED_BILL, userName)); err != nil {
				log.Errorf(c, "failed to daley update bill card on user join: %v", err)
			}
		}
	}
	//return
	//}

	return ShowBillCard(whc, isEditMessage, bill, "")
}
