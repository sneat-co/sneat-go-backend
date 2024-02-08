package dto

//go:generate ffjson $GOFILE

import (
	"encoding/json"
	"time"

	"github.com/crediterra/money"
	"github.com/pquerna/ffjson/ffjson"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/models"
	"github.com/strongo/decimal"
)

type UserMeDto struct {
	UserID       string
	FullName     string `json:",omitempty"`
	GoogleUserID string `json:",omitempty"`
	FbUserID     string `json:",omitempty"`
	VkUserID     int64  `json:",omitempty"`
	ViberUserID  string `json:",omitempty"`
}

type ApiAcknowledgeDto struct {
	Status   string
	UnixTime int64
}

type ApiReceiptDto struct {
	ID       string `json:"ID"`
	Code     string
	Transfer ApiReceiptTransferDto
	SentVia  string
	SentTo   string `json:",omitempty"`
}

type ApiUserDto struct {
	ID   string `json:"ID"`
	Name string `json:",omitempty"`
}

type ApiReceiptTransferDto struct {
	// TODO: We are not replacing with TransferDto as it has From/To => Creator optimisation. Think if we can reuse.
	ID             string `json:"ID"`
	Amount         money.Amount
	From           ContactDto
	DtCreated      time.Time
	To             ContactDto
	IsOutstanding  bool
	Creator        ApiUserDto
	CreatorComment string             `json:",omitempty"`
	Acknowledge    *ApiAcknowledgeDto `json:",omitempty"`
}

type ContactDto struct {
	ID     string `json:",omitempty"` // TODO: Document why it can be empty?
	UserID string `json:",omitempty"`
	Name   string `json:",omitempty"`
	// Note string `json:",omitempty"`
	Comment string `json:",omitempty"`
}

func NewContactDto(counterpartyInfo models.TransferCounterpartyInfo) ContactDto {
	dto := ContactDto{
		ID:      counterpartyInfo.ContactID,
		UserID:  counterpartyInfo.UserID,
		Name:    counterpartyInfo.Name(),
		Comment: counterpartyInfo.Comment,
	}
	if dto.Name == models.NoName {
		dto.Name = ""
	}
	return dto
}

type BillDto struct {
	// TODO: Generate ffjson
	ID      string
	Name    string
	Amount  money.Amount
	Members []BillMemberDto
}

type BillMemberDto struct {
	UserID     string `json:",omitempty"`
	ContactID  string `json:",omitempty"`
	Amount     decimal.Decimal64p2
	Paid       decimal.Decimal64p2 `json:",omitempty"`
	Share      int                 `json:",omitempty"`
	Percent    decimal.Decimal64p2 `json:",omitempty"`
	Adjustment decimal.Decimal64p2 `json:",omitempty"`
}

type ContactListDto struct {
	ContactDto
	Status  string
	Balance *json.RawMessage `json:",omitempty"`
}

type EmailInfo struct {
	Address     string
	IsConfirmed bool
}

type PhoneInfo struct {
	Number      int64
	IsConfirmed bool
}

type ContactDetailsDto struct {
	ContactListDto
	Email *EmailInfo `json:",omitempty"`
	Phone *PhoneInfo `json:",omitempty"`
	TransfersResultDto
	Groups []ContactGroupDto `json:",omitempty"`
}

type TransfersResultDto struct {
	HasMoreTransfers bool           `json:",omitempty"`
	Transfers        []*TransferDto `json:",omitempty"`
}

type TransferDto struct {
	Id            string
	Created       time.Time
	Amount        money.Amount
	IsReturn      bool
	CreatorUserID string
	From          *ContactDto
	To            *ContactDto
	Due           time.Time `json:",omitempty"`
}

func (t TransferDto) String() string {
	if b, err := ffjson.Marshal(t); err != nil {
		return err.Error()
	} else {
		return string(b)
	}
}

func TransfersToDto(userID string, transfers []models.Transfer) []*TransferDto {
	transfersDto := make([]*TransferDto, len(transfers))
	for i, transfer := range transfers {
		transfersDto[i] = TransferToDto(userID, transfer)
	}
	return transfersDto
}

type CreateTransferResponse struct {
	Error               string           `json:",omitempty"`
	Transfer            *TransferDto     `json:",omitempty"`
	CounterpartyBalance *json.RawMessage `json:",omitempty"`
	UserBalance         *json.RawMessage `json:",omitempty"`
}

func TransferToDto(userID string, transfer models.Transfer) *TransferDto {
	transferDto := TransferDto{
		Id:            transfer.ID,
		Amount:        transfer.Data.GetAmount(),
		Created:       transfer.Data.DtCreated,
		CreatorUserID: transfer.Data.CreatorUserID,
		IsReturn:      transfer.Data.IsReturn,
		Due:           transfer.Data.DtDueOn,
	}

	from := NewContactDto(*transfer.Data.From())
	to := NewContactDto(*transfer.Data.To())

	switch userID {
	case "0":
		transferDto.From = &from
		transferDto.To = &to
	case from.UserID:
		transferDto.To = &to
	case to.UserID:
		transferDto.From = &from
	default:
		transferDto.From = &from
		transferDto.To = &to
	}

	return &transferDto
}

type GroupDto struct {
	ID           string
	Name         string
	Status       string
	Note         string           `json:",omitempty"`
	MembersCount int              `json:",omitempty"`
	Members      []GroupMemberDto `json:",omitempty"`
}

type GroupMemberDto struct {
	ID        string
	UserID    string `json:",omitempty"`
	ContactID string `json:",omitempty"`
	Name      string `json:",omitempty"`
}

type ContactGroupDto struct {
	ID           string
	Name         string
	MemberID     string
	MembersCount int
}

type CounterpartyDto struct {
	Id      string
	UserID  string `json:",omitempty"`
	Name    string
	Balance *json.RawMessage `json:",omitempty"`
}
type Record struct {
	Id                     string
	Name                   string
	Counterparties         []CounterpartyDto
	Transfers              int
	CountOfReceiptsCreated int
	InvitedByUser          *struct {
		Id   string
		Name string
	} `json:",omitempty"`
	// InvitedByUserID int64 `json:",omitempty"`
	// InvitedByUserName string `json:",omitempty"`
	Balance         *json.RawMessage `json:",omitempty"`
	TelegramUserIDs []int64          `json:",omitempty"`
}
