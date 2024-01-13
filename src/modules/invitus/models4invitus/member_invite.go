package models4invitus

import (
	"github.com/strongo/validation"
	"time"
)

// MemberInviteBrief record
type MemberInviteBrief struct {
	ID         string    `json:"id" firestore:"id"`
	To         InviteTo  `json:"to" firestore:"to"`
	CreateTime time.Time `json:"createTime" firestore:"createTime"`
}

// Validate validates MemberInviteBrief record
func (v *MemberInviteBrief) Validate() error {
	if v.ID == "" {
		return validation.NewErrRecordIsMissingRequiredField("ID")
	}
	if v.CreateTime.IsZero() {
		return validation.NewErrRecordIsMissingRequiredField("CreateTime")
	}
	if err := v.To.Validate(); err != nil {
		return validation.NewErrBadRecordFieldValue("to", err.Error())
	}
	return nil
}
