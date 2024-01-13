package dal4listus

import (
	"github.com/dal-go/dalgo/record"
)

// ListusChat is not used by bots framework
type ListusChat = record.DataWithID[string, ListusChatData]
