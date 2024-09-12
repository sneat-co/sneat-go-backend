package facade4anybot

import "errors"

var (
	ErrBadInput               = errors.New("bad data input")
	ErrInvalidAcknowledgeType = errors.New("invalid acknowledge type")
	ErrSelfAcknowledgement    = errors.New("transfer not allowed to be accepted by creator")
	ErrLoginAlreadySigned     = errors.New("this login code already used to sign in")
	ErrLoginExpired           = errors.New("this login code has expired")
)
