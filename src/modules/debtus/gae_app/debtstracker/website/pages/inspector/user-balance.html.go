// Code generated by hero.
// source: /Users/astec/go_workspace/src/github.com/sneat-co/sneat-go-backend/debtusbot/gae_app/debtusbot/website/pages/inspector/user-balance.html
// DO NOT EDIT!
package inspector

import (
	"bytes"
	"fmt"

	"github.com/shiyanhui/hero"
)

func renderUserBalance(title string, balances balancesByCurrency, showTransfers bool, buf *bytes.Buffer) {
	buf.WriteString(`
<div class="col-sm">
    <h3>`)
	hero.EscapeHTML(title, buf)
	buf.WriteString(`</h3>
    `)
	if balances.err != nil {
		buf.WriteString(`
    <div style="color:red">Error: `)
		hero.EscapeHTML(balances.err.Error(), buf)
		buf.WriteString(`</div>
    `)
	}
	buf.WriteString(`
    <table class="table table-bordered">
        <thead>
        <tr>
            <th>Currency</th>
            <th class=d>User</th>
            <th class=d>Contacts</th>
            `)
	if showTransfers {
		buf.WriteString(`
            <th class=d>Transfers</th>
            `)
	}
	buf.WriteString(`
        </tr>
        </thead>
        <tbody>
        `)
	for currency, balance := range balances.byCurrency {
		if balance.user == balance.contacts {
			buf.WriteString(`
        <tr>`)
		} else {
			buf.WriteString(`
        <tr class="table-danger">`)
		}
		if balance.userContactBalanceErr != nil {
			buf.WriteString(`
            <td colspan="4" style="color:red">`)
			hero.EscapeHTML(balance.userContactBalanceErr.Error(), buf)
			buf.WriteString(`</td>
            `)
		} else {
			buf.WriteString(`
            <td>`)
			hero.EscapeHTML(fmt.Sprintf("%v", currency), buf)
			buf.WriteString(`</td>
            <td class=d>`)
			hero.EscapeHTML(fmt.Sprintf("%v", balance.user), buf)
			buf.WriteString(`</td>
            <td class=d>`)
			hero.EscapeHTML(fmt.Sprintf("%v", balance.contacts), buf)
			buf.WriteString(`</td>
            `)
			if showTransfers {
				buf.WriteString(`
            <td class=d>`)
				hero.EscapeHTML(fmt.Sprintf("%v", balance.transfers), buf)
				buf.WriteString(`</td>
            `)
			}
		}
		buf.WriteString(`
        </tr>
        `)
	}
	buf.WriteString(`
        </tbody>
    </table>
</div>`)

}
