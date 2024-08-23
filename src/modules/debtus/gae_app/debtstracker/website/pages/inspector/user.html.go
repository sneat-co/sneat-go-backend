// Code generated by hero.
// source: /Users/astec/go_workspace/src/github.com/sneat-co/sneat-go-backend/debtusbot/gae_app/debtusbot/website/pages/inspector/user.html
// DO NOT EDIT!
package inspector

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/models4debtus"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/dbo4userus"
	"io"
	"time"

	"github.com/shiyanhui/hero"
)

func renderUserPage(
	now time.Time,
	user dbo4userus.UserEntry,
	debtusSpace models4debtus.DebtusSpaceEntry,
	userBalances balances,
	contactsMissingInJson, contactsMissedByQuery, matchedContacts map[string]contactWithBalances,
	contactInfosNotFoundInDb map[string]*models4debtus.DebtusContactBrief,
	w io.Writer) {
	_buffer := hero.GetBuffer()
	defer hero.PutBuffer(_buffer)
	_buffer.WriteString(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>`)
	hero.EscapeHTML(user.ID, _buffer)

	_buffer.WriteString(`</title>
    <link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/4.0.0-beta.2/css/bootstrap.min.css" integrity="sha384-PsH8R72JQ3SOdhVi3uxftmaW6Vc51MKb0q5P2rRUpPvrszuE4W1povHYgTpBfshb" crossorigin="anonymous">
    <style>
        .d {text-align: right;}
        td.d {font-family: Courier New;}
        .center {text-align: center;}
    </style>
</head>
<body>
<div class="container-fluid p-4">
`)
	_buffer.WriteString(`

<div class="row">
    <div class="col">
        <h1>UserEntry # `)
	hero.EscapeHTML(user.ID, _buffer)
	_buffer.WriteString(`</h1>

        <table class="table">
            <thead>
            <tr class=thead-light>
                <th>First name</th>
                <th>Last name</th>
                <th>Nickname</th>
                <th>Screen name</th>
                <th>Full name</th>
            </tr>
            </thead>
            <tbody>
            <tr>
                <td>`)
	hero.EscapeHTML(user.Data.Names.FirstName, _buffer)
	_buffer.WriteString(`</td>
                <td>`)
	hero.EscapeHTML(user.Data.Names.LastName, _buffer)
	_buffer.WriteString(`</td>
                <td>`)
	hero.EscapeHTML(user.Data.Names.NickName, _buffer)
	_buffer.WriteString(`</td>
                <td>`)
	hero.EscapeHTML(user.Data.Names.ScreenName, _buffer)
	_buffer.WriteString(`</td>
                <td>`)
	hero.EscapeHTML(user.Data.Names.GetFullName(), _buffer)
	_buffer.WriteString(`</td>
            </tr>
            </tbody>
        </table>
    </div>
</div>

<div class="row">
    `)

	renderUserBalance("Balance (no interest)", userBalances.withoutInterest, false, _buffer)
	renderUserBalance("Balance with interest", userBalances.withInterest, false, _buffer)

	_buffer.WriteString(`
</div>

`)
	if len(contactInfosNotFoundInDb) > 0 {
		_buffer.WriteString(`
<h3 class=row>Contacts not found by ContactID: `)
		hero.FormatInt(int64(len(contactInfosNotFoundInDb)), _buffer)
		_buffer.WriteString(`</h3>
<table>
    <thead>
    <tr>
        <th>Name</th>
        <th>Telegram</th>
    </thead>
    <tbody>")
    for _, contactInfo := range contactInfosNotFoundInDb {
    fmt.Fprintf(w, "
    <tr>
        <td>%v</td>
        <td>%v</td>
        <td>%v</td>
    </tr>
    ", contactInfo.ContactID, contactInfo.Name, contactInfo.TgUserID)
    }
    fmt.Fprintln(w, "
    </tbody>
</table>
`)
	}

	if len(contactsMissingInJson) > 0 {
		heroContactsBlock(now, "Contacts missed in cache", contactsMissingInJson, _buffer)
	}

	if len(contactsMissedByQuery) > 0 {
		heroContactsBlock(now, "Contacts missed by query", contactsMissedByQuery, _buffer)
	}

	if len(matchedContacts) > 0 {
		heroContactsBlock(now, "Matched contacts", matchedContacts, _buffer)
	}

	_buffer.WriteString(`
</div>
</body>
</html>`)
	w.Write(_buffer.Bytes())

}
