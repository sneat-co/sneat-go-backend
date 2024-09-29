package botcmds4splitus

import (
	"errors"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/sneat-co/sneat-core-modules/spaceus/dbo4spaceus"
	"github.com/sneat-co/sneat-mod-debtus-go/debtus/debtusbots/profiles/shared_space"
	"net/url"
)

const joinSpaceCommandCode = "join-space"

var joinSpaceCommand = shared_space.SpaceCallbackCommand(joinSpaceCommandCode,
	func(whc botsfw.WebhookContext, callbackUrl *url.URL, space dbo4spaceus.SpaceEntry) (m botsfw.MessageFromBot, err error) {
		err = errors.New("joinSpaceCommand is not implemented")
		return
		//ctx := whc.Context()
		//
		//userID := whc.AppUserID()
		//
		//m.Format = botsfw.MessageFormatHTML
		//
		//if space.Data.HasUserID(userID) {
		//	user := dbo4userus.NewUserEntry(userID)
		//	if err = facade4userus.GetUserByIdOBSOLETE(ctx, nil, user.Record); err != nil {
		//		return
		//	}
		//	whc.LogRequest()
		//	callbackAnswer := tgbotapi.NewCallback("", whc.Translate(trans.ALERT_TEXT_YOU_ARE_ALREADY_MEMBER_OF_THE_GROUP))
		//	callbackAnswer.ShowAlert = true
		//	m.BotMessage = telegram.CallbackAnswer(callbackAnswer)
		//	return
		//}
		//
		//userContext := facade.NewUserContext(userID)
		//err = facade4userus.RunUserWorker(c, userContext, func(ctx context.Context, tx dal.ReadwriteTransaction, userWorkerParams *facade4userus.UserWorkerParams) (err error) {
		//	appUser := userWorkerParams.User
		//	_, changed, memberIndex, member, members := space.Data.AddOrGetMember(userID, "", appUser.Data.GetFullName())
		//	tgUserID := strconv.FormatInt(int64(whc.Input().GetSender().GetID().(int)), 10)
		//	if member.TgUserID == "" {
		//		member.TgUserID = tgUserID
		//		changed = true
		//	} else {
		//		if tgUserID != member.TgUserID {
		//			logus.Errorf(c, "tgUserID:%v != member.TgUserID:%v", tgUserID, member.TgUserID)
		//		}
		//	}
		//	switch space.Data.GetSplitMode() {
		//	case models4splitus.SplitModeEqually:
		//		var shares int
		//		if space.Data.MembersCount > 0 {
		//			shares = space.Data.GetGroupMembers()[0].Shares
		//		} else {
		//			shares = 1
		//		}
		//		if member.Shares != shares {
		//			member.Shares = shares
		//			changed = true
		//		}
		//	case models4splitus.SplitModeShare:
		//		if member.Shares != 0 {
		//			member.Shares = 0
		//			changed = true
		//		}
		//	}
		//	if changed {
		//		members[memberIndex] = member
		//		space.Data.SetGroupMembers(members)
		//		if err = dtdal.Group.SaveGroup(c, tx, space); err != nil {
		//			return err
		//		}
		//	} else {
		//		logus.Debugf(c, "GroupEntry member not changed")
		//	}
		//	if userChanged := appUser.Data.AddGroup(space, whc.GetBotCode()); userChanged {
		//		userWorkerParams.User.Record.MarkAsChanged()
		//	}
		//	if len(members) > 1 {
		//		groupUsersCount := 0
		//		for _, m := range members {
		//			if m.UserID != "" {
		//				groupUsersCount += 1
		//			}
		//		}
		//		if groupUsersCount > 1 {
		//			if err = facade4splitus.Group.DelayUpdateGroupUsers(c, space.ContactID); err != nil {
		//				return err
		//			}
		//		}
		//	}
		//	return err
		//})
		//
		//if m, err := showGroupMembers(whc, space, true); err != nil {
		//	return m, err
		//} else if _, err = whc.Responder().SendMessage(c, m, botsfw.BotAPISendMessageOverHTTPS); err != nil {
		//	return m, err
		//}
		//
		//m.Text = whc.Translate(trans.MESSAGE_TEXT_USER_JOINED_GROUP, fmt.Sprintf(`<a href="tg://user?id=%v">%v</a>`, whc.MustBotChatID(), appUser.Data.GetFullName()))
		//
		//return
	},
)
