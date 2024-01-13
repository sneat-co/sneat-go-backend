package facade4invitus

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/sneat-co/sneat-go-backend/src/modules/invitus/models4invitus"
	"github.com/sneat-co/sneat-go-core"
	"github.com/sneat-co/sneat-go-core/capturer"
	"github.com/sneat-co/sneat-go-core/emails"
	"html/template"
)

const inviteEmailTemplateText = `
<p>
	InviteDto from: <i>{{.fromHTML}}</i>
</p>

<p>
	To join <b>{{.team.Title}}</b> please follow the link:
</p>

<p>
	https://{{.hostPath}}/join/{{.team.Type}}?id={{.id}}#pin={{.pinCode}}
</p>

<p>You personal PIN code to join the team is: <b><code>{{.pinCode}}</code></b></p>

<p>https://sneat.app/ - A family app that saves you time & money.</p>

<p>P.S. If any issues feel free to get <a href="mailto:help@sneat.app">help@sneat.app</a></p>
`

var inviteEmailTemplate = template.Must(template.New("inviteEmail").Parse(inviteEmailTemplateText))

func sendInviteEmail(ctx context.Context, id string, invite *models4invitus.PersonalInviteDto) (messageID string, err error) {
	if invite.To.Address == "" {
		return "", errors.New("missing required field: invite.To.Address")
	}
	templateData := make(map[string]interface{})
	if core.IsInProd() {
		templateData["hostPath"] = "sneat.app/pwa"
	} else {
		templateData["hostPath"] = "localhost:4200"
	}
	templateData["id"] = id
	if invite.From.Address == "" {
		templateData["fromHTML"] = invite.From.Title
	} else {
		templateData["fromHTML"] = fmt.Sprintf(`<a href="mailto:%v">%v</a>`, invite.From.Address, invite.From.Title)
	}
	templateData["invite"] = invite
	templateData["team"] = invite.Team
	templateData["pinCode"] = invite.Pin
	buf := new(bytes.Buffer)
	if err := inviteEmailTemplate.Execute(buf, templateData); err != nil {
		return "", fmt.Errorf("failed to create email message body: %w", err)
	}

	msg := emails.Email{
		From:    fmt.Sprintf(`"%v" <inviter@sneat.app>`, invite.From.Title),
		To:      []string{invite.To.Address},
		Subject: fmt.Sprintf("You are invited by %v to join %v", invite.From.Title, invite.Team.Title),
		HTML:    buf.String(),
		//ReplyTo: nil,
	}
	var response emails.Sent
	if response, err = emails.Send(msg); err != nil {
		err = capturer.CaptureError(ctx, err)
		return
	}
	if response != nil {
		messageID = response.MessageID()
	}
	return
}
