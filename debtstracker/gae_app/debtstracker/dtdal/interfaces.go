package dtdal

import (
	"context"
	"errors"
	"fmt"
	"github.com/bots-go-framework/bots-fw-telegram-models/botsfwtgmodels"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/crediterra/money"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/auth"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/models"
	"github.com/strongo/decimal"
	"github.com/strongo/gotwilio"
	"github.com/strongo/strongoapp"
	"math/rand"
	"net/http"
	"regexp"
	"sync"
	"time"
)

type TransferSource interface {
	PopulateTransfer(t *models.TransferData)
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
	InsertReward(c context.Context, tx dal.ReadwriteTransaction, rewardEntity *models.RewardDbo) (reward models.Reward, err error)
}

type TransferDal interface {
	GetTransfersByID(c context.Context, tx dal.ReadSession, transferIDs []string) ([]models.TransferEntry, error)
	LoadTransfersByUserID(c context.Context, userID string, offset, limit int) (transfers []models.TransferEntry, hasMore bool, err error)
	LoadTransfersByContactID(c context.Context, contactID string, offset, limit int) (transfers []models.TransferEntry, hasMore bool, err error)
	LoadTransferIDsByContactID(c context.Context, contactID string, limit int, startCursor string) (transferIDs []string, endCursor string, err error)
	LoadOverdueTransfers(c context.Context, tx dal.ReadSession, userID string, limit int) (transfers []models.TransferEntry, err error)
	LoadOutstandingTransfers(c context.Context, tx dal.ReadSession, periodEnds time.Time, userID, contactID string, currency money.CurrencyCode, direction models.TransferDirection) (transfers []models.TransferEntry, err error)
	LoadDueTransfers(c context.Context, tx dal.ReadSession, userID string, limit int) (transfers []models.TransferEntry, err error)
	LoadLatestTransfers(c context.Context, offset, limit int) ([]models.TransferEntry, error)
	DelayUpdateTransferWithCreatorReceiptTgMessageID(c context.Context, botCode string, transferID string, creatorTgChatID, creatorTgReceiptMessageID int64) error
	DelayUpdateTransfersWithCounterparty(c context.Context, creatorCounterpartyID, counterpartyCounterpartyID string) error
	DelayUpdateTransfersOnReturn(c context.Context, returnTransferID string, transferReturnUpdates []TransferReturnUpdate) (err error)
}

type ReceiptDal interface {
	UpdateReceipt(c context.Context, tx dal.ReadwriteTransaction, receipt models.Receipt) error
	GetReceiptByID(c context.Context, tx dal.ReadSession, id string) (models.Receipt, error)
	MarkReceiptAsSent(c context.Context, receiptID, transferID string, sentTime time.Time) error
	CreateReceipt(c context.Context, data *models.ReceiptData) (receipt models.Receipt, err error)
	DelayedMarkReceiptAsSent(c context.Context, receiptID, transferID string, sentTime time.Time) error
	DelayCreateAndSendReceiptToCounterpartyByTelegram(c context.Context, env string, transferID string, userID string) error
}

var ErrReminderAlreadyRescheduled = errors.New("reminder already rescheduled")

type ReminderDal interface {
	DelayDiscardReminders(c context.Context, transferIDs []string, returnTransferID string) error
	DelayCreateReminderForTransferUser(c context.Context, transferID string, userID string) error
	SaveReminder(c context.Context, tx dal.ReadwriteTransaction, reminder models.Reminder) (err error)
	GetReminderByID(c context.Context, tx dal.ReadSession, id string) (models.Reminder, error)
	RescheduleReminder(c context.Context, reminderID string, remindInDuration time.Duration) (oldReminder, newReminder models.Reminder, err error)
	SetReminderStatus(c context.Context, reminderID string, returnTransferID string, status string, when time.Time) (reminder models.Reminder, err error)
	DelaySetReminderIsSent(c context.Context, reminderID string, sentAt time.Time, messageIntID int64, messageStrID, locale, errDetails string) error
	SetReminderIsSent(c context.Context, reminderID string, sentAt time.Time, messageIntID int64, messageStrID, locale, errDetails string) error
	SetReminderIsSentInTransaction(c context.Context, tx dal.ReadwriteTransaction, reminder models.Reminder, sentAt time.Time, messageIntID int64, messageStrID, locale, errDetails string) (err error)
	GetActiveReminderIDsByTransferID(c context.Context, tx dal.ReadSession, transferID int) ([]int, error)
	GetSentReminderIDsByTransferID(c context.Context, tx dal.ReadSession, transferID int) ([]int, error)
}

type CreateUserData struct {
	//FbUserID     string
	//GoogleUserID string
	//VkUserID     int64
	FirstName  string
	LastName   string
	ScreenName string
	Nickname   string
}

func CreateUserEntity(createUserData CreateUserData) (user *models.DebutsAppUserDataOBSOLETE) {
	return &models.DebutsAppUserDataOBSOLETE{
		//FbUserID: createUserData.FbUserID,
		//VkUserID: createUserData.VkUserID,
		//GoogleUniqueUserID: createUserData.GoogleUserID,
		ContactDetails: models.ContactDetails{
			FirstName:  createUserData.FirstName,
			LastName:   createUserData.LastName,
			ScreenName: createUserData.ScreenName,
			Nickname:   createUserData.Nickname,
		},
	}
}

type UserDal interface {
	GetUserByStrID(c context.Context, userID string) (models.AppUser, error)
	GetUserByVkUserID(c context.Context, vkUserID int64) (models.AppUser, error)
	CreateAnonymousUser(c context.Context) (models.AppUser, error)
	CreateUser(c context.Context, userEntity *models.DebutsAppUserDataOBSOLETE) (models.AppUser, error)
	DelaySetUserPreferredLocale(c context.Context, delay time.Duration, userID string, localeCode5 string) error
	DelayUpdateUserHasDueTransfers(c context.Context, userID string) error
	SetLastCurrency(c context.Context, userID string, currency money.CurrencyCode) error
	DelayUpdateUserWithBill(c context.Context, userID string, billID string) error
	DelayUpdateUserWithContact(c context.Context, userID, contactID string) error
}

type PasswordResetDal interface {
	GetPasswordResetByID(c context.Context, tx dal.ReadSession, id int) (models.PasswordReset, error)
	CreatePasswordResetByID(c context.Context, tx dal.ReadwriteTransaction, entity *models.PasswordResetData) (models.PasswordReset, error)
	SavePasswordResetByID(c context.Context, tx dal.ReadwriteTransaction, record models.PasswordReset) (err error)
}

type EmailDal interface {
	InsertEmail(c context.Context, tx dal.ReadwriteTransaction, entity *models.EmailData) (models.Email, error)
	UpdateEmail(c context.Context, tx dal.ReadwriteTransaction, email models.Email) error
	GetEmailByID(c context.Context, tx dal.ReadSession, id int64) (models.Email, error)
}

type FeedbackDal interface {
	GetFeedbackByID(c context.Context, tx dal.ReadSession, feedbackID int64) (feedback models.Feedback, err error)
}

type ContactDal interface {
	GetLatestContacts(whc botsfw.WebhookContext, tx dal.ReadSession, limit, totalCount int) (contacts []models.ContactEntry, err error)
	InsertContact(c context.Context, tx dal.ReadwriteTransaction, contactEntity *models.DebtusContactDbo) (contact models.ContactEntry, err error)
	//CreateContact(c context.Context, userID int64, contactDetails models.ContactDetails) (contact models.ContactEntry, user models.AppUser, err error)
	//CreateContactWithinTransaction(c context.Context, user models.AppUser, contactUserID, counterpartyCounterpartyID int64, contactDetails models.ContactDetails, balanced money.Balanced) (contact models.ContactEntry, err error)
	//UpdateContact(c context.Context, contactID int64, values map[string]string) (contactEntity *models.DebtusContactDbo, err error)
	GetContactIDsByTitle(c context.Context, tx dal.ReadSession, userID string, title string, caseSensitive bool) (contactIDs []string, err error)
	GetContactsWithDebts(c context.Context, tx dal.ReadSession, userID string) (contacts []models.ContactEntry, err error)
}

type BillsHolderGetter func(c context.Context) (billsHolder dal.Record, err error)

type BillDal interface {
	SaveBill(c context.Context, tx dal.ReadwriteTransaction, bill models.Bill) (err error)
	UpdateBillsHolder(c context.Context, tx dal.ReadwriteTransaction, billID string, getBillsHolder BillsHolderGetter) (err error)
}

type SplitDal interface {
	GetSplitByID(c context.Context, splitID int) (split models.Split, err error)
	InsertSplit(c context.Context, splitEntity models.SplitEntity) (split models.Split, err error)
}

type TgGroupDal interface {
	GetTgGroupByID(c context.Context, tx dal.ReadSession, id int64) (tgGroup models.TgGroup, err error)
	SaveTgGroup(c context.Context, tx dal.ReadwriteTransaction, tgGroup models.TgGroup) (err error)
}

type BillScheduleDal interface {
	GetBillScheduleByID(c context.Context, id int64) (billSchedule models.BillSchedule, err error)
	InsertBillSchedule(c context.Context, billScheduleEntity *models.BillScheduleEntity) (billSchedule models.BillSchedule, err error)
	UpdateBillSchedule(c context.Context, billSchedule models.BillSchedule) (err error)
}

type GroupDal interface {
	GetGroupByID(c context.Context, tx dal.ReadSession, groupID string) (group models.GroupEntry, err error)
	InsertGroup(c context.Context, tx dal.ReadwriteTransaction, groupEntity *models.GroupDbo) (group models.GroupEntry, err error)
	SaveGroup(c context.Context, tx dal.ReadwriteTransaction, group models.GroupEntry) (err error)
	DelayUpdateGroupWithBill(c context.Context, groupID, billID string) error
}

//type GroupMemberDal interface {
//	GetGroupMemberByID(c context.Context, groupMemberID int64) (groupMember models.GroupMember, err error)
//	CreateGroupMember(c context.Context, groupMemberEntity *models.GroupMemberData) (groupMember models.GroupMember, err error)
//}

type UserGoogleDal interface {
	GetUserGoogleByID(c context.Context, googleUserID string) (userGoogle models.UserAccount, err error)
	DeleteUserGoogle(c context.Context, googleUserID string) (err error)
	SaveUserGoogle(c context.Context, userGoogle models.UserAccount) (err error)
}

type UserVkDal interface {
	GetUserVkByID(c context.Context, vkUserID int64) (userGoogle models.UserVk, err error)
	SaveUserVk(c context.Context, userVk models.UserVk) (err error)
}

type UserEmailDal interface {
	GetUserEmailByID(c context.Context, tx dal.ReadSession, email string) (userEmail models.UserEmailEntry, err error)
	SaveUserEmail(c context.Context, tx dal.ReadwriteTransaction, userEmail models.UserEmailEntry) (err error)
}

type UserGooglePlusDal interface {
	GetUserGooglePlusByID(c context.Context, id string) (userGooglePlus models.UserGooglePlus, err error)
	SaveUserGooglePlusByID(c context.Context, userGooglePlus models.UserGooglePlus) (err error)
}

type UserFacebookDal interface {
	GetFbUserByFbID(c context.Context, fbAppOrPageID, fbUserOrPageScopeID string) (fbUser models.UserFacebook, err error)
	SaveFbUser(c context.Context, tx dal.ReadwriteTransaction, fbUser models.UserFacebook) (err error)
	DeleteFbUser(c context.Context, fbAppOrPageID, fbUserOrPageScopeID string) (err error)
	//CreateFbUserRecord(c context.Context, fbUserID string, appUserID int64) (fbUser models.UserFacebook, err error)
}

type LoginPinDal interface {
	GetLoginPinByID(c context.Context, tx dal.ReadSession, loginID int) (loginPin models.LoginPin, err error)
	SaveLoginPin(c context.Context, tx dal.ReadwriteTransaction, loginPin models.LoginPin) (err error)
	CreateLoginPin(c context.Context, tx dal.ReadwriteTransaction, channel, gaClientID string, createdUserID string) (loginPin models.LoginPin, err error)
}

type LoginCodeDal interface {
	NewLoginCode(c context.Context, userID string) (code int, err error)
	ClaimLoginCode(c context.Context, code int) (userID string, err error)
}

type TwilioDal interface {
	GetLastTwilioSmsesForUser(c context.Context, tx dal.ReadSession, userID string, to string, limit int) (result []models.TwilioSms, err error)
	SaveTwilioSms(
		c context.Context,
		smsResponse *gotwilio.SmsResponse,
		transfer models.TransferEntry,
		phoneContact models.PhoneContact,
		userID string,
		tgChatID int64,
		smsStatusMessageID int,
	) (twiliosSms models.TwilioSms, err error)
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
	GetInvite(c context.Context, tx dal.ReadSession, inviteCode string) (models.Invite, error)
	ClaimInvite(c context.Context, userID string, inviteCode, claimedOn, claimedVia string) (err error)
	ClaimInvite2(c context.Context, inviteCode string, invite models.Invite, claimedByUserID string, claimedOn, claimedVia string) (err error)
	CreatePersonalInvite(ec strongoapp.ExecutionContext, userID string, inviteBy models.InviteBy, inviteToAddress, createdOnPlatform, createdOnID, related string) (models.Invite, error)
	CreateMassInvite(ec strongoapp.ExecutionContext, userID string, inviteCode string, maxClaimsCount int32, createdOnPlatform string) (invite models.Invite, err error)
}

type AdminDal interface {
	DeleteAll(c context.Context, botCode, botChatID string) error
	LatestUsers(c context.Context) (users []models.AppUser, err error)
}

type UserBrowserDal interface {
	SaveUserBrowser(c context.Context, userID string, userAgent string) (userBrowser models.UserBrowser, err error)
}

type UserOneSignalDal interface {
	SaveUserOneSignal(c context.Context, userID string, oneSignalUserID string) (userOneSignal models.UserOneSignal, err error)
}

type UserGaClientDal interface {
	SaveGaClient(c context.Context, gaClientId, userAgent, ipAddress string) (gaClient models.GaClient, err error)
}

type TgChatDal interface {
	GetTgChatByID(c context.Context, tgBotID string, tgChatID int64) (tgChat models.DebtusTelegramChat, err error)
	DoSomething( // TODO: WTF name?
		c context.Context,
		userTask *sync.WaitGroup,
		currency string,
		tgChatID int64,
		authInfo auth.AuthInfo,
		user models.AppUser,
		sendToTelegram func(tgChat botsfwtgmodels.TgChatData) error,
	) (err error)
}

type TgUserDal interface {
	FindByUserName(c context.Context, tx dal.ReadSession, userName string) (tgUsers []botsfwtgmodels.TgBotUser, err error)
}

//type TaskQueueDal interface {
//	CallDelayFunc(c context.Context, queueName, subPath, key string, f interface{}, args ...interface{}) error
//}

var (
	DB             dal.DB
	Contact        ContactDal
	User           UserDal
	UserFacebook   UserFacebookDal
	UserGoogle     UserGoogleDal
	UserGooglePlus UserGooglePlusDal

	PasswordReset PasswordResetDal
	Email         EmailDal
	UserEmail     UserEmailDal
	UserBrowser   UserBrowserDal
	UserOneSignal UserOneSignalDal
	UserGaClient  UserGaClientDal
	Feedback      FeedbackDal
	Bill          BillDal
	Receipt       ReceiptDal
	Group         GroupDal
	Reminder      ReminderDal
	TgGroup       TgGroupDal
	Transfer      TransferDal
	LoginPin      LoginPinDal
	LoginCode     LoginCodeDal
	Twilio        TwilioDal
	Invite        InviteDal
	Admin         AdminDal
	TgChat        TgChatDal
	TgUser        TgUserDal
	HttpClient    func(c context.Context) *http.Client
	BotHost       botsfw.BotHost
	HttpAppHost   strongoapp.HttpAppHost

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
