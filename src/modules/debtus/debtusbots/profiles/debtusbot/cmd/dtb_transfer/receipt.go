package dtb_transfer

import (
	"bytes"
	"fmt"
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw/botinput"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/sneat-co/debtstracker-translations/trans"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/facade4debtus"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/gae_app/debtstracker/analytics"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/gae_app/debtstracker/dtdal"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/gae_app/general"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/models4debtus"
	"github.com/strongo/i18n"
	"github.com/strongo/logus"
	"html/template"
	"net/url"
	"strings"
	"time"

	"github.com/bots-go-framework/bots-fw-telegram"
)

//func InlineAcceptTransfer(whc botsfw.WebhookContext) (m botsfw.MessageFromBot, err error) {
//	inlineQuery := whc.InputInlineQuery()
//	m.TelegramInlineCongig = &tgbotapi.InlineConfig{
//		InlineQueryID: inlineQuery.GetInlineQueryID(),
//		SwitchPMText: "Accept transfer",
//		SwitchPMParameter: "accept?transfer=ABC",
//	}
//	return m, err
//}

const CREATE_RECEIPT_IF_NO_INLINE_CHOSEN_NOTIFICATION = "create-receipt"

var CreateReceiptIfNoInlineNotificationCommand = botsfw.Command{
	Code:       CREATE_RECEIPT_IF_NO_INLINE_CHOSEN_NOTIFICATION,
	InputTypes: []botinput.WebhookInputType{botinput.WebhookInputCallbackQuery},
	CallbackAction: func(whc botsfw.WebhookContext, callbackUrl *url.URL) (m botsfw.MessageFromBot, err error) {
		return OnInlineChosenCreateReceipt(whc, whc.Input().(botinput.WebhookCallbackQuery).GetInlineMessageID(), callbackUrl)
	},
}

func InlineSendReceipt(whc botsfw.WebhookContext) (m botsfw.MessageFromBot, err error) {
	ctx := whc.Context()
	logus.Debugf(ctx, "InlineSendReceipt()")
	inlineQuery := whc.Input().(botinput.WebhookInlineQuery)
	query := inlineQuery.GetQuery()
	values, err := url.ParseQuery(query[len("receipt?"):])
	if err != nil {
		return m, err
	}
	idParam := values.Get("id")
	if cleanID := strings.Trim(idParam, " \",.;!@#$%^&*(){}[]`~?/\\|"); cleanID != idParam {
		logus.Debugf(ctx, "Unclean receipt ContactID: %v, cleaned: %v", idParam, cleanID)
		idParam = cleanID
	}
	transferID := idParam
	if transferID == "" {
		return m, fmt.Errorf("missing transfer ContactID")
	}
	var transfer models4debtus.TransferEntry
	transfer, err = facade4debtus.Transfers.GetTransferByID(ctx, nil, transferID)
	if err != nil {
		logus.Infof(ctx, "Faield to get transfer by ContactID: %v", transferID)
		return m, err
	}

	logus.Debugf(ctx, "Loaded transfer: %v", transfer)
	creator := whc.Input().GetSender()

	m.BotMessage = telegram.InlineBotMessage(tgbotapi.InlineConfig{
		InlineQueryID: inlineQuery.GetInlineQueryID(),
		//SwitchPmText: "Accept invite",
		//SwitchPmParameter: "invite?code=ABC",
		Results: []interface{}{
			tgbotapi.InlineQueryResultArticle{
				ID:          query,
				Type:        "article",                                                          // TODO: Move to constructor
				ThumbURL:    "https://debtstracker-io.appspot.com/img/debtstracker-512x512.png", //TODO: Replace with receipt image
				ThumbHeight: 512,
				ThumbWidth:  512,
				Title:       fmt.Sprintf(whc.Translate(trans.INLINE_RECEIPT_TITLE), transfer.Data.GetAmount()),
				Description: whc.Translate(trans.INLINE_RECEIPT_DESCRIPTION),
				InputMessageContent: tgbotapi.InputTextMessageContent{
					Text:      getInlineReceiptMessageText(whc, whc.GetBotCode(), whc.Locale().Code5, fmt.Sprintf("%v %v", creator.GetFirstName(), creator.GetLastName()), ""),
					ParseMode: "HTML",
				},
				ReplyMarkup: &tgbotapi.InlineKeyboardMarkup{
					InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{
						{
							{
								Text:         whc.Translate(trans.COMMAND_TEXT_WAIT_A_SECOND),
								CallbackData: fmt.Sprintf("%s?id=%s", CREATE_RECEIPT_IF_NO_INLINE_CHOSEN_NOTIFICATION, transferID),
							},
						},
					},
				},
			},
		},
	})
	logus.Debugf(ctx, "MessageFromBot: %v", m)

	//logus.Debugf(ctx, "Calling botApi.Send(inlineConfig=%v)", inlineConfig)
	//
	//botApi := &tgbotapi.BotAPI{
	//	Token:  whc.GetBotSettings().Token,
	//	Debug:  true,
	//	Client: whc.GetHTTPClient(),
	//}
	//mes, err := botApi.AnswerInlineQuery(inlineConfig)
	//if err != nil {
	//	logus.Errorf(ctx, "Failed to send inline results: %v", err)
	//}
	//s, err := json.Marshal(mes)
	//if err != nil {
	//	logus.Errorf(ctx, "Failed to marshal inline results response: %v, %v", err, mes)
	//}
	//logus.Infof(ctx, "botApi.Send(inlineConfig): %v", string(s))
	return m, err
}

func getInlineReceiptMessageText(t i18n.SingleLocaleTranslator, botCode, localeCode5, creator string, receiptID string) string {
	data := map[string]interface{}{
		"Creator":  creator,
		"SiteLink": template.HTML(`<a href="https://debtus.app/#utm_source=telegram&utm_medium=bot&utm_campaign=receipt-inline">DebtsTracker.IO</a>`),
	}
	if receiptID != "" {
		data["ReceiptUrl"] = GetUrlForReceiptInTelegram(botCode, receiptID, localeCode5)
	}
	var buf bytes.Buffer
	if receiptID == "" {
		buf.WriteString(t.Translate(trans.INLINE_RECEIPT_GENERATING_MESSAGE, data))
	} else {
		buf.WriteString(t.Translate(trans.INLINE_RECEIPT_MESSAGE, data))
	}

	//buf.WriteString("\n\n" + t.Translate(trans.INLINE_RECEIPT_FOOTER, data))

	if receiptID != "" {
		buf.WriteString("\n\n" + t.Translate(trans.INLINE_RECEIPT_CHOOSE_LANGUAGE, data))
	}

	return buf.String()
}

func OnInlineChosenCreateReceipt(whc botsfw.WebhookContext, inlineMessageID string, queryUrl *url.URL) (m botsfw.MessageFromBot, err error) {
	ctx := whc.Context()

	logus.Debugf(ctx, "OnInlineChosenCreateReceipt(queryUrl: %v)", queryUrl)
	transferID := queryUrl.Query().Get("id")
	creator := whc.Input().GetSender()
	creatorName := fmt.Sprintf("%v %v", creator.GetFirstName(), creator.GetLastName())

	transfer, err := facade4debtus.Transfers.GetTransferByID(ctx, nil, transferID)
	if err != nil {
		return m, err
	}
	receiptData := models4debtus.NewReceiptEntity(whc.AppUserID(), transferID, transfer.Data.Counterparty().UserID, whc.Locale().Code5, telegram.PlatformID, "", general.CreatedOn{
		CreatedOnID:       whc.GetBotCode(), // TODO: Replace with method call.
		CreatedOnPlatform: whc.BotPlatform().ID(),
	})
	receipt, err := dtdal.Receipt.CreateReceipt(ctx, receiptData)
	if err != nil {
		return m, err
	}

	if err = dtdal.Receipt.DelayedMarkReceiptAsSent(ctx, receipt.ID, transferID, time.Now()); err != nil {
		logus.Errorf(ctx, "Failed DelayedMarkReceiptAsSent: %v", err)
	}
	if m, err = showReceiptAnnouncement(whc, receipt.ID, creatorName); err != nil {
		return m, err
	}

	if err = analytics.ReceiptSentFromBot(whc, "telegram"); err != nil {
		logus.Errorf(ctx, "Failed to send analytics.ReceiptSentFromBot: %v", err)
		err = nil
	}

	//_, err = whc.Responder().SendMessage(ctx, m, botsfw.BotAPISendMessageOverHTTPS)
	//if err != nil {
	//	logus.Errorf(ctx, "Failed to send inline response: %v", err.Error())
	//}
	//m = whc.NewMessage("")
	return
}
