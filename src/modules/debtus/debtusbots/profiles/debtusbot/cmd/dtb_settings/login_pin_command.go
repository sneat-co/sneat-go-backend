package dtb_settings

import (
	"errors"
	"fmt"
	"github.com/bots-go-framework/bots-fw/botinput"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/common4debtus"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/facade4debtus"
	"strconv"
	"strings"
)

var LoginPinCommand = botsfw.Command{
	Code: "LoginPin",
	Matcher: func(cmd botsfw.Command, whc botsfw.WebhookContext) bool {
		return false
		//if whc.BotPlatform().ContactID() == viber.PlatformID && whc.InputType() == botsfw.WebhookInputText {
		//	context := whc.Input().(viber.WebhookInputConversationStarted).GetContext()
		//	return strings.HasPrefix(context, "login-")
		//} else {
		//	return false
		//}
	},
	Action: func(whc botsfw.WebhookContext) (m botsfw.MessageFromBot, err error) {
		mt := whc.Input().(botinput.WebhookTextMessage).Text()
		context := strings.Split(mt, " ")[0]
		contextParams := strings.Split(context, "_")
		var (
			loginID int
			//gacID string
			lang string
		)
		if len(contextParams) < 2 || len(contextParams) > 3 {
			return m, fmt.Errorf("len(contextParams): %v", len(contextParams))
		}
		for _, p := range contextParams {
			switch {
			case strings.HasPrefix(p, "login-"):
				if loginID, err = strconv.Atoi(p[len("login-"):]); err != nil {
					err = errors.New(whc.Translate("Parameter 'login_id'  should be an integer."))
					return m, err
				}
			case strings.HasPrefix(p, "lang-"):
				lang = common4debtus.Locale2to5(p[len("lang-"):])
				if err = whc.SetLocale(lang); err != nil {
					return m, err
				}
				whc.ChatData().SetPreferredLanguage(lang)
				//case strings.HasPrefix(p,"gac-"):
				//	gacID = p[len("gac-"):]
			}
		}
		ctx := whc.Context()
		//goland:noinspection GoDeprecation
		if pinCode, err := facade4debtus.AuthFacade.AssignPinCode(ctx, loginID, whc.AppUserID()); err != nil {
			return m, err
		} else {
			return whc.NewMessage(fmt.Sprintf("Login PIN code: %v", pinCode)), nil
		}
	},
}
