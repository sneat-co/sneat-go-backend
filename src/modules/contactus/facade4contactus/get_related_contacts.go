package facade4contactus

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/dal4contactus"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
)

const RelatedAsChild = "child"
const RelatedAsParent = "parent"

func GetRelatedContacts(
	ctx context.Context,
	tx dal.ReadTransaction,
	teamID, relatedAs string,
	deepness, maxDeepness int,
	contacts []dal4contactus.ContactEntry,
) (
	related []dal4contactus.ContactEntry,
	err error,
) {
	switch relatedAs {
	case RelatedAsChild, RelatedAsParent: // OK
	default:
		return nil, fmt.Errorf("unknown relatedAs: [%s]", relatedAs)
	}
	var directlyRelated []dal4contactus.ContactEntry
	for _, contact := range contacts {
		for relatedContactID, relatedContact := range contact.Data.Contacts {
			if relatedContact.RelatedAs == relatedAs {
				itemID := dbmodels.TeamItemID(relatedContactID).ItemID()
				if _, found := dal4contactus.FindContactEntryByContactID(related, itemID); !found {
					if _, found = dal4contactus.FindContactEntryByContactID(directlyRelated, itemID); !found {
						directlyRelated = append(related, dal4contactus.NewContactEntry(teamID, itemID))
					}
				}
			}
		}
	}

	if len(directlyRelated) > 0 {
		records := make([]dal.Record, len(directlyRelated))
		for i, c := range directlyRelated {
			records[i] = c.Record
		}
		if err := tx.GetMulti(ctx, records); err != nil {
			return nil, fmt.Errorf("failed to get contacts related as %v: %w", relatedAs, err)
		}
		if maxDeepness < 0 || deepness < maxDeepness {
			indirectlyRelated, err := GetRelatedContacts(ctx, tx, teamID, relatedAs, deepness+1, maxDeepness, directlyRelated)
			if err != nil {
				return nil, fmt.Errorf("failed to get indirectly related contacts: %w", err)
			}
			related = append(related, indirectlyRelated...)
		}
		related = append(related, directlyRelated...)
	}
	return
}
