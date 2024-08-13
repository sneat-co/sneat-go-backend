package admin

import (
	"context"
	"fmt"
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/gae_app/debtstracker/dtdal"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/models4debtus"
)

func SendFeedbackToAdmins(c context.Context, botToken string, feedback models4debtus.Feedback) (err error) {
	bot := tgbotapi.NewBotAPIWithClient(botToken, dtdal.HttpClient(c))
	text := fmt.Sprintf("%v user #%s @%v (rate=%v):\n%v", feedback.CreatedOnPlatform, feedback.UserStrID, feedback.CreatedOnID, feedback.Rate, feedback.Text)
	message := tgbotapi.NewMessageToChannel("-1001128307094", text)
	message.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		[]tgbotapi.InlineKeyboardButton{
			{Text: "Reply to feedback", URL: fmt.Sprintf("https://debtstracker.io/app/#/reply-to-feedback/%d", feedback.ID)},
		},
	)
	_, err = bot.Send(message)
	return
}
