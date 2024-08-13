package dtdal

import (
	"context"
	"errors"
	"fmt"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/crediterra/money"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/auth/facade4auth"
	"github.com/sneat-co/sneat-go-backend/src/auth/models4auth"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/dto4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/models4debtus"
	"github.com/sneat-co/sneat-go-backend/src/modules/splitus/models4splitus"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/dbo4userus"
	"github.com/strongo/decimal"
	"github.com/strongo/gotwilio"
	"github.com/strongo/strongoapp"
	"math/rand"
	"net/http"
	"regexp"
	"time"
)

type TransferSource interface {
	PopulateTransfer(t *models4debtus.TransferData)
}

const (
	AckAccept  = "accept"
	AckDecline = "decline"
)

//var (
//	CrossGroupTransaction  = dal.CrossGroupTransaction
//	SingleGroupTransaction = db.SingleGroupTransaction
//)

type TransferReturnUpdate struct {
	TransferID     string
	ReturnedAmount decimal.Decimal64p2
}

type RewardDal interface {
	//GetRewardByID(c context.Context, rewardID int64) (reward models.Reward, err error)
	InsertReward(c context.Context, tx dal.ReadwriteTransaction, rewardEntity *models4debtus.RewardDbo) (reward models4debtus.Reward, err error)
}

type TransferDal interface {
	GetTransfersByID(c context.Context, tx dal.ReadSession, transferIDs []string) ([]models4debtus.TransferEntry, error)
	LoadTransfersByUserID(c context.Context, userID string, offset, limit int) (transfers []models4debtus.TransferEntry, hasMore bool, err error)
	LoadTransfersByContactID(c context.Context, contactID string, offset, limit int) (transfers []models4debtus.TransferEntry, hasMore bool, err error)
	LoadTransferIDsByContactID(c context.Context, contactID string, limit int, startCursor string) (transferIDs []string, endCursor string, err error)
	LoadOverdueTransfers(c context.Context, tx dal.ReadSession, userID string, limit int) (transfers []models4debtus.TransferEntry, err error)
	LoadOutstandingTransfers(c context.Context, tx dal.ReadSession, periodEnds time.Time, userID, contactID string, currency money.CurrencyCode, direction models4debtus.TransferDirection) (transfers []models4debtus.TransferEntry, err error)
	LoadDueTransfers(c context.Context, tx dal.ReadSession, userID string, limit int) (transfers []models4debtus.TransferEntry, err error)
	LoadLatestTransfers(c context.Context, offset, limit int) ([]models4debtus.TransferEntry, error)
	DelayUpdateTransferWithCreatorReceiptTgMessageID(c context.Context, botCode string, transferID string, creatorTgChatID, creatorTgReceiptMessageID int64) error
	DelayUpdateTransfersWithCounterparty(c context.Context, creatorCounterpartyID, counterpartyCounterpartyID string) error
	DelayUpdateTransfersOnReturn(c context.Context, returnTransferID string, transferReturnUpdates []TransferReturnUpdate) (err error)
}

type ReceiptDal interface {
	UpdateReceipt(c context.Context, tx dal.ReadwriteTransaction, receipt models4debtus.ReceiptEntry) error
	GetReceiptByID(c context.Context, tx dal.ReadSession, id string) (models4debtus.ReceiptEntry, error)
	MarkReceiptAsSent(c context.Context, receiptID, transferID string, sentTime time.Time) error
	CreateReceipt(c context.Context, data *models4debtus.ReceiptDbo) (receipt models4debtus.ReceiptEntry, err error)
	DelayedMarkReceiptAsSent(c context.Context, receiptID, transferID string, sentTime time.Time) error
	DelayCreateAndSendReceiptToCounterpartyByTelegram(c context.Context, env string, transferID string, userID string) error
}

var ErrReminderAlreadyRescheduled = errors.New("reminder already rescheduled")

type ReminderDal interface {
	DelayDiscardReminders(c context.Context, transferIDs []string, returnTransferID string) error
	DelayCreateReminderForTransferUser(c context.Context, transferID string, userID string) error
	SaveReminder(c context.Context, tx dal.ReadwriteTransaction, reminder models4debtus.Reminder) (err error)
	GetReminderByID(c context.Context, tx dal.ReadSession, id string) (models4debtus.Reminder, error)
	RescheduleReminder(c context.Context, reminderID string, remindInDuration time.Duration) (oldReminder, newReminder models4debtus.Reminder, err error)
	SetReminderStatus(c context.Context, reminderID string, returnTransferID string, status string, when time.Time) (reminder models4debtus.Reminder, err error)
	DelaySetReminderIsSent(c context.Context, reminderID string, sentAt time.Time, messageIntID int64, messageStrID, locale, errDetails string) error
	SetReminderIsSent(c context.Context, reminderID string, sentAt time.Time, messageIntID int64, messageStrID, locale, errDetails string) error
	SetReminderIsSentInTransaction(c context.Context, tx dal.ReadwriteTransaction, reminder models4debtus.Reminder, sentAt time.Time, messageIntID int64, messageStrID, locale, errDetails string) (err error)
	GetActiveReminderIDsByTransferID(c context.Context, tx dal.ReadSession, transferID int) ([]int, error)
	GetSentReminderIDsByTransferID(c context.Context, tx dal.ReadSession, transferID int) ([]int, error)
}

type EmailDal interface {
	InsertEmail(c context.Context, tx dal.ReadwriteTransaction, entity *models4auth.EmailData) (models4auth.Email, error)
	UpdateEmail(c context.Context, tx dal.ReadwriteTransaction, email models4auth.Email) error
	GetEmailByID(c context.Context, tx dal.ReadSession, id int64) (models4auth.Email, error)
}

type FeedbackDal interface {
	GetFeedbackByID(c context.Context, tx dal.ReadSession, feedbackID int64) (feedback models4debtus.Feedback, err error)
}

type ContactDal interface {
	GetLatestContacts(whc botsfw.WebhookContext, tx dal.ReadSession, spaceID string, limit, totalCount int) (contacts []models4debtus.DebtusSpaceContactEntry, err error)
	InsertContact(c context.Context, tx dal.ReadwriteTransaction, contactEntity *models4debtus.DebtusSpaceContactDbo) (contact models4debtus.DebtusSpaceContactEntry, err error)
	GetContactIDsByTitle(c context.Context, tx dal.ReadSession, spaceID, userID string, title string, caseSensitive bool) (contactIDs []string, err error)
	GetContactsWithDebts(c context.Context, tx dal.ReadSession, spaceID, userID string) (contacts []models4debtus.DebtusSpaceContactEntry, err error)
}

type BillsHolderGetter func(c context.Context) (billsHolder dal.Record, err error)

type SplitDal interface {
	GetSplitByID(c context.Context, splitID int) (split models4splitus.Split, err error)
	InsertSplit(c context.Context, splitEntity models4splitus.SplitEntity) (split models4splitus.Split, err error)
}

type TgGroupDal interface {
	GetTgGroupByID(c context.Context, tx dal.ReadSession, id int64) (tgGroup models4auth.TgGroup, err error)
	SaveTgGroup(c context.Context, tx dal.ReadwriteTransaction, tgGroup models4auth.TgGroup) (err error)
}

type BillScheduleDal interface {
	GetBillScheduleByID(c context.Context, id int64) (billSchedule models4splitus.BillSchedule, err error)
	InsertBillSchedule(c context.Context, billScheduleEntity *models4splitus.BillScheduleEntity) (billSchedule models4splitus.BillSchedule, err error)
	UpdateBillSchedule(c context.Context, billSchedule models4splitus.BillSchedule) (err error)
}

//type GroupMemberDal interface {
//	GetGroupMemberByID(c context.Context, groupMemberID int64) (groupMember models.GroupMember, err error)
//	CreateGroupMember(c context.Context, groupMemberEntity *models.GroupMemberData) (groupMember models.GroupMember, err error)
//}

type TwilioDal interface {
	GetLastTwilioSmsesForUser(c context.Context, tx dal.ReadSession, userID string, to string, limit int) (result []models4debtus.TwilioSms, err error)
	SaveTwilioSms(
		c context.Context,
		smsResponse *gotwilio.SmsResponse,
		transfer models4debtus.TransferEntry,
		phoneContact dto4contactus.PhoneContact,
		userID string,
		tgChatID int64,
		smsStatusMessageID int,
	) (twiliosSms models4debtus.TwilioSms, err error)
}

const LetterBytes = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789" // Removed 1, I and 0, O as can be messed with l/1 and 0.
var InviteCodeRegex = regexp.MustCompile(fmt.Sprintf("[%v]+", LetterBytes))

var random = rand.New(rand.NewSource(time.Now().UnixNano()))

func RandomCode(n uint8) string {
	b := make([]byte, n)
	lettersCount := len(LetterBytes)
	for i := range b {
		b[i] = LetterBytes[random.Intn(lettersCount)]
	}
	return string(b)
}

type InviteDal interface {
	GetInvite(c context.Context, tx dal.ReadSession, inviteCode string) (models4debtus.Invite, error)
	ClaimInvite(c context.Context, userID string, inviteCode, claimedOn, claimedVia string) (err error)
	ClaimInvite2(c context.Context, inviteCode string, invite models4debtus.Invite, claimedByUserID string, claimedOn, claimedVia string) (err error)
	CreatePersonalInvite(ec strongoapp.ExecutionContext, userID string, inviteBy models4debtus.InviteBy, inviteToAddress, createdOnPlatform, createdOnID, related string) (models4debtus.Invite, error)
	CreateMassInvite(ec strongoapp.ExecutionContext, userID string, inviteCode string, maxClaimsCount int32, createdOnPlatform string) (invite models4debtus.Invite, err error)
}

type AdminDal interface {
	DeleteAll(c context.Context, botCode, botChatID string) error
	LatestUsers(c context.Context) (users []dbo4userus.UserEntry, err error)
}

//type TaskQueueDal interface {
//	CallDelayFunc(c context.Context, queueName, subPath, key string, f interface{}, args ...interface{}) error
//}

var (
	DB      dal.DB
	Contact ContactDal

	UserGoogle facade4auth.UserGoogleDal

	Feedback FeedbackDal
	//Bill      BillDal
	Receipt   ReceiptDal
	Reminder  ReminderDal
	TgGroup   TgGroupDal
	Transfer  TransferDal
	LoginPin  facade4auth.LoginPinDal
	LoginCode facade4auth.LoginCodeDal
	Twilio    TwilioDal
	Invite    InviteDal
	Admin     AdminDal

	HttpClient  func(c context.Context) *http.Client
	BotHost     botsfw.BotHost
	HttpAppHost strongoapp.HttpAppHost

	//Split        SplitDal
	//BillSchedule BillScheduleDal
	//Reward RewardDal
	//TaskQueue		   TaskQueueDal
	//UserVk         UserVkDal

)

func InsertWithRandomStringID(c context.Context, tx dal.ReadwriteTransaction, record dal.Record) error {
	_, _, _ = c, tx, record
	return errors.New("TODO: use dalgo")
}
