package common4debtus

//var store = sessions.NewCookieStore([]byte("Very-secret: 8468df92-fe04-432d-9f56-55ca9ffc20a3"))
//
//func GetSession(r *http.Request) (*DTSession, error) {
//	s, err := store.Get(r, "session")
//	if err != nil {
//		logus.Warningf(appengine.NewContext(r), err.Error())
//	}
//	return &DTSession{Session: s}, nil
//}
//
//type DTSession struct {
//	*sessions.Session
//}

//const (
//	session_PARAM_USER_ID = "UserID"
//	cookie_USER_ID = "UserID"
//)
//
//func (s *DTSession) UserID(w http.ResponseWriter, r http.Request) int64 {
//	if uid, ok := s.Values[session_PARAM_USER_ID]; ok {
//		if cookie, err := r.Cookie(cookie_USER_ID); err == nil {
//			if len(cookie.Value) > 0 {
//				uidStr := strconv.FormatInt(uid, 10)
//				if cookie.Value != uidStr {
//					logus.Warningf(appengine.NewContext(r), "Cookie(UserID):%v != Session(UserID):%v", cookie.Value, uid)
//					s.setUserIdCookie(w, uid)
//				}
//			} else {
//				s.setUserIdCookie(w, uid)
//			}
//		}
//		return uid.(int64)
//	}
//
//	return 0
//}
//
//func (s *DTSession) setUserIdCookie(w http.ResponseWriter, id int64) {
//	http.SetCookie(w, http.Cookie{Name: cookie_USER_ID, Value: strconv.FormatInt(id, 10)})
//}
//
//func (s *DTSession) SetUserID(v int64, w http.ResponseWriter)  {
//	s.Values[session_PARAM_USER_ID] = v
//	s.setUserIdCookie(w, v)
//}

//func CreateSecret(transferID, userID int64) (string) {
//	toSign := "54db6494-381f-4b04-b748-a5f90c274cfd:" + fmt.Sprintf("%s:%s", EncodeID(transferID), EncodeID(userID))
//	signature := sha1.Sum([]byte(toSign))
//	return base64.RawURLEncoding.EncodeToString(signature[:])
//}
