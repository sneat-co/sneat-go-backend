package token4auth

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/strongo/logus"
	"net/http"
	"strings"
)

//func getTokenSecret() []byte { // TODO: implement getting token that is good for Firebase auth
//	return []byte("very-secret-abc")
//}

const SecretPrefix = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9." // TODO: Document purpose / intended usage

func IssueTokenLegacy(userID string, issuer string) string {
	panic("legacy code")
	//switch userID {
	//case "":
	//	panic("IssueAuthToken(userID - empty)")
	//case "0":
	//	panic("IssueAuthToken(userID == 0)")
	//}
	//
	//// Create a new token object, specifying signing method and the claims
	//// you would like it to contain.
	//token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
	//	"foo": "bar",
	//	"nbf": time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
	//})
	//
	//// Sign and get the complete encoded token as a string using the secret
	//secret := getTokenSecret()
	//tokenString, err := token.SignedString(secret)
	//if err != nil {
	//	panic(fmt.Sprintf("faield to sign: %v", err))
	//}
	////claims.SetIssuedAt(time.Now())
	////claims.SetSubject(strconv.FormatInt(userID, 10))
	////if isAdmin {
	////	claims.Set("admin", true)
	////}
	//
	//if issuer != "" {
	//	if len(issuer) > 100 {
	//		if len(issuer) <= 1000 {
	//			panic("IssueAuthToken() => len(issuer) > 20, issuer: " + issuer)
	//		} else {
	//			panic("IssueAuthToken() => len(issuer) > 20, issuer[:1000]: " + issuer[:1000])
	//		}
	//
	//	}
	//	//claims.SetIssuer(issuer)
	//}
	//
	////token := jws.NewJWT(claims, crypto.SigningMethodHS256)
	////signature, err := token.Serialize(secret)
	////if err != nil {
	////	panic(err.Error())
	////}
	//return tokenString[len(SecretPrefix):]
	////return string(signature[len(SECRET_PREFIX):])
}

type AuthInfo struct {
	UserID  string
	IsAdmin bool
	Issuer  string
}

var ErrNoToken = errors.New("No authorization token")

func Authenticate(w http.ResponseWriter, r *http.Request, required bool) (authInfo AuthInfo, token *jwt.Token, err error) {
	c := r.Context()
	s := r.URL.Query().Get("secret")
	if s == "" {
		if a := r.Header.Get("Authorization"); strings.HasPrefix(a, "Bearer ") {
			s = a[7:]
		}
	}

	defer func() {
		if err != nil && required {
			w.WriteHeader(http.StatusUnauthorized)
			w.Header().Add("Access-Control-Allow-Origin", "*")
			_, _ = w.Write([]byte(err.Error()))
		}
	}()

	if s == "" {
		err = ErrNoToken
		return
	}

	if strings.Count(s, ".") == 1 {
		s = SecretPrefix + s
	}

	logus.Debugf(r.Context(), "JWT token: [%v]", s)

	if token, err = jwt.Parse(s, func(token *jwt.Token) (interface{}, error) {
		return nil, nil
	}); err != nil {
		logus.Debugf(c, "Tried to parse: [%v]", s)
		return
	}

	if !token.Valid {
		err = fmt.Errorf("invalid token: %v", token)
		return
	}
	if claims, ok := token.Claims.(SneatClaims); ok {
		if claims.Issuer != "" {
			authInfo.Issuer = claims.Issuer
		} else {
			err = errors.New("token is missing 'issuer' claim")
			return
		}
		authInfo.UserID = claims.Subject
		authInfo.IsAdmin = claims.Admin
	}
	return
}

type SneatClaims struct {
	jwt.RegisteredClaims
	Admin bool `json:"admin"`
}
