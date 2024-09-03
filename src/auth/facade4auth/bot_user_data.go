package facade4auth

import "github.com/bots-go-framework/bots-fw/botsfwconst"

type BotUserData struct {
	PlatformID   botsfwconst.Platform
	BotID        string
	BotUserID    string
	FirstName    string
	LastName     string
	Username     string
	PhotoURL     string
	LanguageCode string
}
