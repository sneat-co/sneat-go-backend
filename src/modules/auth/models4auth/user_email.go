package models4auth

import (
	"github.com/dal-go/dalgo/record"
	"github.com/strongo/strongoapp/appuser"
	"github.com/strongo/strongoapp/person"
	"golang.org/x/crypto/bcrypt"
	"strconv"
	"strings"
	"time"
)

const UserEmailKind = "UserEmailEntry"

type UserEmailDbo struct {
	appuser.AccountDataBase
	appuser.OwnedByUserWithID
	person.NameFields
	IsConfirmed        bool     `firestore:"isConfirmed"`
	PasswordBcryptHash []byte   `firestore:"passwordBcryptHash"`
	Providers          []string `firestore:"providers,omitempty"` // E.g. facebook, vk, user
}

type UserEmailEntry struct {
	record.DataWithID[string, *UserEmailDbo]
}

//var _ user.AccountRecord = (*UserEmailEntry)(nil)

func (userEmail UserEmailEntry) UserAccount() appuser.AccountKey {
	return appuser.AccountKey{Provider: "email", ID: userEmail.ID}
}

func (userEmail UserEmailEntry) Kind() string {
	return UserEmailKind
}

func (UserEmailEntry) NewEntity() interface{} {
	return new(UserEmailDbo)
}

func GetEmailID(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func (userEmail UserEmailEntry) GetEmail() string {
	return userEmail.ID
}

func NewUserEmailData(userID int64, isConfirmed bool, provider string) *UserEmailDbo {
	entity := &UserEmailDbo{
		OwnedByUserWithID: appuser.NewOwnedByUserWithID(strconv.FormatInt(userID, 10), time.Now()),
		IsConfirmed:       isConfirmed,
	}
	entity.AddProvider(provider)
	return entity
}

const pwdSole = "85d80e53-"

func (entity *UserEmailDbo) CheckPassword(password string) error {
	return bcrypt.CompareHashAndPassword(entity.PasswordBcryptHash, []byte(pwdSole+password))
}

func (entity *UserEmailDbo) SetPassword(password string) (err error) {
	entity.PasswordBcryptHash, err = bcrypt.GenerateFromPassword([]byte(pwdSole+password), 0)
	return
}

func (entity *UserEmailDbo) AddProvider(v string) (changed bool) {
	for _, p := range entity.Providers {
		if p == v {
			return
		}
	}
	entity.Providers = append(entity.Providers, v)
	changed = true
	return
}

//func (entity *UserEmailDbo) Load(ps []datastore.Property) error {
//	return datastore.LoadStruct(entity, ps)
//}

//func (entity *UserEmailDbo) Save() (properties []datastore.Property, err error) {
//	if properties, err = datastore.SaveStruct(entity); err != nil {
//		return
//	}
//	//return gaedb.CleanProperties(properties, map[string]gaedb.IsOkToRemove{
//	//	"DtUpdated":          gaedb.IsZeroTime,
//	//	"FirstName":          gaedb.IsEmptyString,
//	//	"LastName":           gaedb.IsEmptyString,
//	//	"NickName":           gaedb.IsEmptyString,
//	//	"PasswordBcryptHash": gaedb.IsEmptyByteArray,
//	//})
//	return nil, nil
//}
