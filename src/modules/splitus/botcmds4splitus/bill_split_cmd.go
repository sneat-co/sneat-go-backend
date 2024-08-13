package botcmds4splitus

import (
	"bytes"
	"fmt"
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw-telegram"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/debtstracker-translations/trans"
	"github.com/sneat-co/sneat-go-backend/src/modules/splitus/briefs4splitus"
	"github.com/sneat-co/sneat-go-backend/src/modules/splitus/facade4splitus"
	"github.com/sneat-co/sneat-go-backend/src/modules/splitus/models4splitus"
	"net/url"
)

const billSharesCommandCode = "bill_shares"

var billSharesCommand = billCallbackCommand(billSharesCommandCode,
	func(whc botsfw.WebhookContext, tx dal.ReadwriteTransaction, callbackUrl *url.URL, bill models4splitus.BillEntry) (m botsfw.MessageFromBot, err error) {
		whc.LogRequest()
		ctx := whc.Context()
		members := bill.Data.GetBillMembers()
		if bill.Data.Currency == "" {
			m.BotMessage = telegram.CallbackAnswer(tgbotapi.NewCallback("", whc.Translate(trans.MESSAGE_TEXT_ASK_BILL_CURRENCY)))
			return
		}
		var billID string
		return editSplitCallbackAction(
			whc, callbackUrl,
			billID,
			billCallbackCommandData(billSharesCommandCode, bill.ID),
			billCardCallbackCommandData(bill.ID),
			trans.MESSAGE_TEXT_ASK_HOW_TO_SPLIT_IN_GROP,
			members,
			bill.Data.TotalAmount(),
			func(buffer *bytes.Buffer) error {
				return writeBillCardTitle(ctx, bill, whc.GetBotCode(), buffer, whc)
			},
			func(memberID string, addValue int) (member *briefs4splitus.BillMemberBrief, err error) {
				for i, m := range members {
					if m.ID == memberID {
						m.Shares += addValue
						if m.Shares < 0 {
							m.Shares = 0
						}
						members[i] = m
						bill.Data.SplitMode = models4splitus.SplitModeShare
						member = m
						if err = bill.Data.SetBillMembers(members); err != nil {
							return
						}
						if err = facade4splitus.SaveBill(ctx, tx, bill); err != nil {
							return
						}
						return
					}
				}
				err = fmt.Errorf("member not found by ContactID: %v", member.ID)
				return
			},
		)
	},
)
