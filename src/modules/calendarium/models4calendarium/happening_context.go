package models4calendarium

import "github.com/dal-go/dalgo/record"

type HappeningContext struct {
	record.WithID[string]
	Dto *HappeningDto
}
