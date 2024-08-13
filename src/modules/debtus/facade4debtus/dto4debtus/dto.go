package dto4debtus

//go:generate ffjson $GOFILE

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/dto4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/models4debtus"
	"time"

	"github.com/crediterra/money"
	"github.com/pquerna/ffjson/ffjson"
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
	ID       string `json:"ContactID"`
	Code     string
	Transfer ApiReceiptTransferDto
	SentVia  string
	SentTo   string `json:",omitempty"`
}

type ApiUserDto struct {
	ID   string `json:"ContactID"`
	Name string `json:",omitempty"`
}

type ApiReceiptTransferDto struct {
	// TODO: We are not replacing with TransferDto as it has From/To => Creator optimisation. Think if we can reuse.
	ID             string `json:"ContactID"`
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

func NewContactDto(counterpartyInfo models4debtus.TransferCounterpartyInfo) ContactDto {
	dto := ContactDto{
		ID:      counterpartyInfo.ContactID,
		UserID:  counterpartyInfo.UserID,
		Name:    counterpartyInfo.Name(),
		Comment: counterpartyInfo.Comment,
	}
	if dto.Name == dto4contactus.NoName {
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
	Status  string        `json:"status"`
	Balance money.Balance `json:"balance,omitempty"`
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

func TransfersToDto(userID string, transfers []models4debtus.TransferEntry) []*TransferDto {
	transfersDto := make([]*TransferDto, len(transfers))
	for i, transfer := range transfers {
		transfersDto[i] = TransferToDto(userID, transfer)
	}
	return transfersDto
}

type CreateTransferResponse struct {
	Error               string        `json:",omitempty"`
	Transfer            *TransferDto  `json:",omitempty"`
	CounterpartyBalance money.Balance `json:",omitempty"`
	UserBalance         money.Balance `json:",omitempty"`
}

func TransferToDto(userID string, transfer models4debtus.TransferEntry) *TransferDto {
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
	ID           string           `json:"id"`
	Name         string           `json:"name"`
	Status       string           `json:"status"`
	Note         string           `json:"note,omitempty"`
	MembersCount int              `json:"membersCount,omitempty"`
	Members      []GroupMemberDto `json:"members,omitempty"`
}

type GroupMemberDto struct {
	ID        string `json:"id"`
	UserID    string `json:"userID,omitempty"`
	ContactID string `json:"contactID,omitempty"`
	Name      string `json:"name,omitempty"`
}

type ContactGroupDto struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	MemberID     string `json:"memberID"`
	MembersCount int    `json:"membersCount"`
}

type CounterpartyDto struct {
	ID      string        `json:"id"`
	UserID  string        `json:",omitempty"`
	Name    string        `json:"name"`
	Balance money.Balance `json:"balance,omitempty"`
}

type Record struct {
	Id                     string `json:"id"`
	Name                   string `json:"name"`
	Counterparties         []CounterpartyDto
	Transfers              int `json:"transfers,omitempty"`
	CountOfReceiptsCreated int `json:"countOfReceiptsCreated,omitempty"`
	InvitedByUser          *struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	} `json:"invitedByUser,omitempty"`
	// InvitedByUserID int64 `json:",omitempty"`
	// InvitedByUserName string `json:",omitempty"`
	Balance         money.Balance `json:"balance,omitempty"`
	TelegramUserIDs []int64       `json:"telegramUserIDs,omitempty"`
}
