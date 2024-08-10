package unsorted

import (
	"errors"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/debtstracker-translations/trans"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/api"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/facade2debtus/dto"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/i18n"
	"github.com/strongo/logus"
	"github.com/strongo/strongoapp"
	"net/http"
	"strings"
	"time"

	"context"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/bot"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/bot/platforms/tgbots"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/analytics"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/auth"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/common"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/dtdal"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/facade2debtus"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/models"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/general"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/invites"
)

func NewReceiptTransferDto(c context.Context, transfer models.TransferEntry) dto.ApiReceiptTransferDto {
	creator := transfer.Data.Creator()
	transferDto := dto.ApiReceiptTransferDto{
		ID:             transfer.ID,
		From:           dto.NewContactDto(*transfer.Data.From()),
		To:             dto.NewContactDto(*transfer.Data.To()),
		Amount:         transfer.Data.GetAmount(),
		IsOutstanding:  transfer.Data.IsOutstanding,
		DtCreated:      transfer.Data.DtCreated,
		CreatorComment: creator.Comment,
		Creator: dto.ApiUserDto{ // TODO: Rename field - it can be not a creator in case of bill created by 3d party (paid by not by bill creator)
			ID:   creator.UserID,
			Name: creator.ContactName,
		},
	}
	// Set acknowledge info if presented
	if !transfer.Data.AcknowledgeTime.IsZero() {
		transferDto.Acknowledge = &dto.ApiAcknowledgeDto{
			Status:   transfer.Data.AcknowledgeStatus,
			UnixTime: transfer.Data.AcknowledgeTime.Unix(),
		}
	}
	if transferDto.From.Name == "" {
		logus.Warningf(c, "transferDto.From.Name is empty string")
	}

	if transferDto.To.Name == "" {
		logus.Warningf(c, "transferDto.To.Name is empty string")
	}
	return transferDto
}

func HandleGetReceipt(c context.Context, w http.ResponseWriter, r *http.Request) {
	receiptID := api.GetStrID(c, w, r, "id")
	if receiptID == "" {
		return
	}

	receipt, err := dtdal.Receipt.GetReceiptByID(c, nil, receiptID)
	if api.HasError(c, w, err, models.ReceiptKind, receiptID, http.StatusBadRequest) {
		return
	}

	var transfer models.TransferEntry
	if err = facade.RunReadwriteTransaction(c, func(c context.Context, tx dal.ReadwriteTransaction) (err error) {
		transfer, err = facade2debtus.Transfers.GetTransferByID(c, tx, receipt.Data.TransferID)
		if api.HasError(c, w, err, models.TransfersCollection, receipt.Data.TransferID, http.StatusInternalServerError) {
			return
		}

		if err = facade2debtus.CheckTransferCreatorNameAndFixIfNeeded(c, tx, transfer); api.HasError(c, w, err, models.TransfersCollection, receipt.Data.TransferID, http.StatusInternalServerError) {
			return
		}
		return nil
	}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	sentTo := receipt.Data.SentTo
	if receipt.Data.SentVia == "telegram" {
		lang := r.URL.Query().Get("lang")
		if lang == "" {
			lang = receipt.Data.Lang
		}
		env := dtdal.HttpAppHost.GetEnvironment(c, r)
		if env == strongoapp.UnknownEnv {
			w.WriteHeader(http.StatusBadRequest)
			logus.Warningf(c, "Unknown host")
		}
		botSettings, err := tgbots.GetBotSettingsByLang(env, bot.ProfileDebtus, lang)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logus.Errorf(c, fmt.Errorf("failed to get bot settings by lang: %w", err).Error())
		}
		sentTo = botSettings.Code
	}

	creator := transfer.Data.Creator()

	logus.Debugf(c, "transfer.Creator(): %v", creator)

	receiptDto := dto.ApiReceiptDto{
		ID:       receiptID,
		Code:     "ToDoCODE",
		SentVia:  receipt.Data.SentVia,
		SentTo:   sentTo,
		Transfer: NewReceiptTransferDto(c, transfer),
	}

	api.JsonToResponse(c, w, &receiptDto)
}

//func transferContactToDto(transferContact models.TransferCounterpartyInfo) dto.ContactDto {
//	return dto.NewContactDto(transferContact)
//}

func HandleReceiptAccept(c context.Context, w http.ResponseWriter, _ *http.Request) {
	api.JsonToResponse(c, w, "ok")
}

func HandleReceiptDecline(c context.Context, w http.ResponseWriter, _ *http.Request) {
	api.JsonToResponse(c, w, "ok")
}

const RECEIPT_CHANNEL_DRAFT = "draft"

func HandleSendReceipt(c context.Context, w http.ResponseWriter, r *http.Request, authInfo auth.AuthInfo, user models.AppUser) {
	w.Header().Add("Access-Control-Allow-Origin", "*")

	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("Invalid form data"))
		return
	}

	receiptID := r.FormValue("receipt")
	if receiptID == "" {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("Missing receipt parameter"))
		return
	}

	channel := r.FormValue("by")

	switch channel {
	case "email":
	case "sms":
		api.BadRequestMessage(c, w, "Not implemented yet")
		return
	default:
		api.BadRequestMessage(c, w, "Unsupported channel")
		return
	}

	toAddress := strings.TrimSpace(r.FormValue("to"))

	if toAddress == "" {
		api.BadRequestMessage(c, w, "'To' parameter is not provided")
		return
	}

	if len(toAddress) > 1024 {
		api.BadRequestMessage(c, w, "'To' parameter is too large")
		return
	}

	receipt, err := dtdal.Receipt.GetReceiptByID(c, nil, receiptID)

	if err != nil {
		var status int
		if dal.IsNotFound(err) {
			status = http.StatusBadRequest
		} else {
			status = http.StatusInternalServerError
		}
		api.ErrorAsJson(c, w, status, err)
		return
	}

	transfer, err := facade2debtus.Transfers.GetTransferByID(c, nil, receipt.Data.TransferID)
	if err != nil {
		api.ErrorAsJson(c, w, http.StatusInternalServerError, err)
		return
	}

	if transfer.Data.CreatorUserID != authInfo.UserID && transfer.Data.Counterparty().UserID != authInfo.UserID {
		api.ErrorAsJson(c, w, http.StatusUnauthorized, errors.New("this transfer does not belong to the current user"))
		return
	}

	locale := i18n.GetLocaleByCode5(user.Data.GetPreferredLocale()) // TODO: Get language from request
	translator := i18n.NewSingleMapTranslator(locale, nil /*shared.TheAppContext.GetTranslator(c)*/)

	if _, err = invites.SendReceiptByEmail(c, translator, receipt, user.Data.FullName(), transfer.Data.Counterparty().ContactName, toAddress); err != nil {
		logus.Errorf(c, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		err = fmt.Errorf("failed to send receipt by email: %w", err)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	analytics.ReceiptSentFromApi(c, r, authInfo.UserID, locale.Code5, "api", "email")

	if _, _, err = updateReceiptAndTransferOnSent(c, receiptID, channel, toAddress, locale.Code5); err != nil {
		api.ErrorAsJson(c, w, http.StatusInternalServerError, err)
		return
	}
}

func updateReceiptAndTransferOnSent(c context.Context, receiptID string, channel, sentTo, lang string) (receipt models.Receipt, transfer models.TransferEntry, err error) {

	err = facade.RunReadwriteTransaction(c, func(c context.Context, tx dal.ReadwriteTransaction) error {
		if receipt, err = dtdal.Receipt.GetReceiptByID(c, tx, receiptID); err != nil {
			return err
		}
		if receipt.Data.SentVia == RECEIPT_CHANNEL_DRAFT {
			if transfer, err = facade2debtus.Transfers.GetTransferByID(c, tx, receipt.Data.TransferID); err != nil {
				return err
			}
			receipt.Data.DtSent = time.Now()
			receipt.Data.SentVia = channel
			receipt.Data.SentTo = sentTo
			receipt.Data.Lang = lang
			transfer.Data.ReceiptsSentCount += 1
			transferHasThisReceiptID := false
			for _, rID := range transfer.Data.ReceiptIDs {
				if rID == receiptID {
					transferHasThisReceiptID = true
					break
				}
			}
			if !transferHasThisReceiptID {
				transfer.Data.ReceiptIDs = append(transfer.Data.ReceiptIDs, receiptID)
			}
			if err = tx.SetMulti(c, []dal.Record{receipt.Record, transfer.Record}); err != nil {
				err = fmt.Errorf("failed to save receipt & transfer: %w", err)
			}
		} else if receipt.Data.SentVia == channel {
			logus.Infof(c, "Receipt already has channel '%s'", channel)
		} else {
			logus.Warningf(c, "An attempt to set receipt channel to '%s' when it's alreay '%s'", channel, receipt.Data.SentVia)
		}

		return err
	})
	return
}

func HandleSetReceiptChannel(c context.Context, w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("Invalid form data"))
		return
	}

	receiptID := r.FormValue("receipt")
	if receiptID == "" {
		err := fmt.Errorf("missing receipt parameter")
		logus.Debugf(c, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	channel, err := getReceiptChannel(r)
	if err != nil {
		if errors.Is(err, errUnknownChannel) {
			m := "Unknown channel: " + channel
			logus.Debugf(c, m)
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(m))
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
		}
	}

	logus.Debugf(c, "HandleSetReceiptChannel(receiptID=%s, channel=%s)", receiptID, channel)
	if channel == RECEIPT_CHANNEL_DRAFT {
		m := fmt.Sprintf("Status '%s' is not supported in this method", RECEIPT_CHANNEL_DRAFT)
		logus.Warningf(c, m)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(m))
	}

	if _, _, err = updateReceiptAndTransferOnSent(c, receiptID, channel, "", ""); err != nil {
		if dal.IsNotFound(err) {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte("Receipt not found by ID"))
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
		}
		return
	}
	logus.Infof(c, "Done")
	_, _ = w.Write([]byte("ok"))
}

var errUnknownChannel = errors.New("Unknown channel")

func getReceiptChannel(r *http.Request) (channel string, err error) {
	channel = r.FormValue("channel")
	switch channel {
	case RECEIPT_CHANNEL_DRAFT:
	case "fbm":
	case "vk":
	case "viber":
	case "whatsapp":
	case "line":
	case "telegram":
	default:
		err = errUnknownChannel
	}
	return
}

func HandleCreateReceipt(c context.Context, w http.ResponseWriter, r *http.Request, authInfo auth.AuthInfo) {
	if err := r.ParseForm(); err != nil {
		logus.Debugf(c, "HandleCreateReceipt() => Invalid form data: "+err.Error())
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("Invalid form data"))
		return
	}
	//logus.Debugf(c, "HandleCreateReceipt() => r.Form: %+v", r.Form)
	transferID := r.FormValue("transfer")
	if transferID == "" {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("Missing transfer parameter"))
		return
	}
	transfer, err := facade2debtus.Transfers.GetTransferByID(c, nil, transferID)
	if err != nil {
		if dal.IsNotFound(err) {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			logus.Errorf(c, err.Error())
		}
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	var user models.AppUser
	if user, err = facade2debtus.User.GetUserByID(c, nil, authInfo.UserID); err != nil {
		api.ErrorAsJson(c, w, http.StatusInternalServerError, err)
		return
	}
	channel, err := getReceiptChannel(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if err == errUnknownChannel {
			_, _ = w.Write([]byte("Unknown channel: " + channel))
		}
		return
	}
	creatorUserID := transfer.Data.CreatorUserID // TODO: Get from session?

	lang := user.Data.PreferredLanguage
	if lang == "" {
		if acceptLanguage := r.Header.Get("Accept-Language"); acceptLanguage != "" {
			for _, set1 := range strings.Split(acceptLanguage, ";") {
				for _, al := range strings.Split(set1, ",") {
					switch len(al) {
					case 5:
						if _, ok := trans.SupportedLocalesByCode5[strings.ToLower(al[:2])+"-"+strings.ToUpper(al[4:])]; ok {
							lang = al
							goto langSet
						}
					case 2:
						al = strings.ToLower(al)
						for localeCode := range trans.SupportedLocalesByCode5 {
							if strings.HasPrefix(localeCode, al) {
								lang = localeCode
								goto langSet
							}
						}
					}
				}
			}
		langSet:
		}
	}
	if lang == "" {
		lang = i18n.LocaleCodeEnUS
	}
	receiptData := models.NewReceiptEntity(creatorUserID, transferID, transfer.Data.Counterparty().UserID, lang, channel, "", general.CreatedOn{
		CreatedOnPlatform: "api", // TODO: Replace with actual, pass from client
		CreatedOnID:       r.Host,
	})

	receipt, err := dtdal.Receipt.CreateReceipt(c, receiptData)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logus.Errorf(c, err.Error())
		return
	}

	//user, err = facade2debtus.User.GetUserByID(c, transfer.CreatorUserID)
	//if err != nil {
	//	w.WriteHeader(http.StatusInternalServerError)
	//	err = errors.Wrapf(err, "Failed to get user by ID=%v", transfer.CreatorUserID)
	//	w.Write([]byte(err.Error()))
	//	logus.Warningf(c, err.Error())
	//	return
	//}
	var messageToSend string

	if channel == "telegram" {
		tgBotID := transfer.Data.Creator().TgBotID
		if tgBotID == "" {
			if strings.Contains(r.URL.Host, "dev") {
				tgBotID = "DebtsTrackerDev1Bot"
			} else {
				tgBotID = "DebtsTrackerBot"
			}
		}
		messageToSend = fmt.Sprintf("https://telegram.me/%s?start=send-receipt_%s", tgBotID, receipt.ID) // TODO:
	} else {
		locale := i18n.GetLocaleByCode5(user.Data.GetPreferredLocale())
		translator := i18n.NewSingleMapTranslator(locale, nil /*shared.TheAppContext.GetTranslator(c)*/)
		//ec := strongoapp.NewExecutionContext(c, translator)

		logus.Debugf(c, "r.Host: %s", r.Host)

		templateParams := struct {
			ReceiptURL string
		}{
			ReceiptURL: common.GetReceiptUrl(receipt.ID, r.Host),
		}
		messageToSend, err = common.TextTemplates.RenderTemplate(c, translator, trans.MESSENGER_RECEIPT_TEXT, templateParams) // TODO: Consider using just ExecutionContext instead of (c, translator)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			err = fmt.Errorf("failed to render message template: %w", err)
			logus.Errorf(c, err.Error())
			_, _ = w.Write([]byte(err.Error()))
		}
	}

	api.JsonToResponse(c, w, struct {
		ID   string
		Link string
		Text string
	}{
		receipt.ID,
		// TODO: It seems wrong to use request host!
		//shared.GetReceiptUrlForUser(receipt, receiptData.CreatorUserID, receiptData.CreatedOnPlatform, receiptData.CreatedOnID)
		fmt.Sprintf("https://%s/receipt?id=%s&t=%s", r.Host, receipt.ID, time.Now().Format("20060102-150405")),
		messageToSend,
	})
}
