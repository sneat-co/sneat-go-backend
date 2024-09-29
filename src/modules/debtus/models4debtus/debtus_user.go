package models4debtus

import (
	"errors"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	"github.com/sneat-co/sneat-core-modules/auth/unsorted4auth"
	"github.com/sneat-co/sneat-core-modules/common4all"
	"github.com/sneat-co/sneat-core-modules/core/coremodels"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/strongo/strongoapp/appuser"
	"reflect"
	"time"
)

const AppUserKind = "User"

// Deprecated: Use NewUserKey instead
//func NewAppUserKeyOBSOLETE(appUserId string) *dal.Key {
//	if appUserId == "" {
//		return dal.NewIncompleteKey(AppUserKind, reflect.Int64, nil)
//	}
//	return dal.NewKeyWithID(AppUserKind, appUserId)
//}

// Deprecated: Use UserEntry instead
type AppUserOBSOLETE = record.DataWithID[string, *DebutsAppUserDataOBSOLETE]

func NewAppUserRecord() dal.Record {
	return dal.NewRecordWithIncompleteKey(AppUserKind, reflect.Int64, new(DebutsAppUserDataOBSOLETE))
}

// Deprecated: Use NewUserEntry instead
//func NewAppUserOBSOLETE(id string, data *DebutsAppUserDataOBSOLETE) AppUserOBSOLETE {
//	key := NewAppUserKeyOBSOLETE(id)
//	if data == nil {
//		data = new(DebutsAppUserDataOBSOLETE)
//	}
//	return record.NewDataWithID[string, *DebutsAppUserDataOBSOLETE](id, key, data)
//}

// Deprecated: Use NewUserEntries instead
//func NewAppUsersOBSOLETE(userIDs []string) []AppUserOBSOLETE {
//	users := make([]AppUserOBSOLETE, len(userIDs))
//	for i, id := range userIDs {
//		users[i] = NewAppUserOBSOLETE(id, nil)
//	}
//	return users
//}

// Deprecated: Use UserRecords instead
//func AppUserRecordsOBSOLETE(appUsers []AppUserOBSOLETE) (records []dal.Record) {
//	records = make([]dal.Record, len(appUsers))
//	for i, u := range appUsers {
//		records[i] = u.Record
//	}
//	return
//}

func NewUser(clientInfo common4all.ClientInfo) AppUserOBSOLETE {
	return AppUserOBSOLETE{
		Data: &DebutsAppUserDataOBSOLETE{
			LastUserAgent:     clientInfo.UserAgent,
			LastUserIpAddress: clientInfo.RemoteAddr,
		},
	}
}

// DebutsAppUserDataOBSOLETE is obsolete
// Should be replaced with sneat app user and DebtusSpaceDbo
type DebutsAppUserDataOBSOLETE struct { // TODO: Remove obsolete struct

	DebtusUserDbo // TODO: to be used on it's OWN
	//DebtusSpaceDbo // TODO: to be used on it's OWN

	appuser.BaseUserFields
	unsorted4auth.UserRewardBalance

	SavedCounter int `firestore:"A"` // Indexing to find most active users

	IsAnonymous        bool   `firestore:",omitempty"`
	PasswordBcryptHash []byte `firestore:",omitempty"` // TODO: Obsolete

	//dto4contactus.ContactDetails

	DtAccessGranted time.Time `firestore:",omitempty"`

	coremodels.SmsStats
	DtCreated time.Time
	appuser.WithLastLogin

	InvitedByUserID string `firestore:"invitedByUserID,omitempty"` // TODO: Prevent circular references! see users 6032980589936640 & 5998019824582656

	TelegramUserIDs    []int64 `firestore:"telegramUserIDs,omitempty"`    // TODO: Obsolete
	ViberBotID         string  `firestore:"viberBotID,omitempty"`         // TODO: Obsolete
	ViberUserID        string  `firestore:"viberUserID,omitempty"`        // TODO: Obsolete
	VkUserID           int64   `firestore:"vkUserID,omitempty"`           // TODO: Obsolete
	GoogleUniqueUserID string  `firestore:"googleUniqueUserID,omitempty"` // TODO: Obsolete
	//FbUserID           string `firestore:",omitempty"` // TODO: Obsolete Facebook assigns different IDs to same FB user for FB app & Messenger app.
	//FbmUserID          string `firestore:",omitempty"` // TODO: Obsolete So we would want to keep both IDs?
	// TODO: How do we support multiple FBM botscore? They will have different PSID (PageScopeID)

	OBSOLETE_CounterpartyIDs []int64 `firestore:"counterpartyIDs,omitempty"` // TODO: Remove obsolete

	ContactsCount int    `firestore:"contactsCount,omitempty"` // TODO: Obsolete
	ContactsJson  string `firestore:"contactsJson,omitempty"`  // TODO: Obsolete

	WithGroups
	//

	//
	//DebtCounterpartyIDs    []int64 `firestore:",omitempty"`
	//DebtCounterpartyCount  int     `firestore:",omitempty"`
	//

	dbmodels.WithLastCurrencies

	// Counts
	CountOfInvitesCreated  int `firestore:",omitempty"`
	CountOfInvitesAccepted int `firestore:",omitempty"`

	LastUserAgent     string `firestore:",omitempty"`
	LastUserIpAddress string `firestore:",omitempty"`
}

//func (v *DebutsAppUserDataOBSOLETE) GetFullName() string {
//	return v.FullName()
//}

//func userContactsByStatus(contacts []DebtusContactBrief) (active, archived []DebtusContactBrief) {
//	for _, contact := range contacts {
//		switch contact.Status {
//		case const4debtus.StatusActive:
//			contact.Status = ""
//			active = append(active, contact)
//		case const4debtus.StatusArchived:
//			contact.Status = ""
//			archived = append(archived, contact)
//		case "":
//			panic("DebtusSpaceContactEntry status is not set")
//		default:
//			panic("Unknown status: " + contact.Status)
//		}
//	}
//	return
//}

// Deprecated: Use DebtusSpaceDbo.SetContacts() instead
func (v *DebutsAppUserDataOBSOLETE) SetContacts(contacts []DebtusContactBrief) {
	panic("Use DebtusSpaceDbo.SetContacts() instead")
}

//var _ botsfwmodels.AppUserData = (*DebutsAppUserDataOBSOLETE)(nil)

//func (entity *DebutsAppUserDataOBSOLETE) Load(ps []datastore.Property) (err error) {
//	// Load I and J as usual.
//	p2 := make([]datastore.Property, 0, len(ps))
//	for _, p := range ps {
//		switch p.Name {
//		case "AA":
//			continue // Ignore legacy
//		case "FirstDueTransferOn":
//			continue // Ignore legacy
//		case "ActiveGroupsJson":
//			p.Name = "GroupsJsonActive"
//		case "ActiveGroupsCount":
//			p.Name = "GroupsCountActive"
//		case "CounterpartiesCount":
//			p.Name = "ContactsCount"
//		case "ContactsCount":
//			continue // Ignore legacy
//		case "FbUserID":
//			if v, ok := p.Value.(string); ok && v != "" {
//				entity.AddAccount(user.AuthAccount{
//					Provider: "fb",
//					ContactID:       v,
//				})
//			}
//			continue
//		case "FmbUserID":
//			if v, ok := p.Value.(string); ok && v != "" {
//				entity.AddAccount(user.AuthAccount{
//					Provider: "fbm",
//					ContactID:       v,
//				})
//			}
//			continue
//		case "FbmUserID":
//			if v, ok := p.Value.(string); ok && v != "" {
//				entity.AddAccount(user.AuthAccount{
//					Provider: "fbm",
//					ContactID:       v,
//				})
//			}
//			continue
//		case "ViberUserID":
//			continue
//		case "ViberBotID":
//			continue
//		case "TelegramUserID":
//			if telegramUserID, ok := p.Value.(int64); ok && telegramUserID != 0 {
//				entity.AccountsOfUser.Accounts = append(entity.AccountsOfUser.Accounts, "telegram::"+strconv.FormatInt(telegramUserID, 10))
//			}
//			continue
//		case "TelegramUserIDs":
//			switch p.Value.(type) {
//			case int64:
//				if id := p.Value.(int64); id != 0 {
//					entity.AccountsOfUser.Accounts = append(entity.AccountsOfUser.Accounts, "telegram::"+strconv.FormatInt(id, 10))
//				}
//			default:
//				err = fmt.Errorf("unknown type of TelegramUserIDs value: %T", p.Value)
//				return
//			}
//			continue
//		case "GoogleUniqueUserID":
//			if v, ok := p.Value.(string); ok && v != "" {
//				entity.AddAccount(user.AuthAccount{
//					Provider: "google",
//					App:      "debtusbot",
//					ContactID:       v,
//				})
//			}
//		default:
//			if p.Name == "CounterpartiesJson" {
//				p.Name = "ContactsJson"
//			}
//			if p.Name == "ContactsJson" {
//				contactsJson := p.Value.(string)
//				if contactsJson != "" {
//					entity.ContactsJson = contactsJson
//					if err := entity.FixObsolete(); err != nil {
//						return err
//					}
//					//panic(fmt.Sprintf("contactsJson: %v\n ContactsJson: %v\n ContactsJsonActive: %v", contactsJson, entity.ContactsJson, entity.ContactsJsonActive))
//					if entity.ContactsCountActive > 0 {
//						p.Name = "ContactsJsonActive"
//						p.Value = entity.ContactsJsonActive
//						p2 = append(p2, p)
//						//
//						p.Name = "ContactsCountActive"
//						p.Value = int64(entity.ContactsCountActive)
//						p2 = append(p2, p)
//					}
//
//					if entity.ContactsCountArchived > 0 {
//						p.Name = "ContactsJsonArchived"
//						p.Value = entity.ContactsJsonArchived
//						p2 = append(p2, p)
//						//
//						p.Name = "ContactsCountArchived"
//						p.Value = int64(entity.ContactsCountArchived)
//						p2 = append(p2, p)
//
//					}
//					continue
//				}
//			}
//		}
//		p2 = append(p2, p)
//	}
//	if err = datastore.LoadStruct(entity, p2); err != nil {
//		return
//	}
//	return
//}

//var userPropertiesToClean = map[string]gaedb.IsOkToRemove{
//	"AA":              gaedb.IsObsolete,
//	"FmbUserID":       gaedb.IsObsolete,
//	"CounterpartyIDs": gaedb.IsObsolete,
//	//
//	"ContactsCount": gaedb.IsZeroInt,   // TODO: Obsolete
//	"ContactsJson":  gaedb.IsEmptyJSON, // TODO: Obsolete
//	//
//	"GroupsCountActive":                   gaedb.IsZeroInt,
//	"GroupsJsonActive":                    gaedb.IsEmptyJSON,
//	"GroupsCountArchived":                 gaedb.IsZeroInt,
//	"GroupsJsonArchived":                  gaedb.IsEmptyJSON,
//	"BillsCountActive":                    gaedb.IsZeroInt,
//	"BillsJsonActive":                     gaedb.IsEmptyJSON,
//	"BillSchedulesCountActive":            gaedb.IsZeroInt,
//	"BillSchedulesJsonActive":             gaedb.IsEmptyJSON,
//	"BalanceCount":                        gaedb.IsZeroInt,
//	"BalanceJson":                         gaedb.IsEmptyString,
//	"CountOfAckTransfersByCounterparties": gaedb.IsZeroInt,
//	"CountOfAckTransfersByUser":           gaedb.IsZeroInt,
//	"CountOfInvitesAccepted":              gaedb.IsZeroInt,
//	"CountOfInvitesCreated":               gaedb.IsZeroInt,
//	"CountOfReceiptsCreated":              gaedb.IsZeroInt,
//	"CountOfTransfers":                    gaedb.IsZeroInt,
//	"ContactsCountActive":                 gaedb.IsZeroInt,
//	"ContactsJsonActive":                  gaedb.IsEmptyJSON,
//	"ContactsCountArchived":               gaedb.IsZeroInt,
//	"ContactsJsonArchived":                gaedb.IsEmptyJSON,
//	"DtAccessGranted":                     gaedb.IsZeroTime,
//	"EmailAddress":                        gaedb.IsEmptyString,
//	"EmailAddressOriginal":                gaedb.IsEmptyString,
//	"FirstName":                           gaedb.IsEmptyString,
//	"HasDueTransfers":                     gaedb.IsFalse,
//	"InvitedByUserID":                     gaedb.IsZeroInt,
//	"IsAnonymous":                         gaedb.IsFalse,
//	"LastName":                            gaedb.IsEmptyString,
//	"LastTransferAt":                      gaedb.IsZeroTime,
//	"LastTransferID":                      gaedb.IsZeroInt,
//	"LastFeedbackAt":                      gaedb.IsZeroTime,
//	"LastFeedbackRate":                    gaedb.IsEmptyString,
//	"LastUserAgent":                       gaedb.IsEmptyString,
//	"LastUserIpAddress":                   gaedb.IsEmptyString,
//	"Nickname":                            gaedb.IsEmptyString,
//	"PhoneNumber":                         gaedb.IsZeroInt,
//	"PhoneNumberConfirmed":                gaedb.IsFalse, // TODO: Duplicate of PhoneNumberIsConfirmed
//	"PhoneNumberIsConfirmed":              gaedb.IsFalse, // TODO: Duplicate of PhoneNumberConfirmed
//	"EmailConfirmed":                      gaedb.IsFalse,
//	"PreferredLanguage":                   gaedb.IsEmptyString,
//	"PrimaryCurrency":                     gaedb.IsEmptyString,
//	"ScreenName":                          gaedb.IsEmptyString,
//	"SmsCost":                             gaedb.IsZeroFloat,
//	"SmsCostUSD":                          gaedb.IsZeroInt,
//	"SmsCount":                            gaedb.IsZeroInt,
//	"Username":                            gaedb.IsEmptyString,
//	"VkUserID":                            gaedb.IsZeroInt,
//	"DtLastLogin":                         gaedb.IsZeroTime,
//	"PasswordBcryptHash":                  gaedb.IsObsolete,
//	"TransfersWithInterestCount":          gaedb.IsZeroInt,
//	//
//	"ViberBotID":         gaedb.IsObsolete,
//	"ViberUserID":        gaedb.IsObsolete,
//	"FbmUserID":          gaedb.IsObsolete,
//	"FbUserID":           gaedb.IsObsolete,
//	"FbUserIDs":          gaedb.IsObsolete,
//	"GoogleUniqueUserID": gaedb.IsObsolete,
//	"TelegramUserID":     gaedb.IsObsolete,
//	"TelegramUserIDs":    gaedb.IsObsolete,
//	//
//}

//func (entity *DebutsAppUserDataOBSOLETE) cleanProps(properties []datastore.Property) ([]datastore.Property, error) {
//	var err error
//	//if properties, err = gaedb.CleanProperties(properties, userPropertiesToClean); err != nil {
//	//	return properties, err
//	//}
//	//if properties, err = entity.UserRewardBalance.cleanProperties(properties); err != nil {
//	//	return properties, err
//	//}
//	return properties, err
//}

var ErrDuplicateContactName = errors.New("user has at least 2 contacts with same name")
var ErrDuplicateTgUserID = errors.New("user has at least 2 contacts with same TgUserID")

func (v *DebutsAppUserDataOBSOLETE) Validate() (err error) {
	//if v.GroupsJsonActive != "" && v.GroupsCountActive == 0 {
	//	return errors.New(`v.GroupsJsonActive != "" && v.GroupsCountActive == 0`)
	//}
	//
	//contacts := v.Contacts()
	//
	//if len(contacts) != v.ContactsCountActive+v.ContactsCountArchived {
	//	panic("len(contacts) != v.ContactsCountActive + v.ContactsCountArchived")
	//}
	//
	//contactsCount := len(contacts)
	//
	//contactsByName := make(map[string]int, contactsCount)
	//contactsByTgUserID := make(map[int64]int, contactsCount)
	//
	////v.TransfersWithInterestCount = 0
	//for i, contact := range contacts {
	//	if contact.ContactID == "" {
	//		panic(fmt.Sprintf("contact[%d].ContactID == 0, contact: %v, contacts: %v", i, contact, contacts))
	//	}
	//	if contact.Name == "" {
	//		panic(fmt.Sprintf("contact[%d].ContactName is Empty string, contact: %v, contacts: %v", i, contact, contacts))
	//	}
	//	if contact.Status == "" {
	//		panic(fmt.Sprintf("contact[%d].Status is Empty string, contact: %v, contacts: %v", i, contact, contacts))
	//	}
	//	{
	//		if duplicateOf, isDuplicate := contactsByName[contact.Name]; isDuplicate {
	//			err = fmt.Errorf("%w: id1=%s, id2=%s => %s", ErrDuplicateContactName, contacts[duplicateOf].ContactID, contact.ContactID, contact.Name)
	//			return
	//		}
	//		contactsByName[contact.Name] = i
	//	}
	//	if contact.TgUserID != 0 {
	//		if duplicateOf, isDuplicate := contactsByTgUserID[contact.TgUserID]; isDuplicate {
	//			err = fmt.Errorf("%s: %d for contacts %s & %s", ErrDuplicateTgUserID, contact.TgUserID, contacts[duplicateOf].ContactID, contact.ContactID)
	//			return
	//		}
	//		contactsByTgUserID[contact.TgUserID] = i
	//	}
	//	//if contact.Transfers != nil {
	//	//	v.TransfersWithInterestCount += len(contact.Transfers.OutstandingWithInterest)
	//	//}
	//}
	return
}

//func (entity *DebutsAppUserDataOBSOLETE) Save() (properties []datastore.Property, err error) {
//	if err = entity.BeforeSave(); err != nil {
//		return
//	}
//
//	//entity.SavedCounter += 1
//	if properties, err = datastore.SaveStruct(entity); err != nil {
//		return
//	}
//	if properties, err = entity.cleanProps(properties); err != nil {
//		return
//	}
//
//	//checkHasProperties(AppUserKind, properties)
//	return properties, err
//}
