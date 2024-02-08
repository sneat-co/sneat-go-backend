package models

import (
	"github.com/dal-go/dalgo/record"
	"github.com/strongo/strongoapp/appuser"
	"github.com/strongo/strongoapp/person"
	"golang.org/x/crypto/bcrypt"
	"strconv"
	"strings"
	"time"
)

const UserEmailKind = "UserEmail"

type UserEmailData struct {
	appuser.AccountDataBase
	appuser.OwnedByUserWithID
	person.NameFields
	IsConfirmed        bool
	PasswordBcryptHash []byte   `datastore:",noindex"`
	Providers          []string `datastore:",noindex"` // E.g. facebook, vk, user
}

type UserEmail struct {
	record.WithID[string]
	Data *UserEmailData
}

//var _ user.AccountRecord = (*UserEmail)(nil)

func (userEmail UserEmail) UserAccount() appuser.AccountKey {
	return appuser.AccountKey{Provider: "email", ID: userEmail.ID}
}

func (userEmail UserEmail) Kind() string {
	return UserEmailKind
}

func (UserEmail) NewEntity() interface{} {
	return new(UserEmailData)
}

func GetEmailID(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func (userEmail UserEmail) GetEmail() string {
	return userEmail.ID
}

func NewUserEmailData(userID int64, isConfirmed bool, provider string) *UserEmailData {
	entity := &UserEmailData{
		OwnedByUserWithID: appuser.NewOwnedByUserWithID(strconv.FormatInt(userID, 10), time.Now()),
		IsConfirmed:       isConfirmed,
	}
	entity.AddProvider(provider)
	return entity
}

const pwdSole = "85d80e53-"

func (entity *UserEmailData) CheckPassword(password string) error {
	return bcrypt.CompareHashAndPassword(entity.PasswordBcryptHash, []byte(pwdSole+password))
}

func (entity *UserEmailData) SetPassword(password string) (err error) {
	entity.PasswordBcryptHash, err = bcrypt.GenerateFromPassword([]byte(pwdSole+password), 0)
	return
}

func (entity *UserEmailData) AddProvider(v string) (changed bool) {
	for _, p := range entity.Providers {
		if p == v {
			return
		}
	}
	entity.Providers = append(entity.Providers, v)
	changed = true
	return
}

//func (entity *UserEmailData) Load(ps []datastore.Property) error {
//	return datastore.LoadStruct(entity, ps)
//}

//func (entity *UserEmailData) Save() (properties []datastore.Property, err error) {
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
