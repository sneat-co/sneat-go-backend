package dtb_settings

import (
	"fmt"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/common4debtus"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/debtusbots/profiles/debtusbot/cmd/dtb_transfer"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/debtusbots/profiles/shared_all"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/gae_app/debtstracker/dtdal"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/models4debtus"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/delays4userus"
	"github.com/strongo/logus"
	"regexp"
	"time"
)

/*
Examples:

	receipt-{ContactID}-view_{LANG_CODE5}_[GA_CLIENT_ID]
*/
var reInviteOrReceiptCodeFromStart = regexp.MustCompile(`^(invite|receipt)-(\w+)(-(view|accept|decline))?(_(\w{2}(-\w{2})?))(_(.+))?$`)

func StartInBotAction(whc botsfw.WebhookContext, startParams []string) (m botsfw.MessageFromBot, err error) {
	if len(startParams) == 1 {
		if matched := reInviteOrReceiptCodeFromStart.FindStringSubmatch(startParams[0]); matched != nil {
			return startByLinkCode(whc, matched)
		}
	}
	err = shared_all.ErrUnknownStartParam
	return
}

func startByLinkCode(whc botsfw.WebhookContext, matches []string) (m botsfw.MessageFromBot, err error) {
	c := whc.Context()
	logus.Debugf(c, "startByLinkCode() => matches: %v", matches)
	chatEntity := whc.ChatData()
	entityType := matches[1]
	entityCode := matches[2]
	operation := matches[4]
	localeCode5 := matches[6]
	//gaClientId := matches[8]
	if localeCode5 != "" {
		if len(localeCode5) == 2 {
			localeCode5 = common4debtus.Locale2to5(localeCode5)
		}
		if err = whc.SetLocale(localeCode5); err != nil {
			return
		}
		chatEntity.SetPreferredLanguage(localeCode5)
		if err = delays4userus.DelaySetUserPreferredLocale(c, time.Second, whc.AppUserID(), localeCode5); err != nil {
			return
		}
	}
	switch entityType {
	case "receipt":
		return startReceipt(whc, entityCode, operation, localeCode5)
	case "invite":
		return startInvite(whc, entityCode, operation, localeCode5)
	default:
		err = shared_all.ErrUnknownStartParam
	}
	return
}

func startInvite(whc botsfw.WebhookContext, inviteCode, operation, localeCode5 string) (m botsfw.MessageFromBot, err error) {
	c := whc.Context()
	var invite models4debtus.Invite
	if invite, err = dtdal.Invite.GetInvite(c, nil, inviteCode); err != nil {
		if dal.IsNotFound(err) {
			return whc.NewMessage(fmt.Sprintf("Unknown invite code: %v", inviteCode)), nil
		}
		return
	}
	logus.Debugf(c, "Invite(%v): ClaimedCount=%v, MaxClaimsCount=%v", inviteCode, invite.Data.ClaimedCount, invite.Data.MaxClaimsCount)
	if invite.Data.MaxClaimsCount == 0 || invite.Data.ClaimedCount < invite.Data.MaxClaimsCount {
		return handleInviteOnStart(whc, inviteCode, invite)
	} else {
		m = whc.NewMessage(fmt.Sprintf("Known & already claimed invite code: %v", inviteCode))
	}
	return m, err
}

func startReceipt(whc botsfw.WebhookContext, receiptID, operation, localeCode5 string) (m botsfw.MessageFromBot, err error) {
	c := whc.Context()
	if receiptID == "" {
		return m, fmt.Errorf("receiptID is empty")
	} else if _, err = dtdal.Receipt.GetReceiptByID(c, nil, receiptID); err != nil {
		return
	}
	switch operation {
	case "view":
		if err = whc.SetLocale(localeCode5); err != nil {
			return
		}
		return dtb_transfer.ShowReceipt(whc, receiptID)
	default:
		return dtb_transfer.AcknowledgeReceipt(whc, receiptID, operation)
	}
}
