package facade4logist

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/coremodules/contactus/dal4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/logistus/dbo4logist"
	"github.com/sneat-co/sneat-go-backend/src/modules/logistus/dto4logist"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/slice"
	"github.com/strongo/strongoapp/with"
	"github.com/strongo/validation"
	"strings"
)

// SetOrderCounterparties sets order Contacts
func SetOrderCounterparties(
	ctx context.Context,
	userCtx facade.UserContext,
	request dto4logist.SetOrderCounterpartiesRequest,
) (orderCounterparties []*dbo4logist.OrderCounterparty, err error) {
	//for i := range request.Contacts {
	//	request.Contacts[i].Instructions = strings.TrimSpace(request.Contacts[i].Instructions)
	//}
	err = RunOrderWorker(ctx, userCtx, request.OrderRequest, func(ctx context.Context, tx dal.ReadwriteTransaction, params *OrderWorkerParams) (err error) {
		orderCounterparties, err = setOrderCounterpartyTxWorker(ctx, userCtx.GetUserID(), tx, params, request)
		return
	})
	return orderCounterparties, err
}

func setOrderCounterpartyTxWorker(
	ctx context.Context,
	userID string,
	tx dal.ReadwriteTransaction,
	params *OrderWorkerParams,
	request dto4logist.SetOrderCounterpartiesRequest,
) (orderCounterparties []*dbo4logist.OrderCounterparty, err error) {
	if slice.Index(params.SpaceWorkerParams.Space.Data.UserIDs, userID) < 0 {
		return nil, facade.ErrUnauthorized
	}
	order := params.Order
	contacts := make([]dal4contactus.ContactEntry, 0)
	recordsToGet := []dal.Record{
		order.Record, // we need to get an order record for updating
		// And we need to get all Contacts that are going to be added as Contacts
	}

counterparties:
	for _, counterparty := range request.Counterparties {
		for _, contact := range contacts {
			if contact.ID == counterparty.ContactID {
				continue counterparties
			}
		}
		contact := dal4contactus.NewContactEntry(request.SpaceID, counterparty.ContactID)
		contacts = append(contacts, contact)
		recordsToGet = append(recordsToGet, contact.Record)
	}

	if err := tx.GetMulti(ctx, recordsToGet); err != nil {
		return nil, fmt.Errorf("failed to get order or contact records: %w", err)
	}

	getCounterpartyContact := func(contactID string) (contact dal4contactus.ContactEntry) {
		for _, c := range contacts {
			if c.ID == contactID {
				return c
			}
		}
		return contact
	}

	for _, counterparty := range request.Counterparties {
		_, orderCounterparty := order.Dto.GetCounterpartyByRoleAndContactID(counterparty.Role, counterparty.ContactID)
		contact := getCounterpartyContact(counterparty.ContactID)
		if orderCounterparty != nil {
			if counterparty.RefNumber != "" {
				orderCounterparty.RefNumber = strings.TrimSpace(counterparty.RefNumber)
			}
			//orderCounterparty.Address = contact.Data.Address
			//if counterparty.Instructions != "" {
			//	orderCounterparty.Instructions = strings.TrimSpace(counterparty.Instructions)
			//}
			continue
		}
		if contact.ID == "" {
			return nil, fmt.Errorf("`%v` contact not found by ContactID: %s", counterparty.Role, counterparty.ContactID)
		}

		_, orderContact := order.Dto.GetContactByID(counterparty.ContactID)
		isNewOrderContact := orderContact == nil
		if isNewOrderContact {
			orderContact = &dbo4logist.OrderContact{
				ID:        counterparty.ContactID,
				Type:      contact.Data.Type,
				Title:     contact.Data.Title,
				ParentID:  contact.Data.ParentID,
				CountryID: contact.Data.CountryID,
			}
			if orderContact.CountryID == "" {
				if counterparty.Role == dbo4logist.CounterpartyRoleShip {
					orderContact.CountryID = with.UnknownCountryID
				} else {
					return nil, validation.NewErrBadRecordFieldValue("contact.Data.CountryID", "only contacts with type=ship can have empty country ContactID")
				}
			}
			//if contact.Data.Address != nil {
			//	orderContact.Address = *contact.Data.Address
			//}
			order.Dto.Contacts = append(order.Dto.Contacts, orderContact)
			params.Changed.Contacts = true

		}
		newCounterparty := &dbo4logist.OrderCounterparty{
			ContactID: contact.ID,
			CountryID: contact.Data.CountryID,
			Title:     contact.Data.Title,
			Role:      counterparty.Role,
			RefNumber: strings.TrimSpace(counterparty.RefNumber),
			//Instructions: strings.TrimSpace(counterparty.Instructions),
		}
		if isNewOrderContact && contact.Data.ParentID != "" {
			if err := setOrderCounterpartyParent(ctx, tx, order, contact.Data.ParentID, newCounterparty); err != nil {
				return nil, err
			}
		}
		if contact.Data.ParentID != "" {
			parentRole := getParentRoleByChildRole(newCounterparty.Role)
			if parentRole == "" {
				return nil, fmt.Errorf("unsupported child counterparty role: %s", newCounterparty.Role)
			}
			newCounterparty.Parent = &dbo4logist.CounterpartyParent{
				ContactID: contact.Data.ParentID,
				Role:      parentRole,
			}
		}

		orderCounterparties = append(orderCounterparties, newCounterparty)
		var i int
		var oldCounterparty *dbo4logist.OrderCounterparty
		if counterparty.Role != dbo4logist.CounterpartyRoleShippingLine {
			i, oldCounterparty = order.Dto.GetCounterpartyByRole(counterparty.Role)
		}
		if oldCounterparty != nil {
			if oldCounterparty.ContactID == newCounterparty.ContactID && newCounterparty.RefNumber == "" {
				newCounterparty.RefNumber = strings.TrimSpace(oldCounterparty.RefNumber)
			} else {
				newCounterparty.RefNumber = strings.TrimSpace(newCounterparty.RefNumber)
			}
			//if oldCounterparty.ContactID == newCounterparty.ContactID && newCounterparty.Instructions == "" {
			//	newCounterparty.RefNumber = strings.TrimSpace(oldCounterparty.Instructions)
			//} else {
			//	newCounterparty.Instructions = strings.TrimSpace(newCounterparty.Instructions)
			//}
			if *oldCounterparty == *newCounterparty {
				continue
			}
		}
		if i < 0 {
			order.Dto.Counterparties = append(order.Dto.Counterparties, newCounterparty)
		} else {
			// Replace oldCounterparty by newCounterparty
			order.Dto.Counterparties[i] = newCounterparty

			// Removes oldContact by index from order.Data.Contacts if it is not used by other counterparties
			if newCounterparty.ContactID != oldCounterparty.ContactID && len(order.Dto.WithCounterparties.GetCounterpartiesByContactID(oldCounterparty.ContactID)) == 0 {
				if j, _ := order.Dto.WithOrderContacts.GetContactByID(oldCounterparty.ContactID); j >= 0 {
					order.Dto.Contacts = append(order.Dto.Contacts[:j], order.Dto.Contacts[j+1:]...)
					params.Changed.Contacts = true
				}
			}
		}
	}

	order.Dto.UpdateKeys()
	if err := order.Dto.Validate(); err != nil {
		return nil, fmt.Errorf("order record is not valid: %w", err)
	}

	params.Changed.Counterparties = true
	return orderCounterparties, nil
}

func getParentRoleByChildRole(childRole dbo4logist.CounterpartyRole) dbo4logist.CounterpartyRole {
	switch childRole {
	case dbo4logist.CounterpartyRoleShip:
		return dbo4logist.CounterpartyRoleShippingLine
	default:
		return ""
	}
}
func setOrderCounterpartyParent(
	ctx context.Context,
	tx dal.ReadwriteTransaction,
	order dbo4logist.Order,
	parentID string,
	childCounterparty *dbo4logist.OrderCounterparty,
) error {
	parentRole := getParentRoleByChildRole(childCounterparty.Role)
	if parentRole == "" {
		return fmt.Errorf("unsupported child counterparty role: %s", childCounterparty.Role)
	}
	_, parentCounterparty := order.Dto.GetCounterpartyByRoleAndContactID(parentRole, parentID)
	if parentCounterparty != nil {
		return nil
	}
	_, parentOrderContact := order.Dto.GetContactByID(parentID)
	if parentOrderContact == nil {
		parentContact := dal4contactus.NewContactEntry(order.Dto.SpaceID, parentID)
		if err := tx.Get(ctx, parentContact.Record); err != nil {
			return fmt.Errorf("failed to get parent contact record: %w", err)
		}
		parentOrderContact = &dbo4logist.OrderContact{
			ID:        parentContact.ID,
			Type:      parentContact.Data.Type,
			Title:     parentContact.Data.Title,
			CountryID: parentContact.Data.CountryID,
		}
		order.Dto.Contacts = append(order.Dto.Contacts, parentOrderContact)
	}
	order.Dto.Counterparties = append(order.Dto.Counterparties, &dbo4logist.OrderCounterparty{
		ContactID: parentOrderContact.ID,
		Role:      parentRole,
		Title:     parentOrderContact.Title,
		CountryID: parentOrderContact.CountryID,
	})
	return nil
}
