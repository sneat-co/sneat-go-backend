package models4splitus

import (
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/const4debtus"
	"github.com/sneat-co/sneat-go-backend/src/modules/splitus/briefs4splitus"
	"github.com/strongo/validation"
	"strings"

	"github.com/crediterra/money"
	"github.com/pquerna/ffjson/ffjson"
)

const GroupKind = "Group"

type GroupEntry = record.DataWithID[string, *GroupDbo]

func NewGroupEntry(id string, data *GroupDbo) GroupEntry {
	key := NewGroupKey(id)
	if data == nil {
		data = new(GroupDbo)
	}
	return GroupEntry{
		WithID: record.WithID[string]{
			ID:     id,
			Key:    key,
			Record: dal.NewRecordWithData(key, data),
		},
		Data: data,
	}
}

func NewGroupKey(id string) *dal.Key {
	if id == "" {
		key, err := dal.NewKeyWithOptions(GroupKind, dal.WithRandomStringID(dal.RandomLength(const4debtus.GroupIdLen)))
		if err != nil {
			panic(err.Error())
		}
		return key
	}
	return dal.NewKeyWithID(GroupKind, id)
}

type GroupDbo struct {
	CreatorUserID string
	//IsUser2User         bool   `datastore:",noindex"`
	Name            string             `firestore:"name"`
	Note            string             `firestore:"note,omitempty"`
	DefaultCurrency money.CurrencyCode `firestore:"defaultCurrency,omitempty"`
	//
	telegramGroups      []briefs4splitus.GroupTgChatJson
	TelegramGroupsCount int    `firestore:"TgGroupsCount,omitempty"`
	TelegramGroupsJson  string `firestore:"TgGroupsJson,omitempty"`
	//
	BillsHolder
}

func (entity *GroupDbo) GetTelegramGroups() (tgGroups []briefs4splitus.GroupTgChatJson, err error) {
	if entity.telegramGroups != nil {
		return entity.telegramGroups, nil
	}
	if entity.TelegramGroupsJson != "" {
		if err = ffjson.Unmarshal([]byte(entity.TelegramGroupsJson), &tgGroups); err != nil {
			return
		} else if len(tgGroups) != entity.TelegramGroupsCount {
			err = fmt.Errorf("len([]GroupTgChatJson) != entity.TelegramGroupsCount: %d != %d", len(tgGroups), entity.TelegramGroupsCount)
			return
		}
		entity.telegramGroups = tgGroups
	}
	return
}

func (entity *GroupDbo) SetTelegramGroups(tgGroups []briefs4splitus.GroupTgChatJson) (changed bool) {
	if data, err := ffjson.Marshal(tgGroups); err != nil {
		panic(err.Error())
	} else {
		if s := string(data); s != entity.TelegramGroupsJson {
			entity.TelegramGroupsJson = s
			changed = true
		}
		if l := len(tgGroups); l != entity.TelegramGroupsCount {
			entity.TelegramGroupsCount = l
			changed = true
		}
	}
	return
}

//func (entity *GroupDbo) Load(ps []datastore.Property) (err error) {
//	if ps, err = gaedb.CleanProperties(ps, map[string]gaedb.IsOkToRemove{
//		"Status": gaedb.IsObsolete,
//	}); err != nil {
//		return
//	}
//	if err = datastore.LoadStruct(entity, ps); err != nil {
//		return err
//	}
//	return nil
//}

//var groupPropertiesToClean = map[string]gaedb.IsOkToRemove{
//	"DefaultCurrency":       gaedb.IsEmptyString,
//	"MembersCount":          gaedb.IsZeroInt,
//	"MemberLastID":          gaedb.IsZeroInt,
//	"MembersJson":           gaedb.IsEmptyJSON,
//	"Note":                  gaedb.IsEmptyString,
//	"OutstandingBillsJson":  gaedb.IsEmptyJSON,
//	"OutstandingBillsCount": gaedb.IsZeroInt,
//	"TgGroupsCount":         gaedb.IsZeroInt,
//	"TgGroupsJson":          gaedb.IsEmptyJSON,
//}

func (entity *GroupDbo) Validate() error {
	if entity.CreatorUserID == "" {
		return validation.NewErrRecordIsMissingRequiredField("CreatorUserID")
	}
	if strings.TrimSpace(entity.Name) == "" {
		return validation.NewErrRecordIsMissingRequiredField("Name")
	}
	//if err := entity.validateMembers(entity.GetGroupMembers(), entity.MembersCount); err != nil {
	//	return err
	//}
	//ps, err := datastore.SaveStruct(entity)
	//if ps, err = gaedb.CleanProperties(ps, groupPropertiesToClean); err != nil {
	//	return ps, fmt.Errorf("%w: failed to clean properties for *GroupDbo", err)
	//}
	//checkHasProperties(GroupKind, ps)
	return nil
}
