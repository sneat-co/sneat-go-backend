package facade4listus

import (
	"errors"
	"github.com/sneat-co/sneat-go-backend/src/modules/listus/models4listus"
	"github.com/strongo/random"
)

func generateRandomListItemID(items []*models4listus.ListItemBrief, initialID string) (id string, err error) {
	isDuplicateID := func() bool {
		for _, item := range items {
			if item.ID == id {
				return true
			}
		}
		return false
	}
	id = initialID
	if !isDuplicateID() {
		return
	}
next:
	for i := 0; i <= 100; i++ {
		id = random.ID(3)
		if isDuplicateID() {
			continue next
		}
		return
	}
	return "", errors.New("too many attempts to generate random ContactID")
}
