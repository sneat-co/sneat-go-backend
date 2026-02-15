package dbo4calendarium

import "github.com/dal-go/dalgo/record"

type HappeningEntry = record.DataWithID[string, *HappeningDbo]
