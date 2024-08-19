// Code generated by hero.
// source: /Users/astec/go_workspace/src/github.com/sneat-co/sneat-go-backend/debtus/gae_app/debtus/website/pages/inspector/api4transfers-page.html
// DO NOT EDIT!
package inspector

import (
	"fmt"
	"github.com/crediterra/money"
	"github.com/shiyanhui/hero"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/models4debtus"
	"io"
)

func renderTransfersPage(contact models4debtus.DebtusSpaceContactEntry, currency money.CurrencyCode, balancesWithoutInterest, balancesWithInterest balanceRow, transfers []models4debtus.TransferEntry, w io.Writer) {
	_buffer := hero.GetBuffer()
	defer hero.PutBuffer(_buffer)
	_buffer.WriteString(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>`)
	_buffer.WriteString(`DebtusSpaceContactEntry # `)
	hero.EscapeHTML(contact.ID, _buffer)
	_buffer.WriteString(`: `)
	hero.EscapeHTML(fmt.Sprintf("%v", currency), _buffer)

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
	_buffer.WriteString(`: `)
	hero.EscapeHTML(fmt.Sprintf("%v", currency), _buffer)
	_buffer.WriteString(`</h1>
</div>

<div class="row">
    <table class="table">
        <thead>
        <tr>
            <th scope="col">Balances</th>
            <th scope="col">Currency</th>
            <th scope="col" class=d>User</th>
            <th scope="col" class=d>DebtusSpaceContactEntry</th>
            <th scope="col" class=d>Transfers</th>
        </tr>
        </thead>
        <tbody>
        <tr>
            <th>Without interest</th>
            <th>`)
	hero.EscapeHTML(fmt.Sprintf("%v", currency), _buffer)
	_buffer.WriteString(`</th>
            <td class=d>`)
	hero.EscapeHTML(fmt.Sprintf("%v", balancesWithoutInterest.user), _buffer)
	_buffer.WriteString(`</td>
            <td class=d>`)
	hero.EscapeHTML(fmt.Sprintf("%v", balancesWithoutInterest.contacts), _buffer)
	_buffer.WriteString(`</td>
            <td class=d>`)
	hero.EscapeHTML(fmt.Sprintf("%v", balancesWithoutInterest.transfers), _buffer)
	_buffer.WriteString(`</td>
        </tr>
        <tr>
            <th>With interest</th>
            <th>`)
	hero.EscapeHTML(fmt.Sprintf("%v", currency), _buffer)
	_buffer.WriteString(`</th>
            <td class=d>`)
	hero.EscapeHTML(fmt.Sprintf("%v", balancesWithInterest.user), _buffer)
	_buffer.WriteString(`</td>
            <td class=d>`)
	hero.EscapeHTML(fmt.Sprintf("%v", balancesWithInterest.contacts), _buffer)
	_buffer.WriteString(`</td>
            <td class=d></td>
        </tr>
        </tbody>
    </table>
</div>


<div class="row">
    <h2>Transfers</h2>
    <table class="table">
        <thead>
        <tr>
            <th>#</th>
            <th>ContactID</th>
            <th>Created at</th>
            <th>Created on</th>
            <th>From</th>
            <th>To</th>
            <th>IsReturn</th>
            <th>IsOutstanding</th>
            <th>Interest</th>
            <th class=d>Amount</th>
            <th class=d>Returned</th>
        </tr>
        </thead>
        <tbody>
        `)
	for i, transfer := range transfers {
		_buffer.WriteString(`
        <tr>
            <td class="d">`)
		hero.FormatInt(int64(i+1), _buffer)
		_buffer.WriteString(`</td>
            <td class="d"><a href="transfer?id=`)
		hero.EscapeHTML(transfer.ID, _buffer)
		_buffer.WriteString(`">`)
		hero.EscapeHTML(transfer.ID, _buffer)
		_buffer.WriteString(`</a></td>
            <td>`)
		hero.EscapeHTML(fmt.Sprintf("%v", transfer.Data.DtCreated), _buffer)
		_buffer.WriteString(`</td>
            <td>
                `)
		if transfer.Data.CreatedOnPlatform == "telegram" {
			_buffer.WriteString(`
                <a href="https://t.me/`)
			hero.EscapeHTML(transfer.Data.CreatedOnID, _buffer)
			_buffer.WriteString(`">@`)
			hero.EscapeHTML(transfer.Data.CreatedOnID, _buffer)
			_buffer.WriteString(`</a>
                `)
		}
		if transfer.Data.CreatedOnPlatform != "telegram" {
			hero.EscapeHTML(transfer.Data.CreatedOnID, _buffer)
			_buffer.WriteString(`@`)
			hero.EscapeHTML(transfer.Data.CreatedOnPlatform, _buffer)
		}
		_buffer.WriteString(`
            </td>
            <td>`)
		hero.EscapeHTML(transfer.Data.From().Name(), _buffer)
		_buffer.WriteString(`</td>
            <td>`)
		hero.EscapeHTML(transfer.Data.To().Name(), _buffer)
		_buffer.WriteString(`</td>
            <td>`)
		hero.FormatBool(transfer.Data.IsReturn, _buffer)
		_buffer.WriteString(`</td>
            <td>`)
		hero.FormatBool(transfer.Data.IsOutstanding, _buffer)
		_buffer.WriteString(`</td>
            <td>
                `)
		if transfer.Data.InterestPercent != 0 {
			hero.EscapeHTML(fmt.Sprintf("%v", transfer.Data.InterestPercent), _buffer)
			_buffer.WriteString(`%
                per `)
			hero.FormatInt(int64(transfer.Data.InterestPeriod), _buffer)
			_buffer.WriteString(` days
                `)
			if transfer.Data.InterestMinimumPeriod > 1 {
				_buffer.WriteString(`
                minimum for `)
				hero.FormatInt(int64(transfer.Data.InterestMinimumPeriod), _buffer)
				_buffer.WriteString(` days
                `)
			}
		}
		_buffer.WriteString(`
            </td>
            <td class=d>`)
		hero.EscapeHTML(fmt.Sprintf("%v", transfer.Data.AmountInCents), _buffer)
		_buffer.WriteString(`</td>
            <td class=d>`)
		hero.EscapeHTML(fmt.Sprintf("%v", transfer.Data.AmountReturned()), _buffer)
		_buffer.WriteString(`</td>
        </tr>
        `)
	}
	_buffer.WriteString(`
        </tbody>
    </table>
</div>
`)

	_buffer.WriteString(`
</div>
</body>
</html>`)
	w.Write(_buffer.Bytes())

}