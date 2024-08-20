// Code generated by hero.
// source: /Users/astec/go_workspace/src/github.com/sneat-co/sneat-go-backend/debtusbot/gae_app/debtusbot/website/pages/inspector/contact.html
// DO NOT EDIT!
package inspector

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/models4debtus"
	"io"

	"github.com/shiyanhui/hero"
)

func RenderContactPage(contact models4debtus.DebtusSpaceContactEntry, w io.Writer) {
	_buffer := hero.GetBuffer()
	defer hero.PutBuffer(_buffer)
	_buffer.WriteString(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>`)
	_buffer.WriteString(`DebtusSpaceContactEntry # `)
	hero.EscapeHTML(contact.ID, _buffer)

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
            <h1>DebtusSpaceContactEntry # `)
	hero.EscapeHTML(contact.ID, _buffer)
	_buffer.WriteString(`</h1>
        </div>

`)

	_buffer.WriteString(`
</div>
</body>
</html>`)
	w.Write(_buffer.Bytes())

}
