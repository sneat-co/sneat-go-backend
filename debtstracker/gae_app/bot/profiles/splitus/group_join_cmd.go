package splitus

import (
	"fmt"
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/debtstracker-translations/trans"
	"net/url"
	"strconv"

	"context"
	"github.com/bots-go-framework/bots-fw-telegram"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/bot/profiles/shared_group"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/dtdal"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/facade"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/models"
	"github.com/strongo/log"
)

const joinGroupCommanCode = "join-group"

var joinGroupCommand = shared_group.GroupCallbackCommand(joinGroupCommanCode,
	func(whc botsfw.WebhookContext, callbackUrl *url.URL, group models.GroupEntry) (m botsfw.MessageFromBot, err error) {
		c := whc.Context()

		userID := whc.AppUserID()
		var appUser models.AppUser
		if group.Data.UserIsMember(userID) {
			if appUser, err = dtdal.User.GetUserByStrID(c, userID); err != nil {
				return
			}
			whc.LogRequest()
			callbackAnswer := tgbotapi.NewCallback("", whc.Translate(trans.ALERT_TEXT_YOU_ARE_ALREADY_MEMBER_OF_THE_GROUP))
			callbackAnswer.ShowAlert = true
			m.BotMessage = telegram.CallbackAnswer(callbackAnswer)
		} else {
			var db dal.DB
			if db, err = facade.GetDatabase(c); err != nil {
				return
			}
			err = db.RunReadwriteTransaction(c, func(c context.Context, tx dal.ReadwriteTransaction) error {
				if appUser, err = dtdal.User.GetUserByStrID(c, userID); err != nil {
					return err
				}
				_, changed, memberIndex, member, members := group.Data.AddOrGetMember(userID, "", appUser.Data.FullName())
				tgUserID := strconv.FormatInt(int64(whc.Input().GetSender().GetID().(int)), 10)
				if member.TgUserID == "" {
					member.TgUserID = tgUserID
					changed = true
				} else {
					if tgUserID != member.TgUserID {
						log.Errorf(c, "tgUserID:%v != member.TgUserID:%v", tgUserID, member.TgUserID)
					}
				}
				switch group.Data.GetSplitMode() {
				case models.SplitModeEqually:
					var shares int
					if group.Data.MembersCount > 0 {
						shares = group.Data.GetGroupMembers()[0].Shares
					} else {
						shares = 1
					}
					if member.Shares != shares {
						member.Shares = shares
						changed = true
					}
				case models.SplitModeShare:
					if member.Shares != 0 {
						member.Shares = 0
						changed = true
					}
				}
				if changed {
					members[memberIndex] = member
					group.Data.SetGroupMembers(members)
					if err = dtdal.Group.SaveGroup(c, tx, group); err != nil {
						return err
					}
				} else {
					log.Debugf(c, "GroupEntry member not changed")
				}
				if userChanged := appUser.Data.AddGroup(group, whc.GetBotCode()); userChanged {
					if err = facade.User.SaveUser(c, tx, appUser); err != nil {
						return err
					}
				}
				if len(members) > 1 {
					groupUsersCount := 0
					for _, m := range members {
						if m.UserID != "" {
							groupUsersCount += 1
						}
					}
					if groupUsersCount > 1 {
						if err = facade.Group.DelayUpdateGroupUsers(c, group.ID); err != nil {
							return err
						}
					}
				}
				return err
			})

			if m, err := showGroupMembers(whc, group, true); err != nil {
				return m, err
			} else if _, err = whc.Responder().SendMessage(c, m, botsfw.BotAPISendMessageOverHTTPS); err != nil {
				return m, err
			}

			m.Text = whc.Translate(trans.MESSAGE_TEXT_USER_JOINED_GROUP, fmt.Sprintf(`<a href="tg://user?id=%v">%v</a>`, whc.MustBotChatID(), appUser.Data.FullName()))
		}

		m.Format = botsfw.MessageFormatHTML
		return
	},
)
