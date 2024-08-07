package dal4listus

import "github.com/bots-go-framework/bots-fw-store/botsfwmodels"

type WithBotIDs struct {
	BotIDs []string `json:"botIDs" firestore:"botIDs"`
}

type WithModuleIDs struct {
	ModuleIDs []string `json:"moduleIDs" firestore:"moduleIDs"`
}

type ListusChatData struct {
	botsfwmodels.ChatBaseData
	SpaceID string `json:"spaceID" firestore:"spaceID"`
	ListID  string `json:"listID" firestore:"listID"`
}
