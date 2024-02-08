package auth

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/strongo/log"
	"google.golang.org/appengine/v2"
	"net/http"
	"strings"
	"time"
)

var secret = []byte("very-secret-abc")

const SECRET_PREFIX = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9."

func IssueToken(userID string, issuer string, isAdmin bool) string {
	switch userID {
	case "":
		panic("IssueToken(userID - empty)")
	case "0":
		panic("IssueToken(userID == 0)")
	}

	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"foo": "bar",
		"nbf": time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(secret)
	if err != nil {
		panic(fmt.Sprintf("faield to sign: %v", err))
	}
	//claims.SetIssuedAt(time.Now())
	//claims.SetSubject(strconv.FormatInt(userID, 10))
	//if isAdmin {
	//	claims.Set("admin", true)
	//}

	if issuer != "" {
		if len(issuer) > 100 {
			if len(issuer) <= 1000 {
				panic("IssueToken() => len(issuer) > 20, issuer: " + issuer)
			} else {
				panic("IssueToken() => len(issuer) > 20, issuer[:1000]: " + issuer[:1000])
			}

		}
		//claims.SetIssuer(issuer)
	}

	//token := jws.NewJWT(claims, crypto.SigningMethodHS256)
	//signature, err := token.Serialize(secret)
	//if err != nil {
	//	panic(err.Error())
	//}
	return tokenString[len(SECRET_PREFIX):]
	//return string(signature[len(SECRET_PREFIX):])
}

type AuthInfo struct {
	UserID  string
	IsAdmin bool
	Issuer  string
}

var ErrNoToken = errors.New("No authorization token")

func Authenticate(w http.ResponseWriter, r *http.Request, required bool) (authInfo AuthInfo, token *jwt.Token, err error) {
	c := appengine.NewContext(r)
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
		s = SECRET_PREFIX + s
	}

	log.Debugf(appengine.NewContext(r), "JWT token: [%v]", s)

	if token, err = jwt.Parse(s, func(token *jwt.Token) (interface{}, error) {
		return nil, nil
	}); err != nil {
		log.Debugf(c, "Tried to parse: [%v]", s)
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
