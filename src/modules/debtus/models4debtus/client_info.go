package models4debtus

import "net/http"

type ClientInfo struct {
	UserAgent  string
	RemoteAddr string
}

func NewClientInfoFromRequest(r *http.Request) ClientInfo {
	return ClientInfo{
		UserAgent:  r.UserAgent(),
		RemoteAddr: r.RemoteAddr,
	}
}
