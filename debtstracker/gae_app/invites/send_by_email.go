package invites

import (
	"bytes"
	"context"
	"fmt"
	"github.com/sneat-co/debtstracker-translations/trans"
	"github.com/sneat-co/sneat-go-core/emails"
	"github.com/strongo/i18n"
	"github.com/strongo/strongoapp"
	"html/template"

	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/common"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/emailing"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/models"
)

type InviteTemplateParams struct {
	ToName     string
	FromName   string
	InviteCode string
	InviteURL  string
	ReceiptURL template.HTML
	TgBot      string
	Utm        string
}

func SendInviteByEmail(ec strongoapp.ExecutionContext, translator i18n.SingleLocaleTranslator, fromName, toEmail, toName, inviteCode, telegramBotID, utmSource string) (emailID string, err error) {
	//cred := credentials.NewStaticCredentials(, , "")
	//credStaticProvider := credentials.StaticProvider{}
	//credStaticProvider.AccessKeyID = "******"
	//credStaticProvider.SecretAccessKey = "*******"
	//credStaticProvider.ProviderName = "Static"
	//htmlTemplate, err := template.New("html").Parse(Translate(EMAIL_INVITE_HTML, whc))
	//if err != nil {
	//	return err
	//}
	//var html bytes.Buffer
	//htmlTemplate.Execute(&html)

	templateParams := InviteTemplateParams{
		ToName:     toName,
		FromName:   fromName,
		InviteCode: inviteCode,
		TgBot:      telegramBotID,
		Utm: common.UtmParams{
			Source:   utmSource,
			Medium:   string(models.InviteByEmail),
			Campaign: common.UTM_CAMPAIGN_ONBOARDING_INVITE,
		}.String(),
	}

	c := ec.Context()

	subject, err := emailing.GetEmailText(c, translator, trans.EMAIL_INVITE_SUBJ, templateParams)
	if err != nil {
		return "", err
	}

	bodyText, err := emailing.GetEmailText(c, translator, trans.EMAIL_INVITE_TEXT, templateParams)
	if err != nil {
		return "", err
	}

	bodyHtml, err := emailing.GetEmailHtml(c, translator, trans.EMAIL_INVITE_HTML, templateParams)
	if err != nil {
		return "", err
	}

	emailMessage := emails.Email{
		From:    "invite@sneat.app",
		To:      []string{toEmail},
		Subject: subject,
		Text:    bodyText,
		HTML:    bodyHtml,
	}
	emailID, err = emailing.SendEmail(c, emailMessage)
	return
}

func SendReceiptByEmail(c context.Context, translator i18n.SingleLocaleTranslator, receipt models.Receipt, fromName, toName, toEmail string) (emailID string, err error) {
	templateParams := struct {
		ToName     string
		FromName   string
		ReceiptID  string
		ReceiptURL template.HTML
	}{
		toName,
		fromName,
		receipt.ID,
		template.HTML(""),
	}

	subject, err := common.TextTemplates.RenderTemplate(c, translator, trans.EMAIL_RECEIPT_SUBJ, templateParams)
	if err != nil {
		return "", err
	}

	bodyText, err := common.TextTemplates.RenderTemplate(c, translator, trans.EMAIL_RECEIPT_BODY_TEXT, templateParams)
	if err != nil {
		return "", err
	}

	receiptURL := common.GetReceiptUrl(receipt.ID, common.GetWebsiteHost(receipt.Data.CreatedOnID))
	//displayUrl := strings.Split(string(templateParams.ReceiptURL), "#")[0]
	templateParams.ReceiptURL = template.HTML(fmt.Sprintf(`<a href="%v">%v</a>`, receiptURL, receiptURL))
	var bodyHtml bytes.Buffer
	if err = common.HtmlTemplates.RenderTemplate(c, &bodyHtml, translator, trans.EMAIL_RECEIPT_BODY_HTML, templateParams); err != nil {
		return "", err
	}

	emailMessage := emails.Email{
		From:    "receipt@debtus.app",
		To:      []string{toEmail},
		Subject: subject,
		Text:    bodyText,
		HTML:    bodyHtml.String(),
	}
	return emailing.SendEmail(c, emailMessage)
}
