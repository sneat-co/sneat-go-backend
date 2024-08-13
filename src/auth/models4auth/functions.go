package models4auth

import "strings"

func FixContactName(contactName string) (isFixed bool, s string) {
	if start := strings.Index(contactName, "("); start > 0 {
		if end := strings.Index(contactName, ")"); end > start {
			if l := len(contactName); end == l-1 {
				if (end-start-1)*2 == len(contactName)-3 {
					if s = contactName[start+1 : end]; s == contactName[:start-1] {
						isFixed = true
						return
					}
				}
			}
		}
	}
	s = contactName
	return
}
