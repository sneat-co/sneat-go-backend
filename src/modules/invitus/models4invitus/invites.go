package models4invitus

import (
	"fmt"
	"github.com/sneat-co/sneat-go-backend/src/modules/teamus/core4teamus"
	"github.com/sneat-co/sneat-go-core"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/sneat-co/sneat-go-core/models/dbprofile"
	"github.com/strongo/validation"
	"net/mail"
	"strings"
	"time"
)

// InviteContact holds invitation contact data
type InviteContact struct {
	Channel   string `json:"channel,omitempty" firestore:"channel,omitempty"`
	Address   string `json:"address,omitempty" firestore:"address,omitempty"`
	Title     string `json:"title,omitempty" firestore:"title,omitempty"`
	UserID    string `json:"userID,omitempty" firestore:"userID,omitempty"`
	MemberID  string `json:"memberID,omitempty" firestore:"memberID,omitempty"`
	ContactID string `json:"contactID,omitempty" firestore:"contactID,omitempty"`
}

func ValidateChannel(v string, required bool) error {
	switch v {
	case "":
		if required {
			return validation.NewErrRecordIsMissingRequiredField("channel")
		}
	case "email", "sms", "link": // known channels
	default:
		return validation.NewErrBadRecordFieldValue("channel", fmt.Sprintf("unsupported value: [%v]", v))
	}
	return nil
}

// Validate returns error if not valid
func (v InviteContact) Validate() error {
	if err := ValidateChannel(v.Channel, false); err != nil {
		return err
	}
	if v.Channel == "email" && v.Address != "" {
		if _, err := mail.ParseAddress(v.Address); err != nil {
			return validation.NewErrBadRequestFieldValue("address", fmt.Errorf("failed to parse email: %w", err).Error())
		}
	}
	return nil
}

// InviteFrom describes who created the invite
type InviteFrom struct {
	InviteContact
}

// Validate returns error if not valid
func (v InviteFrom) Validate() error {
	if v.UserID == "" {
		return validation.NewErrRecordIsMissingRequiredField("userID")
	}
	if v.MemberID == "" {
		return validation.NewErrRecordIsMissingRequiredField("memberID")
	}
	if err := v.InviteContact.Validate(); err != nil {
		return err
	}
	return nil
}

// InviteTo record
type InviteTo struct {
	InviteContact
}

// Validate returns error if not valid
func (v InviteTo) Validate() error {
	if err := v.InviteContact.Validate(); err != nil {
		return err
	}
	if v.Channel == "" {
		return validation.NewErrRecordIsMissingRequiredField("channel")
	}
	if v.Channel == "email" {
		if strings.TrimSpace(v.Address) == "" {
			return validation.NewErrRecordIsMissingRequiredField("address")
		}
		if _, err := mail.ParseAddress(v.Address); err != nil {
			return validation.NewErrBadRecordFieldValue("address", "not a valid email")
		}
	}
	const maxTitleLen = 100
	if len(v.Title) > maxTitleLen {
		return validation.NewErrBadRecordFieldValue("title",
			fmt.Sprintf("contact title should not exceed max length of %v, got: %v",
				maxTitleLen, len(v.Title)))
	}
	//if strings.TrimSpace(v.Title) == "" {
	//	return validation.NewErrRecordIsMissingRequiredField("Title")
	//}
	return nil
}

// InviteToMember an invitation to a member, member ID is validated
type InviteToMember struct {
	InviteTo
}

// Validate returns error if not valid
func (v InviteToMember) Validate() error {
	if v.MemberID == "" {
		return validation.NewErrRecordIsMissingRequiredField("memberID")
	}
	if err := v.InviteTo.Validate(); err != nil {
		return err
	}
	return nil
}

// Joiners defines fields for how many can join and how manu already joined
type Joiners struct {
	Limit  int `json:"limit" firestore:"limit"`
	Joined int `json:"joined" firestore:"joined"`
}

// Validate returns error if not valid
func (v Joiners) Validate() error {
	if v.Limit < 0 {
		return validation.NewErrBadRecordFieldValue("limit", "should be >= 0")
	}
	if v.Joined < 0 {
		return validation.NewErrBadRecordFieldValue("joined", "should be >= 0")
	}
	return nil
}

// InviteTeam a summary on team for which an invite has been created
type InviteTeam struct {
	ID    string               `json:"id,omitempty" firestore:"id,omitempty"`
	Type  core4teamus.TeamType `json:"type" firestore:"type"`
	Title string               `json:"title,omitempty" firestore:"title,omitempty"`
}

// Validate returns error if not valid
func (v InviteTeam) Validate() error {
	//if v.InviteID == "" {
	//	return validation.NewErrRecordIsMissingRequiredField("id")
	//}
	if v.Type == "" {
		return validation.NewErrRecordIsMissingRequiredField("type")
	}
	switch v.Type {
	case "family":
		// Can be empty
	default:
		if v.Title == "" {
			return validation.NewErrRecordIsMissingRequiredField("title")
		}
	}
	return nil
}

// InviteBase base data about invite to be used in InviteBrief & InviteDto
type InviteBase struct {
	Type        string     `json:"type" firestore:"type"` // either "personal" or "mass"
	Channel     string     `json:"channel" firestore:"channel"`
	ComposeOnly bool       `json:"composeOnly" firestore:"composeOnly"`
	From        InviteFrom `json:"from" firestore:"from"`
	To          *InviteTo  `json:"to" firestore:"to"`
}

// Validate returns error if not valid
func (v InviteBase) Validate() error {
	switch v.Type {
	case "":
		return validation.NewErrRecordIsMissingRequiredField("type")
	case "personal":
		if v.To == nil {
			return fmt.Errorf("%w: expected to be either 'personal' or 'mass'", validation.NewErrRecordIsMissingRequiredField("to"))
		}
		// known values
	case "mass":
		if v.To != nil {
			// TODO: we might want to change this to store a distribution channel?
			return validation.NewErrBadRecordFieldValue("to", "mass invite can not have 'to' value for now")
		}
		// known
	default:
		return validation.NewErrBadRecordFieldValue("type", "unknown invite type: "+v.Type)
	}
	if err := ValidateChannel(v.Channel, true); err != nil {
		return err
	}
	return nil
}

// InviteBrief summary about invite
type InviteBrief struct {
	ID   string      `json:"id" firestore:"id"`
	Pin  string      `json:"pin,omitempty" firestore:"pin,omitempty"`
	From *InviteFrom `json:"from,omitempty" firestore:"from,omitempty"`
	To   *InviteTo   `json:"to,omitempty" firestore:"to,omitempty"`
}

// Validate returns error if not valid
func (v InviteBrief) Validate() error {
	if v.ID == "" {
		return validation.NewErrRecordIsMissingRequiredField("id")
	}
	if err := v.From.Validate(); err != nil {
		return validation.NewErrBadRecordFieldValue("from", err.Error())
	}
	if v.To != nil {
		if err := v.To.Validate(); err != nil {
			return validation.NewErrBadRecordFieldValue("to", err.Error())
		}
	}
	return nil
}

// NewInviteBriefFromDto creates brief from DTO
func NewInviteBriefFromDto(id string, dto InviteDto) InviteBrief {
	from := dto.From
	to := *dto.To
	return InviteBrief{ID: id, From: &from, To: &to}
}

// InviteDto record - used in PersonalInviteDto and MassInvite
type InviteDto struct {
	InviteBase
	Status    string               `json:"status" firestore:"status" `
	Pin       string               `json:"pin,omitempty" firestore:"pin,omitempty"`
	TeamID    string               `json:"teamID" firestore:"teamID"`
	MessageID string               `json:"messageId" firestore:"messageId"` // e.g. email message ID from AWS SES
	CreatedAt time.Time            `json:"createdAt" firestore:"createdAt"`
	Created   dbmodels.CreatedInfo `json:"created" firestore:"created"`
	Claimed   *time.Time           `json:"claimed,omitempty" firestore:"claimed,omitempty"`
	Revoked   *time.Time           `json:"revoked" firestore:"revoked,omitempty"`
	Sending   *time.Time           `json:"sending,omitempty" firestore:"sending,omitempty"`
	Sent      *time.Time           `json:"sent,omitempty" firestore:"sent,omitempty"`
	Expires   *time.Time           `json:"expires,omitempty" firestore:"expires,omitempty"`
	Team      InviteTeam           `json:"team" firestore:"team"`
	Roles     []string             `json:"roles,omitempty" firestore:"roles,omitempty"`
	//FromUserID string     `json:"fromUserID" firestore:"fromUserID"`
	//ToUserID   string     `json:"toUserID,omitempty" firestore:"toUserID,omitempty"`
	Message string `json:"message,omitempty" firestore:"message,omitempty"`
}

// Validate validates record
func (v InviteDto) Validate() error {
	if err := v.InviteBase.Validate(); err != nil {
		return err
	}
	switch v.Status {
	case "":
		return validation.NewErrRecordIsMissingRequiredField("status")
	case "active", "accepted", "expired", "rejected": // known statuses
	default:
		return validation.NewErrBadRecordFieldValue("status", "unknown value: "+v.Status)
	}
	if v.TeamID == "" {
		return validation.NewErrRecordIsMissingRequiredField("teamID")
	}
	if err := v.Created.Validate(); err != nil {
		return validation.NewErrBadRecordFieldValue("created", err.Error())
	}
	if err := v.Team.Validate(); err != nil {
		return validation.NewErrBadRecordFieldValue("team", err.Error())
	}
	if err := v.From.Validate(); err != nil {
		return validation.NewErrBadRecordFieldValue("from", err.Error())
	}
	if v.To != nil {
		if err := v.To.Validate(); err != nil {
			return validation.NewErrBadRecordFieldValue("to", err.Error())
		}
	}
	if v.Type == "mass" && len(v.Roles) == 0 {
		return validation.NewErrRecordIsMissingRequiredField("roles")
	}
	if len(v.Roles) == 0 {
		return validation.NewErrRecordIsMissingRequiredField("roles")
	}
	for i, role := range v.Roles {
		if strings.TrimSpace(role) == "" {
			return validation.NewErrRecordIsMissingRequiredField(fmt.Sprintf("roles[%v]", i))
		}
	}
	return nil
}

func (v InviteDto) validateType(expected string) error {
	if v.Type == "" {
		return validation.NewErrRecordIsMissingRequiredField("type")
	}
	if v.Type != expected {
		return validation.NewErrBadRecordFieldValue("type", "expected to have value 'mass', got: "+expected)
	}
	return nil
}

var _ core.Validatable = (*InviteDto)(nil)

// PersonalInviteDto record
type PersonalInviteDto struct {
	InviteDto
	Address string `json:"address,omitempty" firestore:"address,omitempty"` // Can be empty for channel=link

	// in format "<TEAM_ID>:<MEMBER_ID>"
	ToTeamMemberID string `json:"toTeamMemberId" firestore:"toTeamMemberId"`

	ToAvatar *dbprofile.Avatar `json:"toAvatar,omitempty" firestore:"toAvatar,omitempty"`
	Attempts int               `json:"attempts,omitempty" firestore:"attempts,omitempty"`
}

// Validate validates record
func (v PersonalInviteDto) Validate() error {
	if err := v.InviteDto.Validate(); err != nil {
		return err
	}
	if err := v.InviteDto.validateType("personal"); err != nil {
		return err
	}
	if v.ToTeamMemberID == "" {
		return validation.NewErrRecordIsMissingRequiredField("ToTeamMemberID")
	}
	if v.ToTeamMemberID[0] == ':' {
		return validation.NewErrBadRecordFieldValue("memberID", "starts with ':'")
	}
	if v.ToTeamMemberID[len(v.ToTeamMemberID)-1] == ':' {
		return validation.NewErrBadRecordFieldValue("memberID", "ends with ':'")
	}
	switch v.Channel {
	case "":
		return validation.NewErrRecordIsMissingRequiredField("channel")
	case "email":
		if v.Address != "" {
			if address, err := mail.ParseAddress(v.Address); err != nil {
				return validation.NewErrBadRequestFieldValue("address", fmt.Errorf("field channel is 'email': %w", err).Error())
			} else if address.Name != "" {
				return validation.NewErrBadRecordFieldValue("address", "should not have name, only email address")
			} else if v.Address != strings.ToLower(v.Address) {
				return validation.NewErrBadRecordFieldValue("address", "should be in lower case")
			}
		}
	case "link", "sms":
	default:
		return validation.NewErrBadRecordFieldValue("channel", "unknown value: "+v.Channel)
	}
	if !v.ComposeOnly && v.Address == "" {
		return validation.NewErrRecordIsMissingRequiredField("address")
	}
	if v.Pin == "" {
		return validation.NewErrRecordIsMissingRequiredField("pin")
	}
	return nil
}

// MassInvite record
type MassInvite struct {
	InviteDto
	Joiners Joiners `json:"joiners" firestore:"joiners"`
}

// Validate returns error if not valid
func (v MassInvite) Validate() error {
	if err := v.InviteDto.Validate(); err != nil {
		return err
	}
	if err := v.InviteDto.validateType("mass"); err != nil {
		return err
	}
	if err := v.Joiners.Validate(); err != nil {
		return validation.NewErrBadRecordFieldValue("joiners", err.Error())
	}
	return nil
}

// InviteClaim record
type InviteClaim struct {
	Time   time.Time `json:"time" firestore:"time"`
	UserID string    `json:"userId" firestore:"userId"`
}

// InviteCode record
type InviteCode struct {
	Claims []InviteClaim `json:"claims,omitempty" firestore:"claims,omitempty"`
}
