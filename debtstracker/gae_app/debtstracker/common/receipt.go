package common

import (
	"fmt"
	"strings"
)

func GetReceiptUrl(receiptID string, host string) string {
	if receiptID == "" {
		panic("receiptID is empty")
	}
	if host == "" {
		panic("host is empty string")
	} else if !strings.Contains(host, ".") {
		panic("host is not a domain name: " + host)
	}
	return fmt.Sprintf("https://%v/receipt?id=%v", host, receiptID)
}
