package splitus

import (
	"fmt"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/crediterra/money"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/debtstracker-translations/trans"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/facade"
	"net/url"

	"context"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/bot/profiles/shared_all"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/bot/profiles/shared_group"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/dtdal"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/models"
)

const groupSplitCommandCode = "group-split"

var groupSplitCommand = shared_group.GroupCallbackCommand(groupSplitCommandCode,
	func(whc botsfw.WebhookContext, callbackUrl *url.URL, group models.Group) (m botsfw.MessageFromBot, err error) {
		c := whc.Context()

		members := group.Data.GetMembers()
		billMembers := make([]models.BillMemberJson, len(members))
		for i, m := range members {
			billMembers[i].MemberJson = m
		}
		return editSplitCallbackAction(
			whc, callbackUrl,
			"",
			shared_group.GroupCallbackCommandData(groupSplitCommandCode, group.ID),
			shared_group.GroupCallbackCommandData(shared_all.SettingsCommandCode, group.ID),
			trans.MESSAGE_TEXT_ASK_HOW_TO_SPLIT_IN_GROP,
			billMembers,
			money.Amount{},
			nil,
			func(memberID string, addValue int) (member models.BillMemberJson, err error) {
				var db dal.DB
				if db, err = facade.GetDatabase(c); err != nil {
					return
				}
				err = db.RunReadwriteTransaction(c, func(c context.Context, tx dal.ReadwriteTransaction) (err error) {
					if group, err = dtdal.Group.GetGroupByID(c, tx, group.ID); err != nil {
						return
					}
					members := group.Data.GetGroupMembers()
					for i, m := range members {
						if m.ID == memberID {
							m.Shares += addValue
							if m.Shares < 0 {
								m.Shares = 0
							}
							members[i] = m
							group.Data.SetGroupMembers(members)
							if err = dtdal.Group.SaveGroup(c, tx, group); err != nil {
								return
							}
							member = models.BillMemberJson{MemberJson: m.MemberJson}
							return err
						}
					}
					return fmt.Errorf("member not found by ID: %v", member.ID)
				})
				return
			},
		)
	},
)
