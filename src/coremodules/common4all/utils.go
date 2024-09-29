package common4all

import (
	"encoding/base64"
	"encoding/binary"

	"errors"
)

var (
	endian           = binary.BigEndian
	base64UrlEncoder = base64.RawURLEncoding
	ErrEmptyID       = errors.New("Can't decode ContactID as value is empty string")
)

func EncodeID(id int64) string {
	b := make([]byte, 8)
	endian.PutUint64(b, uint64(id))
	return base64UrlEncoder.EncodeToString(b)
}

func EncodeIntID(id int) string {
	b := make([]byte, 8)
	endian.PutUint64(b, uint64(id))
	return base64UrlEncoder.EncodeToString(b)
}

func DecodeID(s string) (int64, error) {
	if s == "" {
		return 0, ErrEmptyID
	}
	b := make([]byte, 8)
	if _, err := base64UrlEncoder.Decode(b, []byte(s)); err != nil {
		return 0, err
	}
	return int64(endian.Uint64(b)), nil
}

func DecodeIntID(s string) (int, error) {
	if s == "" {
		return 0, ErrEmptyID
	}
	b := make([]byte, 8)
	if _, err := base64UrlEncoder.Decode(b, []byte(s)); err != nil {
		return 0, err
	}
	return int(endian.Uint64(b)), nil
}
