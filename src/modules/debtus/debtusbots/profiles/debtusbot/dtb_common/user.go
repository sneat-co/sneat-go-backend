package dtb_common

import (
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/sneat-co/sneat-core-modules/userus/dbo4userus"
)

func GetUser(whc botsfw.WebhookContext) (user dbo4userus.UserEntry, err error) {
	panic("not implemented")
	//var appUser botsfwmodels.BotAppUser
	//if appUser, err = whc.GetAppUser(); err != nil {
	//	return
	//}
	//user.Data = appUser.(*models.DebutsAppUserDataOBSOLETE)
	//user.ContactID, err = strconv.ParseInt(whc.AppUserID(), 10, 64)
	//return
}
