package models4debtus

import (
	"fmt"
	"github.com/pquerna/ffjson/ffjson"
	"github.com/sneat-co/sneat-go-backend/src/modules/splitus/models4splitus"
)

// UserGroupJson is used by SpaceDbo entity to store active groups information in a JSON property user.GroupsJsonActive
type UserGroupJson struct {
	ID           string
	Name         string
	Note         string   `json:",omitempty"`
	Status       string   `json:",omitempty"`
	MembersCount int      `json:",omitempty"`
	TgBots       []string `json:",omitempty"`
}

// Deprecated: use SplitusSpaceDbo
type WithGroups struct {
	GroupsCountActive   int `firestore:",omitempty"`
	GroupsCountArchived int `firestore:",omitempty"`

	GroupsJsonActive   string `firestore:",omitempty"`
	GroupsJsonArchived string `firestore:",omitempty"`
}

func (entity *WithGroups) ActiveGroups() (groups []UserGroupJson) {
	if entity.GroupsJsonActive != "" {
		if err := ffjson.Unmarshal([]byte(entity.GroupsJsonActive), &groups); err != nil {
			panic(fmt.Errorf("failed to unmarhal user.ContactsJson: %w", err))
		}
	}
	return
}

func (entity *WithGroups) SetActiveGroups(groups []UserGroupJson) {
	if len(groups) == 0 {
		entity.GroupsJsonActive = ""
		entity.GroupsCountActive = 0
	} else {
		if data, err := ffjson.Marshal(&groups); err != nil {
			panic(err.Error())
		} else {
			entity.GroupsJsonActive = (string)(data)
			entity.GroupsCountActive = len(groups)
		}
	}
}

// Deprecated: use SplitusSpaceDbo
func (entity *WithGroups) AddGroup(group models4splitus.GroupEntry, tgBot string) (changed bool) {
	groups := entity.ActiveGroups()
	for i, g := range groups {
		if g.ID == group.ID {
			if g.Name != group.Data.Name || g.Note != group.Data.Note /*|| g.MembersCount != group.Data.MembersCount*/ {
				g.Name = group.Data.Name
				g.Note = group.Data.Note
				//g.MembersCount = group.Data.MembersCount
				groups[i] = g
				changed = true
			}
			if tgBot != "" {
				for _, b := range g.TgBots {
					if b == tgBot {
						goto found
					}
				}
				g.TgBots = append(g.TgBots, tgBot)
				changed = true
			found:
			}
			return
		}
	}
	g := UserGroupJson{
		ID:   group.ID,
		Name: group.Data.Name,
		Note: group.Data.Note,
		//MembersCount: group.Data.MembersCount,
	}
	if tgBot != "" {
		g.TgBots = []string{tgBot}
	}
	groups = append(groups, g)
	entity.SetActiveGroups(groups)
	changed = true
	return
}
