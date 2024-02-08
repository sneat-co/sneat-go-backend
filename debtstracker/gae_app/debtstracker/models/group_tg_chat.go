package models

//go:generate ffjson $GOFILE

type GroupTgChatJson struct {
	ChatID         int64  `json:"id,omitempty"`
	Title          string `json:"title,omitempty"`
	ChatInviteLink string `json:"link,omitempty"`
}
