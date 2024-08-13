package facade4logist

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/dal4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/logistus/dbo4logist"
	"github.com/sneat-co/sneat-go-backend/src/modules/logistus/dto4logist"
	"github.com/sneat-co/sneat-go-backend/src/modules/spaceus/dal4spaceus"
	"github.com/sneat-co/sneat-go-core/facade"
	dbmodels2 "github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/strongo/slice"
	"github.com/strongo/strongoapp/with"
	"strconv"
	"time"
)

// CreateOrder creates a new order
func CreateOrder(
	ctx context.Context,
	userCtx facade.UserContext,
	request dto4logist.CreateOrderRequest,
) (
	orderBrief *dbo4logist.OrderBrief, err error,
) {
	err = dal4spaceus.RunSpaceWorker(ctx, userCtx, request.SpaceID,
		func(ctx context.Context, tx dal.ReadwriteTransaction, params *dal4spaceus.SpaceWorkerParams) (err error) {
			orderBrief, err = createOrderTxWorker(ctx, tx, params, params.UserID, request)
			return err
		},
	)
	return orderBrief, err
}

func createOrderTxWorker(
	ctx context.Context, tx dal.ReadwriteTransaction, params *dal4spaceus.SpaceWorkerParams,
	userID string,
	request dto4logist.CreateOrderRequest,
) (orderBrief *dbo4logist.OrderBrief, err error) {
	logistSpace := dbo4logist.NewLogistSpaceEntry(request.SpaceID)
	if err = tx.Get(ctx, logistSpace.Record); err != nil {
		if dal.IsNotFound(err) {
			err = nil // OK
		} else {
			return nil, fmt.Errorf("failed to load logistus team record: %w", err)
		}
	}

	if logistSpace.Data.OrderCounters == nil {
		logistSpace.Data.OrderCounters = make(map[string]dbo4logist.OrderCounter, 1)
	}
	const counterName = "all"
	counter := logistSpace.Data.OrderCounters[counterName]
	counter.LastNumber++
	logistSpace.Data.OrderCounters[counterName] = counter

	if err := logistSpace.Data.Validate(); err != nil {
		return nil, fmt.Errorf("logistus team record is not valid: %w", err)
	}

	orderNumberPrefixed := logistSpace.Data.OrderCounters["all"].Prefix + strconv.Itoa(counter.LastNumber)
	order := dbo4logist.NewOrder(params.Space.ID, orderNumberPrefixed)
	fillOrderDtoFromRequest(order.Dto, request, params, userID)

	if err := addContactsFromCounterparties(ctx, tx, params.Space.ID, order.Dto); err != nil {
		return nil, fmt.Errorf("failed to add contacts from counterparties: %w", err)
	}

	if logistSpace.Record.Exists() {
		logistSpaceUpdates := []dal.Update{
			{Field: "orderCounters.all.lastNumber", Value: counter.LastNumber},
		}
		if err := tx.Update(ctx, logistSpace.Key, logistSpaceUpdates); err != nil {
			return nil, fmt.Errorf("failed to update logistus team record: %w", err)
		}
	} else if err := tx.Insert(ctx, logistSpace.Record); err != nil {
		return nil, fmt.Errorf("failed to insert logistus team record: %w", err)
	}

	//order.Data.OrderIDs = []string{params.Space.ContactID + ":" + orderNumberPrefixed}

	order.Dto.UpdateKeys()
	order.Dto.UpdateDates()
	if err := order.Dto.Validate(); err != nil {
		return nil, fmt.Errorf("order record is not valid before insert: %w", err)
	}
	if err := tx.Insert(ctx, order.Record); err != nil {
		return nil, fmt.Errorf("failed to insert order record: %w", err)
	}

	orderBrief = &dbo4logist.OrderBrief{
		ID:        order.Key.ID.(string),
		OrderBase: order.Dto.OrderBase,
	}

	return orderBrief, err
}

func fillOrderDtoFromRequest(orderDto *dbo4logist.OrderDbo, request dto4logist.CreateOrderRequest, params *dal4spaceus.SpaceWorkerParams, userID string) {
	orderDto.OrderBase = request.Order

	orderDto.Status = "active"
	orderDto.UserIDs = params.Space.Data.UserIDs
	orderDto.SpaceID = params.Space.ID
	orderDto.SpaceIDs = []string{params.Space.ID}
	modified := dbmodels2.Modified{
		By: userID,
		At: time.Now(),
	}
	orderDto.CreatedFields = with.CreatedFields{
		CreatedAtField: with.CreatedAtField{
			CreatedAt: modified.At,
		},
		CreatedByField: with.CreatedByField{
			CreatedBy: modified.By,
		},
	}
	orderDto.UpdatedFields = with.UpdatedFields{
		UpdatedAt: modified.At,
		UpdatedBy: modified.By,
	}

	if orderDto.Route != nil {
		if countryID := orderDto.Route.Origin.CountryID; slice.Index(orderDto.CountryIDs, countryID) < 0 {
			orderDto.CountryIDs = append(orderDto.CountryIDs, countryID)
			orderDto.CountryIDs = append(orderDto.CountryIDs, countryID+":origin")
		}
		if countryID := orderDto.Route.Destination.CountryID; slice.Index(orderDto.CountryIDs, countryID) < 0 {
			orderDto.CountryIDs = append(orderDto.CountryIDs, countryID)
			orderDto.CountryIDs = append(orderDto.CountryIDs, countryID+":destination")
		}
		for _, t := range orderDto.Route.TransitPoints {
			if countryID := t.CountryID; slice.Index(orderDto.CountryIDs, countryID) < 0 {
				orderDto.CountryIDs = append(orderDto.CountryIDs, countryID)
				orderDto.CountryIDs = append(orderDto.CountryIDs, countryID+":transit")
			}
		}
	}
	for size, count := range request.NumberOfContainers {
		for i := 1; i <= count; i++ {
			container := dbo4logist.OrderContainer{
				ID: fmt.Sprintf("%s%d", size, i),
				OrderContainerBase: dbo4logist.OrderContainerBase{
					Type: size,
				},
			}
			orderDto.Containers = append(orderDto.Containers, &container)
		}
	}
}

func addContactsFromCounterparties(ctx context.Context, tx dal.ReadTransaction, spaceID string, order *dbo4logist.OrderDbo) error {
	if len(order.Counterparties) == 0 {
		panic("at least 1 counterparty should be added to a new order")
	}
	contacts := make([]dal4contactus.ContactEntry, 0, len(order.Counterparties))
	records := make([]dal.Record, 0, len(order.Counterparties))
	contactIDs := make([]string, 0, len(order.Counterparties))
	for _, cp := range order.Counterparties {
		if slice.Index(contactIDs, cp.ContactID) < 0 {
			contactIDs = append(contactIDs, cp.ContactID)
			contact := dal4contactus.NewContactEntry(spaceID, cp.ContactID)
			contacts = append(contacts, contact)
			records = append(records, contact.Record)
		}
	}
	if err := tx.GetMulti(ctx, records); err != nil {
		return fmt.Errorf("failed to load contacts: %w", err)
	}

	for _, contact := range contacts {
		_, orderContact := order.GetContactByID(contact.ID)
		if orderContact == nil {
			orderContact = &dbo4logist.OrderContact{
				ID:        contact.ID,
				Type:      contact.Data.Type,
				Title:     contact.Data.Title,
				ParentID:  contact.Data.ParentID,
				CountryID: contact.Data.CountryID,
			}
			//if contact.Data.Address != nil {
			//	orderContact.Address = *contact.Data.Address
			//}
			//if orderContact.Address.CountryID == "" {
			//	orderContact.Address.CountryID = contact.Data.CountryID
			//}
			order.Contacts = append(order.Contacts, orderContact)
		}
	}

	return nil
}
