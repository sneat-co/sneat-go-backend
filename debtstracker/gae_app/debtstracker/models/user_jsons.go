package models

//go:generate ffjson $GOFILE

// This class is used by User entity to store active groups information in a JSON property user.GroupsJsonActive
type UserGroupJson struct {
	ID           string
	Name         string
	Note         string   `json:",omitempty"`
	Status       string   `json:",omitempty"`
	MembersCount int      `json:",omitempty"`
	TgBots       []string `json:",omitempty"`
}
