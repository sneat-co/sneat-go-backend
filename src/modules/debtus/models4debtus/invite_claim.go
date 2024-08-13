package models4debtus

import (
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	"reflect"
	"time"
)

const InviteClaimKind = "InviteClaim"

type InviteClaim struct {
	record.WithID[int64]
	Data *InviteClaimData
}

func NewInviteClaimWithoutID(data *InviteClaimData) InviteClaim {
	return InviteClaim{
		WithID: record.NewWithID[int64](0, NewInviteClaimIncompleteKey(), data),
		Data:   data,
	}
}

func NewInviteClaim(id int64, data *InviteClaimData) InviteClaim {
	return InviteClaim{
		WithID: record.NewWithID(id, NewInviteClaimKey(id), data),
		Data:   data,
	}
}

type InviteClaimData struct {
	InviteCode string // We don't use it as parent key as can be a bottleneck for public invites
	UserID     string
	DtClaimed  time.Time
	ClaimedOn  string // For example: "Telegram"
	ClaimedVia string // For the Telegram it would be bot name
}

func NewInviteClaimIncompleteKey() *dal.Key {
	return dal.NewIncompleteKey(InviteClaimKind, reflect.Int64, nil)
}

func NewInviteClaimKey(claimID int64) *dal.Key {
	if claimID == 0 {
		return dal.NewIncompleteKey(InviteClaimKind, reflect.Int64, nil)
	}
	return dal.NewKeyWithID(InviteClaimKind, claimID)
}

func NewInviteClaimData(inviteCode string, userID string, claimedOn, claimedVia string) *InviteClaimData {
	return &InviteClaimData{
		InviteCode: inviteCode,
		UserID:     userID,
		ClaimedOn:  claimedOn,
		ClaimedVia: claimedVia,
		DtClaimed:  time.Now(),
	}
}
