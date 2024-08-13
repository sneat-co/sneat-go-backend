package models4auth

import (
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	"strconv"
	"time"

	"errors"
)

const LoginCodeKind = "LoginCode"

// LoginCode - TODO check and describe how it is different from LoginPin
type LoginCode struct {
	record.WithID[int]
	Data *LoginCodeData
}

// LoginCodeData is a data structure for LoginCode entity.
// TODO check and describe how it is different from LoginPinData
type LoginCodeData struct {
	Created time.Time
	Claimed time.Time
	UserID  string
}

const CodeLength = 5

func NewLoginCodeKey(code int) *dal.Key {
	return dal.NewKeyWithID(LoginCodeKind, code)
}

func NewLoginCode(code int, data *LoginCodeData) *LoginCode {
	if data == nil {
		data = new(LoginCodeData)
	}
	key := NewLoginCodeKey(code)
	return &LoginCode{
		WithID: record.NewWithID(code, key, data),
		Data:   data,
	}
}

func LoginCodeToString(code int32) string {
	return fmt.Sprintf("%0"+strconv.Itoa(CodeLength)+"d", code)
}

var (
	ErrLoginCodeExpired        = errors.New("code expired")         // TODO: Show we move this to DAL?
	ErrLoginCodeAlreadyClaimed = errors.New("code already claimed") // TODO: Show we move this to DAL?
)
