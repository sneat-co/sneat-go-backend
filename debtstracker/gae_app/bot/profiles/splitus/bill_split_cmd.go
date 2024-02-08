package splitus

import (
	"bytes"
	"fmt"
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw-telegram"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/debtstracker-translations/trans"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/dtdal"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/models"
	"net/url"
)

const billSharesCommandCode = "bill_shares"

var billSharesCommand = billCallbackCommand(billSharesCommandCode,
	func(whc botsfw.WebhookContext, tx dal.ReadwriteTransaction, callbackUrl *url.URL, bill models.Bill) (m botsfw.MessageFromBot, err error) {
		whc.LogRequest()
		c := whc.Context()
		members := bill.Data.GetBillMembers()
		if bill.Data.Currency == "" {
			m.BotMessage = telegram.CallbackAnswer(tgbotapi.NewCallback("", whc.Translate(trans.MESSAGE_TEXT_ASK_BILL_CURRENCY)))
			return
		}
		var billID string
		if bill.Data.MembersCount <= 1 {
			billID = bill.ID
		}
		return editSplitCallbackAction(
			whc, callbackUrl,
			billID,
			billCallbackCommandData(billSharesCommandCode, bill.ID),
			billCardCallbackCommandData(bill.ID),
			trans.MESSAGE_TEXT_ASK_HOW_TO_SPLIT_IN_GROP,
			members,
			bill.Data.TotalAmount(),
			func(buffer *bytes.Buffer) error {
				return writeBillCardTitle(c, bill, whc.GetBotCode(), buffer, whc)
			},
			func(memberID string, addValue int) (member models.BillMemberJson, err error) {
				for i, m := range members {
					if m.ID == memberID {
						m.Shares += addValue
						if m.Shares < 0 {
							m.Shares = 0
						}
						members[i] = m
						bill.Data.SplitMode = models.SplitModeShare
						member = m
						if err = bill.Data.SetBillMembers(members); err != nil {
							return
						}
						if err = dtdal.Bill.SaveBill(c, tx, bill); err != nil {
							return
						}
						return
					}
				}
				err = fmt.Errorf("member not found by ID: %v", member.ID)
				return
			},
		)
	},
)
