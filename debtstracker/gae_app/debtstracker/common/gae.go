package common

import (
	"fmt"
	"google.golang.org/appengine/v2"
	"time"

	"context"
)

//func SignInt64WithExpiry(c context.Context, v int64, expires time.Time) string {
//	var toSign [16]byte
//	expiryBytes := toSign[8:16]
//	endian.PutUint64(toSign[0:8], uint64(v))
//	endian.PutUint64(expiryBytes, uint64(expires.Unix()))
//	_, signature, err := appengine.SignBytes(c, toSign[:])
//	if err != nil {
//		panic(err.Error())
//	}
//	return fmt.Sprintf("%v:%v", base64UrlEncoder.EncodeToString(expiryBytes), base64UrlEncoder.EncodeToString(signature))
//}

func SignStrWithExpiry(c context.Context, v string, expires time.Time) string {
	expiryBytes := make([]byte, 8)
	endian.PutUint64(expiryBytes, uint64(expires.Unix()))
	toSign := append([]byte(v), expiryBytes...)
	_, signature, err := appengine.SignBytes(c, toSign)
	if err != nil {
		panic(err.Error())
	}
	return fmt.Sprintf("%v:%v", base64UrlEncoder.EncodeToString(expiryBytes), base64UrlEncoder.EncodeToString(signature))
}
