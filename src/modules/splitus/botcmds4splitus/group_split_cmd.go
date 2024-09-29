package botcmds4splitus

import (
	"context"
	"fmt"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/crediterra/money"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/debtstracker-translations/trans"
	"github.com/sneat-co/sneat-go-backend/src/coremodules/anybot/cmds4anybot"
	"github.com/sneat-co/sneat-go-backend/src/coremodules/spaceus/dbo4spaceus"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/debtusbots/profiles/shared_space"
	"github.com/sneat-co/sneat-go-backend/src/modules/splitus/briefs4splitus"
	"github.com/sneat-co/sneat-go-core/facade"
	"net/url"
)

const spaceSplitCommandCode = "space-split"

var spaceSplitCommand = shared_space.SpaceCallbackCommand(spaceSplitCommandCode,
	func(whc botsfw.WebhookContext, callbackUrl *url.URL, space dbo4spaceus.SpaceEntry) (m botsfw.MessageFromBot, err error) {
		ctx := whc.Context()

		//members := space.Data.GetMembers()
		billMembers := make([]*briefs4splitus.BillMemberBrief, 0 /*len(members)*/)
		//for i, m := range members {
		//	billMembers[i].MemberBrief = m
		//}
		return editSplitCallbackAction(
			whc, callbackUrl,
			"",
			shared_space.SpaceCallbackCommandData(spaceSplitCommandCode, space.ID),
			shared_space.SpaceCallbackCommandData(cmds4anybot.SettingsCommandCode, space.ID),
			trans.MESSAGE_TEXT_ASK_HOW_TO_SPLIT_IN_GROP,
			billMembers,
			money.Amount{},
			nil,
			func(memberID string, addValue int) (member *briefs4splitus.BillMemberBrief, err error) {
				err = facade.RunReadwriteTransaction(ctx, func(ctx context.Context, tx dal.ReadwriteTransaction) (err error) {
					//if space, err = dtdal.Group.GetGroupByID(ctx, tx, space.ContactID); err != nil {
					//	return
					//}
					//members := space.Data.GetGroupMembers()
					//for i, m := range members {
					//	if m.ContactID == memberID {
					//		m.Shares += addValue
					//		if m.Shares < 0 {
					//			m.Shares = 0
					//		}
					//		members[i] = m
					//		space.Data.SetGroupMembers(members)
					//		if err = dtdal.Group.SaveGroup(ctx, tx, space); err != nil {
					//			return
					//		}
					//		member = briefs4splitus.BillMemberBrief{MemberBrief: m.MemberBrief}
					//		return err
					//	}
					//}
					return fmt.Errorf("not implemented yet: member not found by ContactID: %v", member.ID)
				})
				return
			},
		)
	},
)
